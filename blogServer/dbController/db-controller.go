package dbController

import (
	"methompson.com/blog-microservice/blogServer/logging"
)

type DatabaseController interface {
	InitDatabase() error

	AddUserData(data *UserDataDocument) error
	GetUserDataById(id string) (*UserDataDocument, error)
	EditUserData(data *UserDataDocument) error
	DeleteUserData(id string) error

	AddBlogPost(doc *AddBlogDocument) (id string, err error)
	GetBlogPostById(id string) (*BlogDocument, error)
	GetBlogPostBySlug(slug string) (*BlogDocument, error)
	EditBlogPost(doc *EditBlogDocument) error
	DeleteBlogPost(doc *DeleteBlogDocument) error

	AddRequestLog(log *logging.RequestLogData) error
	AddInfoLog(log *logging.InfoLogData) error
}
