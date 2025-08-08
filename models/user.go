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
	Email        string         `json:"email" gorm:"unique"`
	PasswordHash string         `json:"-"`
	LastLogin    *time.Time     `json:"lastLogin"`
	Groups       []Group        `json:"groups"`
	EmailCron    *string        `json:"emailCron"`
	IsAdmin      bool           `json:"isAdmin"`
	IsVerified   bool           `json:"isVerified"`
	Profile      *Profile       `json:"profile"`
}
type Profile struct {
	gorm.Model
	UserID uint
	Name   *string `json:"name"`
}
