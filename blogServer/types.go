package blogServer

import (
	"time"

	"methompson.com/blog-microservice/blogServer/dbController"
)

type AuthorizationHeader struct {
	Token string `header:"authorization" binding:"required"`
}

type AddBlogBody struct {
	Title          string    `json:"title" binding:"required"`
	Slug           string    `json:"slug" binding:"required"`
	Body           string    `json:"body" binding:"required"`
	Tags           *[]string `json:"tags"`
	AuthorId       string    `json:"authorId" binding:"required"`
	DateAdded      int       `json:"dateAdded" binding:"required"`
	UpdateAuthorId *string   `json:"updateAuthorId"`
	DateUpdated    *int      `json:"dateUpdated"`
}

func (abb *AddBlogBody) GetBlogDocument() *dbController.AddBlogDocument {
	dateAdded := time.Unix(int64(abb.DateAdded), 0)

	var dateUpdated *time.Time
	if abb.DateUpdated != nil {
		t := time.Unix(int64(*abb.DateUpdated), 0)
		dateUpdated = &t
	}

	doc := dbController.AddBlogDocument{
		Title:          abb.Title,
		Slug:           abb.Slug,
		Body:           abb.Body,
		Tags:           abb.Tags,
		AuthorId:       abb.AuthorId,
		DateAdded:      dateAdded,
		UpdateAuthorId: abb.UpdateAuthorId,
		DateUpdated:    dateUpdated,
	}

	return &doc
}

type EditBlogBody struct {
	Id             string    `json:"id" binding:"required"`
	Title          *string   `json:"title"`
	Slug           *string   `json:"slug"`
	Body           *string   `json:"body"`
	Tags           *[]string `json:"tags"`
	AuthorId       *string   `json:"authorId"`
	DateAdded      *int      `json:"dateAdded"`
	UpdateAuthorId *string   `json:"updateAuthorId"`
	DateUpdated    *int      `json:"dateUpdated"`
}

func (ebb *EditBlogBody) GetBlogDocument() *dbController.EditBlogDocument {
	var dateAdded *time.Time
	var dateUpdated *time.Time

	if ebb.DateAdded != nil {
		t := time.Unix(int64(*ebb.DateAdded), 0)
		dateAdded = &t
	}

	if ebb.DateUpdated != nil {
		t := time.Unix(int64(*ebb.DateUpdated), 0)
		dateUpdated = &t
	}

	doc := dbController.EditBlogDocument{
		Id:             ebb.Id,
		Title:          ebb.Title,
		Slug:           ebb.Slug,
		Body:           ebb.Body,
		Tags:           ebb.Tags,
		AuthorId:       ebb.AuthorId,
		DateAdded:      dateAdded,
		UpdateAuthorId: ebb.UpdateAuthorId,
		DateUpdated:    dateUpdated,
	}

	return &doc
}

type DeleteBlogBody struct {
	Id string `json:"id" binding:"required"`
}

func (dbb *DeleteBlogBody) GetBlogDocument() *dbController.DeleteBlogDocument {
	doc := dbController.DeleteBlogDocument{
		Id: dbb.Id,
	}

	return &doc
}
