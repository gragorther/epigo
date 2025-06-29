package db

import (
	"time"

	"github.com/gragorther/epigo/models"
)

type Groups interface {
	DeleteGroupByID(id uint) error
	FindGroupsAndRecipientEmailsByUserID(userID uint) ([]groupWithEmails, error)
	CreateGroupAndRecipientEmails(group *models.Group, recipientEmails *[]models.RecipientEmail) error
	UpdateGroup(group *models.Group, recipientEmails *[]models.RecipientEmail) error
}
type groupWithEmails struct {
	ID              uint      `json:"id"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	RecipientEmails []string  `json:"recipientEmails"`
}
