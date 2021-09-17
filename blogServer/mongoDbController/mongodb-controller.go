package mongoDbController

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
	"methompson.com/blog-microservice/blogServer/dbController"
	"methompson.com/blog-microservice/blogServer/logging"
)

type MongoDbController struct {
	MongoClient *mongo.Client
	dbName      string
}

// getCollection is a convenience function that performs a function used regularly
// throughout the Mongodbc. It accepts a collectionName string for the
// specific collection you want to retrieve, and returns a collection, context and
// cancel function.
func (mdbc *MongoDbController) getCollection(collectionName string) (*mongo.Collection, context.Context, context.CancelFunc) {
	// Write the hash to the database
	collection := mdbc.MongoClient.Database(mdbc.dbName).Collection(collectionName)
	backCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	return collection, backCtx, cancel
}

func (mdbc *MongoDbController) initUserCollection(dbName string) error {
	db := mdbc.MongoClient.Database(dbName)

	jsonSchema := bson.M{
		"bsonType": "object",
		"required": []string{"userId"},
		"properties": bson.M{
			"userId": bson.M{
				"bsonType":    "string",
				"description": "userId must be a string",
			},
			"name": bson.M{
				"bsonType":    "string",
				"description": "name must be a string",
			},
		},
	}

	colOpts := options.CreateCollection().SetValidator(bson.M{"$jsonSchema": jsonSchema})

	createCollectionErr := db.CreateCollection(context.TODO(), "users", colOpts)

	if createCollectionErr != nil {
		return dbController.NewDBError(createCollectionErr.Error())
	}

	models := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "userId", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	opts := options.CreateIndexes().SetMaxTime(2 * time.Second)

	collection, _, _ := mdbc.getCollection("users")
	_, setIndexErr := collection.Indexes().CreateMany(context.TODO(), models, opts)

	if setIndexErr != nil {
		return dbController.NewDBError(setIndexErr.Error())
	}

	return nil
}

func (mdbc *MongoDbController) initBlogCollection(dbName string) error {
	db := mdbc.MongoClient.Database(dbName)

	jsonSchema := bson.M{
		"bsonType": "object",
		"required": []string{"title", "slug", "authorId", "dateAdded", "updateAuthorId", "dateUpdated"},
		"properties": bson.M{
			"title": bson.M{
				"bsonType":    "string",
				"description": "titlr must be a string",
			},
			"slug": bson.M{
				"bsonType":    "string",
				"description": "slug must be a string",
			},
			"body": bson.M{
				"bsonType":    "string",
				"description": "slug must be a string",
			},
			"tags": bson.M{
				"bsonType":    "array",
				"description": "tags must be an array",
			},
			"authorId": bson.M{
				"bsonType":    "string",
				"description": "authorId must be a string",
			},
			"dateAdded": bson.M{
				"bsonType":    "timestamp",
				"description": "dateAdded must be a timestamp",
			},
			"updateAuthorId": bson.M{
				"bsonType":    "string",
				"description": "updateAuthorId must be a string",
			},
			"dateUpdated": bson.M{
				"bsonType":    "timestamp",
				"description": "dateUpdated must be a timestamp",
			},
		},
	}

	collectionName := "blogPosts"

	colOpts := options.CreateCollection().SetValidator(bson.M{"$jsonSchema": jsonSchema})

	createCollectionErr := db.CreateCollection(context.TODO(), collectionName, colOpts)

	if createCollectionErr != nil {
		return dbController.NewDBError(createCollectionErr.Error())
	}

	models := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "slug", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	opts := options.CreateIndexes().SetMaxTime(2 * time.Second)

	collection, _, _ := mdbc.getCollection(collectionName)
	_, setIndexErr := collection.Indexes().CreateMany(context.TODO(), models, opts)

	if setIndexErr != nil {
		return dbController.NewDBError(setIndexErr.Error())
	}

	return nil
}

func (mdbc *MongoDbController) initLoggingCollection(dbName string) error {
	db := mdbc.MongoClient.Database(dbName)

	jsonSchema := bson.M{
		"bsonType": "object",
		"required": []string{"timestamp", "type"},
		"properties": bson.M{
			"timestamp": bson.M{
				"bsonType":    "timestamp",
				"description": "timestamp is required and must be a timestamp",
			},
			"type": bson.M{
				"bsonType":    "string",
				"description": "type is required and must be a string",
			},
		},
	}

	colOpts := options.CreateCollection().SetValidator(bson.M{"$jsonSchema": jsonSchema})
	colOpts.SetCapped(true)
	colOpts.SetSizeInBytes(100000)

	createCollectionErr := db.CreateCollection(context.TODO(), "logging", colOpts)

	if createCollectionErr != nil {
		return dbController.NewDBError(createCollectionErr.Error())
	}

	return nil
}

func (mdbc *MongoDbController) InitDatabase() error {
	userCreationErr := mdbc.initUserCollection(mdbc.dbName)

	if userCreationErr != nil && !strings.Contains(userCreationErr.Error(), "Collection already exists") {
		return userCreationErr
	}

	blogCreationErr := mdbc.initBlogCollection(mdbc.dbName)

	if blogCreationErr != nil && !strings.Contains(blogCreationErr.Error(), "Collection already exists") {
		return blogCreationErr
	}

	loggingCreationErr := mdbc.initLoggingCollection(mdbc.dbName)

	if loggingCreationErr != nil && !strings.Contains(loggingCreationErr.Error(), "Collection already exists") {
		return loggingCreationErr
	}

	return nil
}

func (mdbc *MongoDbController) AddUserData(data *dbController.UserDataDocument) error {
	return errors.New("Unimplemented")
}

func (mdbc *MongoDbController) GetUserDataById(id string) (*dbController.UserDataDocument, error) {
	return nil, errors.New("Unimplemented")
}

func (mdbc *MongoDbController) EditUserData(data *dbController.UserDataDocument) error {
	return errors.New("Unimplemented")
}

func (mdbc *MongoDbController) DeleteUserData(id string) error {
	return errors.New("Unimplemented")
}

func (mdbc *MongoDbController) AddBlogPost(doc *dbController.AddBlogDocument) (string, error) {
	collection, backCtx, cancel := mdbc.getCollection("blogPosts")
	defer cancel()

	dateAdded := primitive.Timestamp{T: uint32(doc.DateAdded.Unix())}

	insert := bson.M{
		"title":     doc.Title,
		"slug":      doc.Slug,
		"body":      doc.Body,
		"authorId":  doc.AuthorId,
		"dateAdded": dateAdded,
	}

	if doc.Tags != nil {
		insert["tags"] = *doc.Tags
	}

	if doc.UpdateAuthorId != nil {
		insert["updateAuthorId"] = *doc.UpdateAuthorId
	} else {
		insert["updateAuthorId"] = doc.AuthorId
	}

	if doc.DateUpdated != nil {
		insert["dateUpdated"] = primitive.Timestamp{T: uint32((*doc.DateUpdated).Unix())}
	} else {
		insert["dateUpdated"] = dateAdded
	}

	insertResult, mdbErr := collection.InsertOne(backCtx, insert)

	if mdbErr != nil {
		err := mdbErr.Error()
		print("Add blog error. Error: " + err + "\n")

		if strings.Contains(err, "duplicate key error") {
			msg := "Duplicate blog post."
			if strings.Contains(err, "slug") {
				msg = msg + " Blog Post with slug '" + doc.Slug + "' already exists."
			}

			return "", dbController.NewDuplicateEntryError(msg)
		}

		return "", dbController.NewDBError(mdbErr.Error())
	}

	objectId, idOk := insertResult.InsertedID.(primitive.ObjectID)

	if !idOk {
		return "", dbController.NewDBError("invalid id returned by database")
	}

	return objectId.Hex(), nil
}

func (mdbc *MongoDbController) GetBlogPostById(id string) (*dbController.BlogDocument, error) {
	return nil, errors.New("Unimplemented")
}

func (mdbc *MongoDbController) GetBlogPostBySlug(slug string) (*dbController.BlogDocument, error) {
	return nil, errors.New("Unimplemented")
}

func (mdbc *MongoDbController) EditBlogPost(doc *dbController.EditBlogDocument) error {
	collection, backCtx, cancel := mdbc.getCollection("blogPosts")
	defer cancel()

	print("Editing Blog Post\n")

	id, idErr := primitive.ObjectIDFromHex(doc.Id)
	if idErr != nil {
		return dbController.NewInvalidInputError("Invalid User ID")
	}

	filter := bson.M{"_id": id}

	values := bson.M{}

	if doc.Title != nil {
		values["title"] = *doc.Title
	}

	if doc.Slug != nil {
		values["slug"] = *doc.Slug
	}

	if doc.Body != nil {
		values["body"] = *doc.Body
	}

	if doc.Tags != nil {
		values["tags"] = *doc.Tags
	}

	if doc.AuthorId != nil {
		values["authorId"] = *doc.AuthorId
	}

	if doc.DateUpdated != nil {
		values["dateAdded"] = primitive.Timestamp{T: uint32((*doc.DateAdded).Unix())}
	}

	if doc.UpdateAuthorId != nil {
		values["updateAuthorId"] = doc.UpdateAuthorId
	}

	if doc.DateUpdated != nil {
		values["dateUpdated"] = primitive.Timestamp{T: uint32((*doc.DateUpdated).Unix())}
	}

	update := bson.M{
		"$set": values,
	}

	result, mdbErr := collection.UpdateOne(backCtx, filter, update)

	if mdbErr != nil {
		err := mdbErr.Error()
		print("Edit blog error. Error: " + err + "\n")

		if strings.Contains(err, "duplicate key error") {
			msg := "Duplicate blog post."
			if strings.Contains(err, "slug") {
				msg = msg + " Blog Post with slug '" + *doc.Slug + "' already exists."
			}

			return dbController.NewDuplicateEntryError(msg)
		}

		return dbController.NewDBError(mdbErr.Error())
	}

	if result.MatchedCount == 0 {
		return dbController.NewInvalidInputError("id did not match any blog posts")
	}

	return nil
}

func (mdbc *MongoDbController) DeleteBlogPost(doc *dbController.DeleteBlogDocument) error {
	return errors.New("Unimplemented")
}

func (mdbc *MongoDbController) AddRequestLog(log *logging.RequestLogData) error {
	return errors.New("Unimplemented")
}

func (mdbc *MongoDbController) AddInfoLog(log *logging.InfoLogData) error {
	return errors.New("Unimplemented")
}
