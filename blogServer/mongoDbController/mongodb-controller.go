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
	"methompson.com/blog-microservice/blogServer/user"
)

const BLOG_COLLECTION = "blogPosts"
const LOGGING_COLLECTION = "logging"
const USER_COLLECTION = "users"

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

func (mdbc *MongoDbController) initBlogCollection(dbName string) error {
	db := mdbc.MongoClient.Database(dbName)

	jsonSchema := bson.M{
		"bsonType": "object",
		"required": []string{"title", "slug", "authorId", "dateAdded", "updateAuthorId", "dateUpdated"},
		"properties": bson.M{
			"title": bson.M{
				"bsonType":    "string",
				"description": "title must be a string",
			},
			"slug": bson.M{
				"bsonType":    "string",
				"description": "slug must be a string",
			},
			"body": bson.M{
				"bsonType":    "string",
				"description": "body must be a string",
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

	colOpts := options.CreateCollection().SetValidator(bson.M{"$jsonSchema": jsonSchema})

	createCollectionErr := db.CreateCollection(context.TODO(), BLOG_COLLECTION, colOpts)

	if createCollectionErr != nil {
		return dbController.NewDBError(createCollectionErr.Error())
	}

	models := []mongo.IndexModel{
		{
			Keys:    bson.M{"slug": 1},
			Options: options.Index().SetUnique(true),
		},
	}

	opts := options.CreateIndexes().SetMaxTime(2 * time.Second)

	collection, _, _ := mdbc.getCollection(BLOG_COLLECTION)
	_, setIndexErr := collection.Indexes().CreateMany(context.TODO(), models, opts)

	if setIndexErr != nil {
		return dbController.NewDBError(setIndexErr.Error())
	}

	return nil
}

func (mdbc *MongoDbController) initUserCollection(dbName string) error {
	db := mdbc.MongoClient.Database(dbName)

	jsonSchema := bson.M{
		"bsonType": "object",
		// "required": []string{"uid"},
		"required": []string{"uid", "name", "email", "active", "role"},
		"properties": bson.M{
			"uid": bson.M{
				"bsonType":    "string",
				"description": "uid must be a string",
			},
			"name": bson.M{
				"bsonType":    "string",
				"description": "name must be a string",
			},
			"email": bson.M{
				"bsonType":    "string",
				"description": "email must be a string",
			},
			"active": bson.M{
				"bsonType":    "bool",
				"description": "active must be a bool",
			},
			"role": bson.M{
				"bsonType":    "string",
				"description": "role must be a string",
			},
		},
	}

	colOpts := options.CreateCollection().SetValidator(bson.M{"$jsonSchema": jsonSchema})

	createCollectionErr := db.CreateCollection(context.TODO(), USER_COLLECTION, colOpts)

	if createCollectionErr != nil {
		return dbController.NewDBError(createCollectionErr.Error())
	}

	models := []mongo.IndexModel{
		{
			Keys:    bson.M{"uid": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.M{"email": 1},
			Options: options.Index().SetUnique(true),
		},
	}

	opts := options.CreateIndexes().SetMaxTime(2 * time.Second)

	collection, _, _ := mdbc.getCollection(USER_COLLECTION)
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

	createCollectionErr := db.CreateCollection(context.TODO(), LOGGING_COLLECTION, colOpts)

	if createCollectionErr != nil {
		return dbController.NewDBError(createCollectionErr.Error())
	}

	return nil
}

func (mdbc *MongoDbController) InitDatabase() error {
	blogCreationErr := mdbc.initBlogCollection(mdbc.dbName)

	if blogCreationErr != nil && !strings.Contains(blogCreationErr.Error(), "Collection already exists") {
		return blogCreationErr
	}

	userCreationErr := mdbc.initUserCollection(mdbc.dbName)

	if userCreationErr != nil && !strings.Contains(userCreationErr.Error(), "Collection already exists") {
		return userCreationErr
	}

	loggingCreationErr := mdbc.initLoggingCollection(mdbc.dbName)

	if loggingCreationErr != nil && !strings.Contains(loggingCreationErr.Error(), "Collection already exists") {
		return loggingCreationErr
	}

	return nil
}

func (mdbc *MongoDbController) AddBlogPost(doc *dbController.AddBlogDocument) (string, error) {
	collection, backCtx, cancel := mdbc.getCollection(BLOG_COLLECTION)
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

func (mdbc *MongoDbController) GetAggregationStages() (projectStage, authorLookupStage, updateAuthorLookupStage *bson.D) {
	ps := bson.D{
		{
			Key: "$project", Value: bson.M{
				"body":           1,
				"slug":           1,
				"title":          1,
				"authorId":       1,
				"updateAuthorId": 1,
				"dateAdded":      1,
				"dateUpdated":    1,
				"tags":           1,
			},
		},
	}

	als := bson.D{{
		Key: "$lookup",
		Value: bson.M{
			"from":         USER_COLLECTION,
			"localField":   "authorId",
			"foreignField": "uid",
			"as":           "author",
		},
	}}

	uals := bson.D{{
		Key: "$lookup",
		Value: bson.M{
			"from":         USER_COLLECTION,
			"localField":   "updateAuthorId",
			"foreignField": "uid",
			"as":           "updateAuthor",
		},
	}}

	return &ps, &als, &uals
}

func (mdbc *MongoDbController) GetBlogPostWithMatcher(matchStage *bson.D) (*dbController.BlogDocument, error) {
	collection, backCtx, cancel := mdbc.getCollection(BLOG_COLLECTION)
	defer cancel()

	projectStage, authorLookupStage, updateAuthorLookupStage := mdbc.GetAggregationStages()

	limitStage := bson.D{{
		Key:   "$limit",
		Value: int32(1),
	}}

	cursor, aggErr := collection.Aggregate(backCtx, mongo.Pipeline{
		*matchStage,
		*projectStage,
		limitStage,
		*authorLookupStage,
		*updateAuthorLookupStage,
	})

	if aggErr != nil {
		return nil, dbController.NewDBError("error getting data from database: " + aggErr.Error())
	}

	var results []BlogDocResult
	if allErr := cursor.All(backCtx, &results); allErr != nil {
		return nil, dbController.NewDBError("error parsing results: " + allErr.Error())
	}

	if len(results) < 1 {
		return nil, dbController.NewNoResultsError("")
	}

	var post *dbController.BlogDocument = results[0].GetBlogDocument()

	return post, nil
}

func (mdbc *MongoDbController) GetBlogPostById(id string) (*dbController.BlogDocument, error) {
	idObj, idObjErr := primitive.ObjectIDFromHex(id)

	if idObjErr != nil {
		return nil, dbController.NewInvalidInputError("invalid id")
	}

	matchStage := bson.D{{Key: "$match", Value: bson.M{
		"_id": idObj,
	}}}

	return mdbc.GetBlogPostWithMatcher(&matchStage)
}

func (mdbc *MongoDbController) GetBlogPostBySlug(slug string) (*dbController.BlogDocument, error) {
	matchStage := bson.D{{Key: "$match", Value: bson.M{
		"slug": slug,
	}}}

	return mdbc.GetBlogPostWithMatcher(&matchStage)
}

func (mdbc *MongoDbController) GetBlogPosts(page int, pagination int) ([]*dbController.BlogDocument, error) {
	collection, backCtx, cancel := mdbc.getCollection(BLOG_COLLECTION)
	defer cancel()

	matchStage := bson.D{{Key: "$match", Value: bson.M{}}}

	projectStage, authorLookupStage, updateAuthorLookupStage := mdbc.GetAggregationStages()

	sortStage := bson.D{{
		Key: "$sort",
		Value: bson.M{
			"dateAdded": -1,
		},
	}}

	limitStage := bson.D{{
		Key:   "$limit",
		Value: int32(pagination),
	}}

	skipStage := bson.D{{
		Key:   "$skip",
		Value: int64((page - 1) * pagination),
	}}

	cursor, aggErr := collection.Aggregate(backCtx, mongo.Pipeline{
		matchStage,
		*projectStage,
		sortStage,
		skipStage,
		limitStage,
		*authorLookupStage,
		*updateAuthorLookupStage,
	})

	if aggErr != nil {
		return nil, dbController.NewDBError("")
	}

	var results []BlogDocResult
	if allErr := cursor.All(backCtx, &results); allErr != nil {
		return nil, errors.New("error parsing results")
	}

	var posts []*dbController.BlogDocument = []*dbController.BlogDocument{}
	for _, v := range results {
		posts = append(posts, v.GetBlogDocument())
	}

	return posts, nil
}

func (mdbc *MongoDbController) EditBlogPost(doc *dbController.EditBlogDocument) error {
	collection, backCtx, cancel := mdbc.getCollection(BLOG_COLLECTION)
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
	collection, backCtx, cancel := mdbc.getCollection(BLOG_COLLECTION)
	defer cancel()

	print("Deleting Blog Post\n")

	id, idErr := primitive.ObjectIDFromHex(doc.Id)
	if idErr != nil {
		return dbController.NewInvalidInputError("Invalid User ID")
	}

	delResult, delErr := collection.DeleteOne(
		backCtx,
		bson.M{
			"_id": id,
		},
	)

	if delResult.DeletedCount == 0 {
		return dbController.NewInvalidInputError("invalid id. no blog posts deleted")
	}

	if delErr != nil {
		return dbController.NewDBError(delErr.Error())
	}

	return nil
}

func (mdbc *MongoDbController) AddUserInformation(info *user.UserInformation) error {
	collection, backCtx, cancel := mdbc.getCollection(USER_COLLECTION)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"uid":    info.Uid,
			"name":   info.Name,
			"email":  info.Email,
			"active": info.Active,
			"role":   info.Role.String(),
		},
	}

	opts := options.Update().SetUpsert(true)

	filter := bson.M{}

	updateResult, mdbErr := collection.UpdateOne(backCtx, filter, update, opts)

	if mdbErr != nil {
		return mdbErr
	}

	if updateResult.MatchedCount == 0 && updateResult.UpsertedCount == 0 {
		return dbController.NewInvalidInputError("id did not match any blog posts")
	}

	return nil
}

func (mdbc *MongoDbController) AddRequestLog(log *logging.RequestLogData) error {
	collection, backCtx, cancel := mdbc.getCollection(LOGGING_COLLECTION)
	defer cancel()

	insert := bson.M{
		"timestamp":    primitive.Timestamp{T: uint32(log.Timestamp.Unix())},
		"type":         log.Type,
		"clientIP":     log.ClientIP,
		"method":       log.Method,
		"path":         log.Path,
		"protocol":     log.Protocol,
		"statusCode":   log.StatusCode,
		"latency":      log.Latency,
		"userAgent":    log.UserAgent,
		"errorMessage": log.ErrorMessage,
	}

	_, mdbErr := collection.InsertOne(backCtx, insert)

	if mdbErr != nil {
		return dbController.NewDBError(mdbErr.Error())
	}

	return nil
}

func (mdbc *MongoDbController) AddInfoLog(log *logging.InfoLogData) error {
	collection, backCtx, cancel := mdbc.getCollection(LOGGING_COLLECTION)
	defer cancel()

	insert := bson.M{
		"timestamp": primitive.Timestamp{T: uint32(log.Timestamp.Unix())},
		"type":      log.Type,
		"message":   log.Message,
	}

	_, mdbErr := collection.InsertOne(backCtx, insert)

	if mdbErr != nil {
		return dbController.NewDBError(mdbErr.Error())
	}

	return nil
}
