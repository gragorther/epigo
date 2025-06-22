package db

import (
	"time"

	"github.com/gragorther/epigo/apperrors"
	"github.com/gragorther/epigo/models"
)

func (h *DBHandler) DeleteGroupByID(id uint) error {
	res := h.DB.Delete(&models.Group{}, id)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

type recipientEmail struct {
	Email   string `json:"email"`
	GroupID uint
}

type listGroupsDTO struct {
	ID              uint             `gorm:"primarykey"`
	CreatedAt       time.Time        `json:"createdAt"`
	UpdatedAt       time.Time        `json:"updatedAt"`
	Name            string           `json:"name"`
	Description     string           `json:"description"`
	RecipientEmails []recipientEmail `json:"recipientEmails" gorm:"foreignKey:GroupID"`
}
type groupWithEmails struct {
	ID              uint      `json:"id"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	RecipientEmails []string  `json:"recipientEmails"`
}

func (h *DBHandler) FindGroupsAndRecipientEmailsByUserID(userID uint) ([]groupWithEmails, error) {
	var groups []listGroupsDTO
	res := h.DB.Model(&models.Group{ID: userID}).Preload("RecipientEmails").Find(&groups)
	if res.Error != nil {
		return nil, res.Error
	}
	var out []groupWithEmails
	for _, g := range groups {
		emails := make([]string, len(g.RecipientEmails))
		for i, re := range g.RecipientEmails {
			emails[i] = re.Email
		}

		out = append(out, groupWithEmails{
			ID:              g.ID,
			CreatedAt:       g.CreatedAt,
			UpdatedAt:       g.UpdatedAt,
			Name:            g.Name,
			Description:     g.Description,
			RecipientEmails: emails,
		})
	}
	return out, nil
}

func (h *DBHandler) CreateGroupAndRecipientEmails(group *models.Group, recipientEmails *[]models.RecipientEmail) error {

	newGroup := models.Group{
		UserID:          group.UserID,
		Name:            group.Name,
		Description:     group.Description,
		RecipientEmails: *recipientEmails,
	}
	err := h.DB.Create(&newGroup).Error

	if err != nil {
		return err
	}
	return nil
}

func (h *DBHandler) UpdateGroup(group *models.Group, recipientEmails *[]models.RecipientEmail) error {
	output := h.DB.Updates(&group)
	if output.Error != nil {
		return output.Error
	}
	if output.RowsAffected < 1 {
		return apperrors.ErrNotFound
	}
	h.DB.Model(&group).Association("RecipientEmails").Replace(recipientEmails)

	return nil
}
