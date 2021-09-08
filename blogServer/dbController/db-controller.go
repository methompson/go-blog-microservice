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

	AddBlogPost(doc *BlogDocument) error
	GetBlogPostById(id string) (*BlogDocument, error)
	GetBlogPostBySlug(slug string) (*BlogDocument, error)
	EditBlogPost(doc *BlogDocument) error
	DeleteBlogPost(id string) error

	AddRequestLog(log *logging.RequestLogData) error
	AddInfoLog(log *logging.InfoLogData) error
}
