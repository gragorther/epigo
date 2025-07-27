package gormdb

import (
	"gorm.io/gorm"
)

type GormDB struct {
	db *gorm.DB
}

func NewGormDB(gormdb *gorm.DB) *GormDB {
	return &GormDB{
		db: gormdb,
	}
}
