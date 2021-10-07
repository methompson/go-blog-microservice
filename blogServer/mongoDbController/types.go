package mongoDbController

import (
	"time"

	"methompson.com/blog-microservice/blogServer/dbController"
)

type UserDocResult struct {
	Id    string `bson:"_id"`
	UID   string `bson:"uid"`
	Name  string `bson:"name"`
	Role  string `bson:"role"`
	Email string `bson:"email"`
}

func (udr *UserDocResult) GetUserDataDoc() *dbController.UserDataDocument {
	doc := dbController.UserDataDocument{
		Id:    udr.Id,
		UID:   udr.UID,
		Name:  udr.Name,
		Role:  udr.Role,
		Email: udr.Email,
	}

	return &doc
}

type BlogDocResult struct {
	Id             string          `bson:"_id"`
	Title          string          `bson:"title"`
	Slug           string          `bson:"slug"`
	Body           string          `bson:"body"`
	Tags           []string        `bson:"tags"`
	Author         []UserDocResult `bson:"author"`
	AuthorId       string          `bson:"authorId"`
	DateAdded      time.Time       `bson:"dateAdded"`
	UpdateAuthor   []UserDocResult `bson:"updateAuthor"`
	UpdateAuthorId string          `bson:"updateAuthorId"`
	DateUpdated    time.Time       `bson:"dateUpdated"`
}

func (bdr *BlogDocResult) GetBlogDocument() *dbController.BlogDocument {
	var author string = ""
	if len(bdr.Author) > 0 {
		author = bdr.Author[0].Name
	}

	var updateAuthor string = ""
	if len(bdr.UpdateAuthor) > 0 {
		updateAuthor = bdr.Author[0].Name
	}

	doc := dbController.BlogDocument{
		Id:             bdr.Id,
		Title:          bdr.Title,
		Slug:           bdr.Slug,
		Body:           bdr.Body,
		Tags:           bdr.Tags,
		Author:         author,
		AuthorId:       bdr.AuthorId,
		DateAdded:      bdr.DateAdded,
		UpdateAuthor:   updateAuthor,
		UpdateAuthorId: bdr.UpdateAuthorId,
		DateUpdated:    bdr.DateUpdated,
	}

	return &doc
}
