package blogServer

import (
	"methompson.com/blog-microservice/blogServer/dbController"
	"methompson.com/blog-microservice/blogServer/logging"
)

type BlogController struct {
	DBController *dbController.DatabaseController
	Loggers      []*logging.BlogLogger
}

// The DatabaseController should already be initialized before getting
// passed to the InitController function
func InitController(dbc *dbController.DatabaseController) BlogController {
	bc := BlogController{
		DBController: dbc,
		Loggers:      make([]*logging.BlogLogger, 0),
	}

	return bc
}

func (bc *BlogController) AddUserData()     {}
func (bc *BlogController) GetUserData()     {}
func (bc *BlogController) EditUserData()    {}
func (bc *BlogController) DeletedUserData() {}

func (bc *BlogController) AddBlogPost()       {}
func (bc *BlogController) GetBlogPostById()   {}
func (bc *BlogController) GetBlogPostBySlug() {}
func (bc *BlogController) EditBlogPost()      {}
func (bc *BlogController) DeleteBlogPost()    {}

func (bc *BlogController) AddLogger()     {}
func (bc *BlogController) AddRequestLog() {}
func (bc *BlogController) AddInfoLog()    {}
