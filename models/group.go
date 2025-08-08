package models

import (
	"time"

	"gorm.io/gorm"
)

type Group struct {
	ID           uint           `json:"ID,omitzero" gorm:"primarykey"`
	UserID       uint           `json:"-"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
	Name         string         `json:"name"`
	Description  *string        `json:"description"`
	Recipients   []Recipient    `json:"recipients"`
	LastMessages []LastMessage  `json:"lastMessages" gorm:"many2many:group_last_messages;"`
}

type APIRecipient struct {
	Email string `json:"email"`
}

type Recipient struct {
	gorm.Model
	APIRecipient
	GroupID uint `json:"groupID"` // group of the email
}
