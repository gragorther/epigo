package db

import (
	"github.com/gragorther/epigo/models"
	"github.com/gragorther/epigo/types"
)

type Groups interface {
	DeleteGroupByID(id uint) error
	FindGroupsAndRecipientEmailsByUserID(userID uint) ([]types.GroupWithEmails, error)
	CreateGroupAndRecipientEmails(group *models.Group, recipientEmails *[]models.RecipientEmail) error
	UpdateGroup(group *models.Group, recipientEmails *[]models.RecipientEmail) error
}
