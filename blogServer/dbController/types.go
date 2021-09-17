package dbController

import (
	"time"
)

type UserDataDocument struct {
	Id     string
	UserId string
	Name   string
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
	AuthorId       string
	DateAdded      time.Time
	UpdateAuthorId string
	DateUpdated    time.Time
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
