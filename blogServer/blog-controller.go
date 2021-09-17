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

func (bc *BlogController) AddBlogPost(blogBody AddBlogBody) (id string, err error) {
	blogDocument := blogBody.GetBlogDocument()

	if !bc.isValidSlug(blogDocument.Slug) {
		blogDocument.Slug = bc.slugify(blogDocument.Title)
	}

	addBlogId, addBlogErr := (*bc.DBController).AddBlogPost(blogDocument)

	if addBlogErr != nil {
		return "", addBlogErr
	}

	return addBlogId, nil
}

func (bc *BlogController) GetBlogPostById()   {}
func (bc *BlogController) GetBlogPostBySlug() {}

// TODO Check that the data is valid (e.g. the slug)
func (bc *BlogController) EditBlogPost(body EditBlogBody) error {
	blogDocument := body.GetBlogDocument()

	if blogDocument.Slug != nil && bc.isValidSlug(*blogDocument.Slug) {
		if blogDocument.Title == nil {
			return NewInputError("invalid slug and no title")
		}

		s := bc.slugify(*blogDocument.Title)
		blogDocument.Slug = &s
	}

	return (*bc.DBController).EditBlogPost(blogDocument)
}

func (bc *BlogController) DeleteBlogPost(body DeleteBlogBody) error {
	blogDocument := body.GetBlogDocument()

	return (*bc.DBController).DeleteBlogPost(blogDocument)
}

func (bc *BlogController) AddLogger()     {}
func (bc *BlogController) AddRequestLog() {}
func (bc *BlogController) AddInfoLog()    {}

func (bc *BlogController) isValidSlug(slug string) bool {
	return true
}

func (bc *BlogController) slugify(slug string) string {
	return slug
}
