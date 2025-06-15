package models

import "gorm.io/gorm"

type LastMessage struct {
	gorm.Model
	Title        string        `json:"title"`
	Content      string        `json:"content"` // markdown
	UserID       int           `json:"userID"`  // every LastMessage belongs to a user (this is the foreign key)
	User         User          // define relationship between LastMessage and User
	SendToGroups []SendToGroup `json:"sendToGroups"`
}

// groups of people the user can send messages to
type SendToGroup struct {
	gorm.Model
	Recipients []string `json:"recipients"`
	UserID     int      `json:"userID"` // defines the user that created this group
	User       User
}
