package gormdb

import (
	"gorm.io/gorm"
)

type gormDBs struct {
	Auth    *AuthDB
	Group   *GroupDB
	User    *UserDB
	Message *MessageDB
}

func NewGormDB(gormdb *gorm.DB) *gormDBs {
	return &gormDBs{
		Auth:    &AuthDB{db: gormdb},
		Message: &MessageDB{db: gormdb},
		User:    &UserDB{db: gormdb},
		Group:   &GroupDB{db: gormdb},
	}
}
