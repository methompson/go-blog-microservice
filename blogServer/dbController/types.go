package dbController

import (
	"time"
)

type UserDataDocument struct {
	Id   string
	Name string
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

type NonceDocument struct {
	NonceHash     string `bson:"hash"`
	RemoteAddress string `bson:"remoteAddress"`
	Time          int    `bson:"time"`
}

type FullUserDocument struct {
	Id           string
	Username     string
	Email        string
	Enabled      bool
	Admin        bool
	PasswordHash string
}

func (fud *FullUserDocument) GetUserDocument() UserDocument {
	return UserDocument{
		Id:       fud.Id,
		Username: fud.Username,
		Email:    fud.Email,
		Enabled:  fud.Enabled,
		Admin:    fud.Admin,
	}
}

type UserDocument struct {
	Id       string
	Username string
	Email    string
	Enabled  bool
	Admin    bool
}

type EditUserDocument struct {
	Id       string
	Username *string
	Email    *string
	Enabled  *bool
	Admin    *bool
}
