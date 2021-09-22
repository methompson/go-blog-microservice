package blogServer

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"

	"methompson.com/blog-microservice/blogServer/constants"
	"methompson.com/blog-microservice/blogServer/dbController"
	"methompson.com/blog-microservice/blogServer/logging"
	"methompson.com/blog-microservice/blogServer/mongoDbController"
)

func MakeAndStartServer() {
	envErr := checkEnvVariables()

	if envErr != nil {
		log.Fatal("Error with environment variables")
	}

	blogServer, srvErr := makeServer()

	if srvErr != nil {
		log.Fatal("Error making server")
	}

	// We run this after creating a server, but before setting routes. Any
	// route set BEFORE this won't actually use this.
	if !DebugMode() {
		errs := configureReleaseLogging(blogServer)

		if len(errs) > 0 {
			for _, err := range errs {
				print(err.Error() + "\n")
			}
		}
		addLogging(blogServer)

		addRecovery(blogServer)
	}

	blogServer.SetRoutes()

	blogServer.StartServer()
}

func configureReleaseLogging(bs *BlogServer) []error {
	errs := make([]error, 0)
	controller := &bs.BlogController

	if os.Getenv(constants.DB_LOGGING) == "true" {
		// We set the logger to a database logger
		// First, we manipulate the pointers in order to add the DBController to the logger
		// in order to log release data to the database.
		var dbController logging.BlogLogger = *controller.DBController
		controller.AddLogger(&dbController)
	}

	if os.Getenv(constants.FILE_LOGGING) == "true" {
		// We can also log to a file
		var fileLogger logging.BlogLogger
		var fileLoggerErr error

		fileLogger, fileLoggerErr = logging.MakeNewFileLogger(os.Getenv(constants.FILE_LOGGING_PATH), "logs.log")

		if fileLoggerErr != nil {
			errs = append(errs, fileLoggerErr)
		}
		controller.AddLogger(&fileLogger)
	}

	if os.Getenv(constants.CONSOLE_LOGGING) == "true" {
		var consoleLogger logging.BlogLogger = &logging.ConsoleLogger{}

		controller.AddLogger(&consoleLogger)
	}

	return errs
}

// TODO figure out recovery
func addRecovery(bs *BlogServer) {
	bs.GinEngine.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		msg := "Unknown Error"
		if err, ok := recovered.(string); ok {
			msg = fmt.Sprintf("error: %s", err)
			c.String(http.StatusInternalServerError, msg)
		}

		errorLog := logging.InfoLogData{
			Timestamp: time.Now(),
			Type:      "error",
			Message:   msg,
		}

		for _, logger := range bs.BlogController.Loggers {
			l := *logger
			l.AddInfoLog(&errorLog)
		}

		c.AbortWithStatus(http.StatusInternalServerError)
	}))
}

func makeServer() (*BlogServer, error) {
	mdbController, mdbControllerErr := mongoDbController.MakeMongoDbController(constants.BLOG_DB_NAME)

	if mdbControllerErr != nil {
		log.Fatal(mdbControllerErr.Error())
	}

	initDbErr := mdbController.InitDatabase()

	if initDbErr != nil {
		log.Fatal("Error Initializing Database: ", initDbErr.Error())
	}

	app, err := makeFirebaseApp()

	if err != nil {
		return nil, errors.New("error making firebase app")
	}

	engine := makeGinEngine()

	// First we assign the pointer-to MongoDbController of mongoDbController to
	// a variable of type DatabaseController. Then we get the pointer-to
	// DatabaseController and assign that to cont. We can use pointer-to
	// DatabaseController to run InitController to initialize the AuthController.
	var passedController dbController.DatabaseController = mdbController
	ptrToCont := &passedController

	srv := BlogServer{
		FirebaseApp:    app,
		BlogController: InitController(ptrToCont),
		GinEngine:      engine,
	}

	return &srv, nil
}

func makeFirebaseApp() (*firebase.App, error) {
	sa := option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	app, err := firebase.NewApp(context.Background(), nil, sa)

	if err != nil {
		return nil, err
	}

	return app, nil
}

func makeGinEngine() *gin.Engine {
	// We run this prior to creating a server. Any gin engine created prior
	// to running SetMode won't include this configuration.
	if !DebugMode() {
		gin.SetMode(gin.ReleaseMode)
	}

	if os.Getenv("GIN_MODE") == "release" {
		return gin.New()
	}

	return gin.Default()
}

func addLogging(as *BlogServer) {
	as.GinEngine.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// address := authUtils.GetRemoteAddressIP(param.Request.RemoteAddr)
		requestData := logging.RequestLogData{
			Timestamp:    param.TimeStamp,
			Type:         "request",
			ClientIP:     param.ClientIP,
			Method:       param.Method,
			Path:         param.Path,
			Protocol:     param.Request.Proto,
			StatusCode:   param.StatusCode,
			Latency:      param.Latency,
			UserAgent:    param.Request.UserAgent(),
			ErrorMessage: param.ErrorMessage,
		}

		for _, logger := range as.BlogController.Loggers {
			l := *logger
			l.AddRequestLog(&requestData)
		}

		return ""
	}))
}

// TODO check environment variables
func checkEnvVariables() error {
	return nil
}

type BlogServer struct {
	FirebaseApp    *firebase.App
	BlogController BlogController
	GinEngine      *gin.Engine
}

func (srv *BlogServer) StartServer() {
	srv.GinEngine.Run()
}

func (srv *BlogServer) ValidateIdToken(header AuthorizationHeader) (*auth.Token, error) {
	ctx := context.Background()
	client, clientErr := srv.FirebaseApp.Auth(ctx)

	if clientErr != nil {
		fmt.Println(clientErr)
		return nil, clientErr
	}

	token, tokenErr := client.VerifyIDToken(ctx, header.Token)

	if tokenErr != nil {
		fmt.Println(tokenErr)
		return nil, tokenErr
	}

	return token, nil
}

func (srv *BlogServer) GetAuthorizationHeader(ctx *gin.Context) (*auth.Token, error) {
	var header AuthorizationHeader

	// No Token Error
	if headerErr := ctx.ShouldBindHeader(&header); headerErr != nil {
		fmt.Println(headerErr)
		ctx.Data(401, "text/html; charset=utf-8", make([]byte, 0))
		return nil, headerErr
	}

	token, tokenErr := srv.ValidateIdToken(header)

	if tokenErr != nil {
		return nil, tokenErr
	}

	return token, nil
}

func (srv *BlogServer) GetRoleFromToken(token *auth.Token) (string, error) {
	roleInt := token.Claims["role"]

	role, ok := roleInt.(string)

	if !ok {
		fmt.Println(role)
		return "", errors.New("role is not a string")
	}

	return role, nil
}

func (srv *BlogServer) GetTokenAndRoleFromHeader(ctx *gin.Context) (*auth.Token, string, error) {
	token, tokenErr := srv.GetAuthorizationHeader(ctx)

	// No Token Error
	if tokenErr != nil {
		return nil, "", tokenErr
	}

	role, roleErr := srv.GetRoleFromToken(token)

	if roleErr != nil {
		return nil, "", roleErr
	}

	return token, role, nil
}

func (srv *BlogServer) CanEditBlog(role string) bool {
	return role == constants.USER_ADMIN || role == constants.USER_EDITOR
}
