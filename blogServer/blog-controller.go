package blogServer

import (
	"github.com/gosimple/slug"
	"methompson.com/blog-microservice/blogServer/dbController"
	"methompson.com/blog-microservice/blogServer/logging"
	"methompson.com/blog-microservice/blogServer/user"
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

func (bc *BlogController) AddUserData(userInfo *user.UserInformation) {}

func (bc *BlogController) GetUserData() *user.UserInformation {
	ui := user.UserInformation{}

	return &ui
}

func (bc *BlogController) EditUserData()    {}
func (bc *BlogController) DeletedUserData() {}

func (bc *BlogController) AddBlogPost(blogBody AddBlogBody) (id string, slug string, err error) {
	blogDocument := blogBody.GetBlogDocument()

	if !bc.isValidSlug(blogDocument.Slug) {
		blogDocument.Slug = bc.slugify(blogDocument.Title)
	}

	addBlogId, addBlogErr := (*bc.DBController).AddBlogPost(blogDocument)

	if addBlogErr != nil {
		return "", "", addBlogErr
	}

	return addBlogId, blogDocument.Slug, nil
}

func (bc *BlogController) GetBlogPostById(id string) (*dbController.BlogDocument, error) {
	return (*bc.DBController).GetBlogPostById(id)
}

func (bc *BlogController) GetBlogPostBySlug(slug string) (*dbController.BlogDocument, error) {
	return (*bc.DBController).GetBlogPostBySlug(slug)
}

func (bc *BlogController) GetBlogPosts(page int, pagination int) ([]*dbController.BlogDocument, error) {
	_pagination := pagination

	if _pagination <= 0 {
		_pagination = 10
	}

	return (*bc.DBController).GetBlogPosts(page, _pagination)
}

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

func (bc *BlogController) AddLogger(logger *logging.BlogLogger) {
	bc.Loggers = append(bc.Loggers, logger)
}

func (bc *BlogController) isValidSlug(_slug string) bool {
	return slug.IsSlug(_slug)
}

func (bc *BlogController) slugify(_slug string) string {
	return slug.Make(_slug)
}
