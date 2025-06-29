package db

import (
	"time"

	"github.com/gragorther/epigo/apperrors"
	"github.com/gragorther/epigo/models"
	"gorm.io/gorm"
)

type groupDB struct {
	DB *gorm.DB
}

func (g *groupDB) DeleteGroupByID(id uint) error {
	err := g.DB.Transaction(func(tx *gorm.DB) error {
		res := tx.Delete(&models.Group{}, id)
		return res.Error
	})
	return err
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

func (g *groupDB) FindGroupsAndRecipientEmailsByUserID(userID uint) ([]groupWithEmails, error) {
	var groups []listGroupsDTO
	res := g.DB.Model(&models.Group{}).Where("user_id = ?", userID).Preload("RecipientEmails").Find(&groups)
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

func (g *groupDB) CreateGroupAndRecipientEmails(group *models.Group, recipientEmails *[]models.RecipientEmail) error {

	newGroup := models.Group{
		UserID:          group.UserID,
		Name:            group.Name,
		Description:     group.Description,
		RecipientEmails: *recipientEmails,
	}
	err := g.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&newGroup).Error
		return err
	})
	return err
}

func (g *groupDB) UpdateGroup(group *models.Group, recipientEmails *[]models.RecipientEmail) error {
	err := g.DB.Transaction(func(tx *gorm.DB) error {
		output := g.DB.Updates(&group)
		if output.Error != nil {
			return output.Error
		}
		if output.RowsAffected < 1 {
			return apperrors.ErrNotFound
		}
		err := g.DB.Model(&group).Association("RecipientEmails").Replace(recipientEmails)

		return err
	})
	return err

}
