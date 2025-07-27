package gormdb

import (
	"gorm.io/gorm"
)

type GormDBs struct {
	Auth    *AuthDB
	Group   *GroupDB
	User    *UserDB
	Message *MessageDB
}

func NewGormDB(gormdb *gorm.DB) *GormDBs {
	return &GormDBs{
		Auth:    &AuthDB{db: gormdb},
		Message: &MessageDB{db: gormdb},
		User:    &UserDB{db: gormdb},
		Group:   &GroupDB{db: gormdb},
	}
}
