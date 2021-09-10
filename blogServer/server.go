package blogServer

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"

	"methompson.com/blog-microservice/blogServer/constants"
	"methompson.com/blog-microservice/blogServer/dbController"
	"methompson.com/blog-microservice/blogServer/mongoDbController"
)

func MakeAndStartServer() {
	envErr := checkEnvVariables()

	if envErr != nil {
		log.Fatal("Error with environment variables")
	}

	srv, srvErr := MakeServer()

	if srvErr != nil {
		log.Fatal("Error making server")
	}

	srv.SetRoutes()

	srv.StartServer()
}

func MakeServer() (*BlogServer, error) {
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

	engine := gin.Default()

	// First we assign the pointer-to MongoDbController of mongoDbController to
	// a variable of type DatabaseController. Then we get the pointer-to DatabaseController and assign that to cont. We can use pointer-to DatabaseController to run InitController to initialize the AuthController.
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
