package gormdb

import (
	"github.com/gragorther/epigo/db"
	"gorm.io/gorm"
)

type gormDBs struct {
	AuthDB
	GroupDB
	UserDB
	MessageDB
}

func NewGormDB(gormdb *gorm.DB) *db.DBHandler {
	return &db.DBHandler{
		Auth:     &AuthDB{db: gormdb},
		Messages: &MessageDB{db: gormdb},
		Users:    &UserDB{db: gormdb},
		Groups:   &GroupDB{db: gormdb},
	}
}
