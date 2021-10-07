package dbController

import (
	"time"
)

type UserDataDocument struct {
	Id    string
	UID   string
	Name  string
	Role  string
	Email string
}

type AddBlogDocument struct {
	Title          string
	Slug           string
	Body           string
	Tags           *[]string
	AuthorId       string
	DateAdded      time.Time
	UpdateAuthorId *string
	DateUpdated    *time.Time
}

type BlogDocument struct {
	Id             string
	Title          string
	Slug           string
	Body           string
	Tags           []string
	Author         string
	AuthorId       string
	DateAdded      time.Time
	UpdateAuthor   string
	UpdateAuthorId string
	DateUpdated    time.Time
}

func (bd *BlogDocument) GetMap() *map[string]interface{} {
	m := make(map[string]interface{})

	m["id"] = bd.Id

	m["title"] = bd.Title
	m["slug"] = bd.Slug
	m["body"] = bd.Body
	m["author"] = bd.Author
	m["authorId"] = bd.AuthorId
	m["dateAdded"] = bd.DateAdded.Unix()
	m["updateAuthor"] = bd.UpdateAuthor
	m["updateAuthorId"] = bd.UpdateAuthorId
	m["dateUpdated"] = bd.DateUpdated.Unix()

	if bd.Tags != nil {
		m["tags"] = bd.Tags
	} else {
		m["tags"] = make([]string, 0)
	}

	return &m
}

type EditBlogDocument struct {
	Id             string
	Title          *string
	Slug           *string
	Body           *string
	Tags           *[]string
	AuthorId       *string
	DateAdded      *time.Time
	UpdateAuthorId *string
	DateUpdated    *time.Time
}

type DeleteBlogDocument struct {
	Id string
}
