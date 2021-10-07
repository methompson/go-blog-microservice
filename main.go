package main

import (
	"flag"
	"syscall"

	"github.com/joho/godotenv"

	"methompson.com/blog-microservice/blogServer"
)

func main() {
	initPtr := flag.Bool("init", false, "Whether to run the admin initialization function")
	initUidPtr := flag.String("adminUid", "", "The UID of the user slated to be an admin")
	initNamePtr := flag.String("adminName", "", "The Name of the user slated to be an admin")

	flag.Parse()

	godotenv.Load()

	syscall.Umask(0)

	// fmt.Println(*initPtr)

	if initPtr != nil && *initPtr {
		blogServer.DoAdminInit(initUidPtr, initNamePtr)
		return
	}

	blogServer.MakeAndStartServer()
}
