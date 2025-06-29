package db

import (
	db "github.com/gragorther/epigo/db/gormdb"
	"gorm.io/gorm"
)

type DB interface {
	Auth
	Groups
	Messages
	Users
}

type DBHandler struct {
	Auth     Auth
	Messages Messages
	Users    Users
	Groups   Groups
}

// NewDBHandler returns a concrete DBHandler with all repos initialized
func NewDBHandler(conn *gorm.DB) *DBHandler {
	return &DBHandler{
		Users:    &db.UserDB{DB: conn},
		Auth:     &db.AuthDB{DB: conn},
		Groups:   &db.GroupDB{DB: conn},
		Messages: &db.MessageDB{DB: conn},
	}
}
