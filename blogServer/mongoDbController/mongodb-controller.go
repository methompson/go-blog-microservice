package mongoDbController

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"

	"methompson.com/blog-microservice/blogServer/dbController"
	"methompson.com/blog-microservice/blogServer/logging"
)

type MongoDbController struct {
	MongoClient *mongo.Client
	dbName      string
}

func (mdb *MongoDbController) InitDatabase() error {
	return errors.New("Unimplemented")
}

func (mdb *MongoDbController) AddUserData(data *dbController.UserDataDocument) error {
	return errors.New("Unimplemented")
}

func (mdb *MongoDbController) GetUserDataById(id string) (*dbController.UserDataDocument, error) {
	return nil, errors.New("Unimplemented")
}

func (mdb *MongoDbController) EditUserData(data *dbController.UserDataDocument) error {
	return errors.New("Unimplemented")
}

func (mdb *MongoDbController) DeleteUserData(id string) error {
	return errors.New("Unimplemented")
}

func (mdb *MongoDbController) AddBlogPost(doc *dbController.BlogDocument) error {
	return errors.New("Unimplemented")
}

func (mdb *MongoDbController) GetBlogPostById(id string) (*dbController.BlogDocument, error) {
	return nil, errors.New("Unimplemented")
}

func (mdb *MongoDbController) GetBlogPostBySlug(slug string) (*dbController.BlogDocument, error) {
	return nil, errors.New("Unimplemented")
}

func (mdb *MongoDbController) EditBlogPost(doc *dbController.BlogDocument) error {
	return errors.New("Unimplemented")
}

func (mdb *MongoDbController) DeleteBlogPost(id string) error {
	return errors.New("Unimplemented")
}

func (mdb *MongoDbController) AddRequestLog(log *logging.RequestLogData) error {
	return errors.New("Unimplemented")
}

func (mdb *MongoDbController) AddInfoLog(log *logging.InfoLogData) error {
	return errors.New("Unimplemented")
}
