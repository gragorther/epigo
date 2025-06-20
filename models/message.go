package models

import "gorm.io/gorm"

type LastMessage struct {
	gorm.Model
	Title   string  `json:"title"`
	Content string  `json:"content"` // markdown
	UserID  uint    `json:"userID"`
	Groups  []Group `json:"groups" gorm:"many2many:group_last_messages;"`
}
