package main

import (
	"syscall"

	"github.com/joho/godotenv"
	"methompson.com/blog-microservice/blogServer"
)

func main() {
	godotenv.Load()

	syscall.Umask(0)
	blogServer.MakeAndStartServer()
}
