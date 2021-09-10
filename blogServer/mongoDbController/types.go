package mongoDbController

import (
	"time"

	"methompson.com/blog-microservice/blogServer/dbController"
)

type UserDocResult struct {
	Id     string `bson:"_id"`
	UserId string `bson:"userId"`
	Name   string `bson:"name"`
}

func (udr *UserDocResult) GetUserDataDoc() *dbController.UserDataDocument {
	doc := dbController.UserDataDocument{
		Id:     udr.Id,
		UserId: udr.UserId,
		Name:   udr.Name,
	}

	return &doc
}

type BlogDocResult struct {
	Id             string    `bson:"_id"`
	Title          string    `bson:"title"`
	Slug           string    `bson:"slug"`
	Body           string    `bson:"body"`
	Tags           []string  `bson:"tags"`
	AuthorId       string    `bson:"authorId"`
	DateAdded      time.Time `bson:"dateAdded"`
	UpdateAuthorId string    `bson:"updateAuthorId"`
	DateUpdated    time.Time `bson:"dateUpdated"`
}

func (bdr *BlogDocResult) GetBlogDocument() *dbController.BlogDocument {
	doc := dbController.BlogDocument{
		Id:             bdr.Id,
		Title:          bdr.Title,
		Slug:           bdr.Slug,
		Body:           bdr.Body,
		Tags:           bdr.Tags,
		AuthorId:       bdr.AuthorId,
		DateAdded:      bdr.DateAdded,
		UpdateAuthorId: bdr.UpdateAuthorId,
		DateUpdated:    bdr.DateUpdated,
	}

	return &doc
}
