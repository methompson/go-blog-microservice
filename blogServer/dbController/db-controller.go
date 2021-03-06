package dbController

import (
	"methompson.com/blog-microservice/blogServer/logging"
	"methompson.com/blog-microservice/blogServer/user"
)

type DatabaseController interface {
	InitDatabase() error

	AddBlogPost(doc *AddBlogDocument) (id string, err error)
	GetBlogPostById(id string) (*BlogDocument, error)
	GetBlogPostBySlug(slug string) (*BlogDocument, error)
	GetBlogPosts(page int, pagination int) ([]*BlogDocument, error)
	EditBlogPost(doc *EditBlogDocument) error
	DeleteBlogPost(doc *DeleteBlogDocument) error

	AddUserInformation(info *user.UserInformation) error

	AddRequestLog(log *logging.RequestLogData) error
	AddInfoLog(log *logging.InfoLogData) error
}
