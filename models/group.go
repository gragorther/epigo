package models

import (
	"time"

	"gorm.io/gorm"
)

type Group struct {
	ID              uint `gorm:"primarykey"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt   `gorm:"index"`
	UserID          uint             `json:"userID"`
	Name            string           `json:"name"`
	Description     string           `json:"description"`
	RecipientEmails []RecipientEmail `json:"recipientEmails"`
	LastMessages    []LastMessage    `json:"lastMessages" gorm:"many2many:group_last_messages;"`
}

type RecipientEmail struct {
	gorm.Model
	GroupID uint   `json:"groupID"` // group of the email
	Email   string `json:"email"`
}
