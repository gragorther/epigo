package models

import (
	"time"

	"gorm.io/gorm"
)

type LastMessage struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Title     string         `json:"title"`
	Content   string         `json:"content"` // markdown
	UserID    uint           `json:"userID"`
	Groups    []Group        `json:"groups" gorm:"many2many:group_last_messages;"`
}
