package blogServer

import (
	"context"
	"fmt"
	"log"

	"methompson.com/blog-microservice/blogServer/user"
)

func DoAdminInit(initUidPtr, initNamePtr *string) {
	fmt.Println("Doing Admin Init")

	if len(*initUidPtr) == 0 {
		log.Fatal("adminUid is required for setting admin role")
	}

	fmt.Println(*initUidPtr)

	// if len(*initEmailPtr) == 0 {
	// 	log.Fatal("adminEmail is required for setting admin role")
	// }

	// fmt.Println(*initEmailPtr)

	mdbController, mdbControllerErr := makeAndInitDatabase()

	if mdbControllerErr != nil {
		log.Fatal("Error Initializing Database: ", mdbControllerErr.Error())
	}

	ctx := context.Background()

	app, appErr := makeFirebaseApp()

	if appErr != nil {
		log.Fatal(appErr.Error())
	}

	client, clientErr := app.Auth(ctx)

	if clientErr != nil {
		log.Fatal(appErr.Error())
	}

	// Get the user initially
	getUser, userErr := client.GetUser(ctx, *initUidPtr)

	if clientErr != nil {
		log.Fatal(userErr.Error())
	}

	if getUser == nil {
		log.Fatal("No user retrieved with that UID")
	}

	email := getUser.UserInfo.Email

	fmt.Println(email)

	info := user.UserInformation{
		Uid:    *initUidPtr,
		Email:  email,
		Name:   *initNamePtr,
		Active: true,
		Role:   user.Admin,
	}

	editUserErr := mdbController.AddUserInformation(&info)

	if editUserErr != nil {
		log.Fatal("Error adding user information: ", editUserErr.Error())
	}

	// Set admin privilege on the user corresponding to uid.
	claims := map[string]interface{}{
		"role":   "admin",
		"name":   *initNamePtr,
		"active": true,
	}
	setClaimsErr := client.SetCustomUserClaims(ctx, *initUidPtr, claims)

	if setClaimsErr != nil {
		log.Fatalf(setClaimsErr.Error())
	}
}
