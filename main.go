package main

import (
	"context"
	"flag"
	"log"
	"os"
	"syscall"

	firebase "firebase.google.com/go/v4"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"

	"methompson.com/blog-microservice/blogServer"
	"methompson.com/blog-microservice/blogServer/constants"
)

func main() {
	initPtr := flag.Bool("init", false, "Whether to run the admin initialization function")
	initUidPtr := flag.String("adminUid", "", "The UID for a user to make an admin")

	flag.Parse()

	godotenv.Load()

	syscall.Umask(0)

	// fmt.Println(*initPtr)

	if initPtr != nil && *initPtr {
		doAdminInit(initUidPtr)
		return
	}

	blogServer.MakeAndStartServer()
}

func doAdminInit(initUidPtr *string) {
	print("Doing Admin Init")

	if initUidPtr == nil {
		log.Fatal("adminUid is required for setting admin role")
	}

	if len(os.Getenv(constants.GOOGLE_APPLICATION_CREDENTIALS)) == 0 {
		log.Fatal("No Google Application credential path listed")
	}

	ctx := context.Background()
	sa := option.WithCredentialsFile(os.Getenv(constants.GOOGLE_APPLICATION_CREDENTIALS))
	app, appErr := firebase.NewApp(ctx, nil, sa)

	if appErr != nil {
		log.Fatal(appErr.Error())
	}

	client, clientErr := app.Auth(ctx)

	if clientErr != nil {
		log.Fatal(appErr.Error())
	}

	// Set admin privilege on the user corresponding to uid.
	claims := map[string]interface{}{"role": "admin"}
	setClaimsErr := client.SetCustomUserClaims(ctx, *initUidPtr, claims)
	if setClaimsErr != nil {
		log.Fatalf(setClaimsErr.Error())
	}
}
