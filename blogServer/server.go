package blogServer

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
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
	mdbController, mdbControllerErr := mongoDbController.MakeMongoDbController(constants.AUTH_DB_NAME)

	if mdbControllerErr != nil {
		log.Fatal(mdbControllerErr.Error())
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

func (srv *BlogServer) ValidateIdToken(header AuthorizationHeader) {
	ctx := context.Background()
	client, err := srv.FirebaseApp.Auth(ctx)

	if err != nil {
		fmt.Println(err)
		return
	}

	token, err := client.VerifyIDToken(ctx, header.Token)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Verified ID Token: %v\n", token)
}
