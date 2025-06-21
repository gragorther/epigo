package db

import (
	"time"

	"github.com/gragorther/epigo/models"
)

func (h *DBHandler) DeleteGroupByID(id uint) error {
	res := h.DB.Delete(&models.Group{}, id)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

type listGroupsOutput struct {
	ID          uint      `gorm:"primarykey"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

func (h *DBHandler) FindGroupsByUserID(userID uint) ([]listGroupsOutput, error) {
	var groups []listGroupsOutput
	res := h.DB.Model(&models.Group{}).Where("user_id = ?", userID).Find(&groups)
	if res.Error != nil {
		return nil, res.Error
	}
	return groups, nil
}

func (h *DBHandler) CreateGroupAndRecipientEmails(group *models.Group, recipientEmails *[]models.RecipientEmail) {
	h.DB.Create(&group).Association("RecipientEmails").Append(&recipientEmails)
}
