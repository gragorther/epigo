package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint `gorm:"primarykey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	Username     string         `json:"username" gorm:"unique"`
	Name         *string        `json:"name"`
	Email        string         `json:"email" gorm:"unique"`
	PasswordHash string         `json:"-"`
	//	Country      string //should probably be a foreign key of another table
	LastLogin *time.Time `json:"lastLogin"`
	Groups    Group
	EmailCron *string `json:"emailCron"`
	IsAdmin   bool    `json:"isAdmin"`
}
