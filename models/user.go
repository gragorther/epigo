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
	Cron         *string        `json:"cron"`
	IsAdmin      bool           `json:"isAdmin"`
	Profile      *Profile       `json:"profile"`
	Alive        bool
	LastMessages []LastMessage
}
type Profile struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	UserID    uint
	Name      *string `json:"name" sql:"name"`
}
