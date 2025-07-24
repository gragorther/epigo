package gormdb

import "gorm.io/gorm"

type gormDBs struct {
	AuthDB
	GroupDB
	UserDB
	MessageDB
}

func NewGormDB(db *gorm.DB) *gormDBs {
	return &gormDBs{
		AuthDB{db: db},
		GroupDB{db: db},
		UserDB{db: db},
		MessageDB{db: db},
	}
}
