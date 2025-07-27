package gormdb

import (
	"time"

	"github.com/gragorther/epigo/apperrors"
	"github.com/gragorther/epigo/models"
	"github.com/gragorther/epigo/types"
	"gorm.io/gorm"
)

type GroupDB struct {
	db *gorm.DB
}

func (g *GormDB) DeleteGroupByID(id uint) error {
	err := g.db.Transaction(func(tx *gorm.DB) error {
		//res := tx.Delete(&models.Group{}, id)
		group := models.Group{ID: id}
		res := tx.Select("RecipientEmails").Delete(&group)
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

func (g *GormDB) FindGroupsAndRecipientEmailsByUserID(userID uint) ([]types.GroupWithEmails, error) {
	var groups []listGroupsDTO
	res := g.db.Model(&models.Group{}).Where("user_id = ?", userID).Preload("RecipientEmails").Find(&groups)
	if res.Error != nil {
		return nil, res.Error
	}
	var out []types.GroupWithEmails
	for _, g := range groups {
		emails := make([]string, len(g.RecipientEmails))
		for i, re := range g.RecipientEmails {
			emails[i] = re.Email
		}

		out = append(out, types.GroupWithEmails{
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

func (g *GormDB) CreateGroupAndRecipientEmails(group *models.Group, recipientEmails *[]models.RecipientEmail) error {

	newGroup := models.Group{
		UserID:          group.UserID,
		Name:            group.Name,
		Description:     group.Description,
		RecipientEmails: *recipientEmails,
	}
	err := g.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&newGroup).Error
		return err
	})
	return err
}

func (g *GormDB) UpdateGroup(group *models.Group, recipientEmails *[]models.RecipientEmail) error {
	err := g.db.Transaction(func(tx *gorm.DB) error {
		output := tx.Updates(&group)
		if output.Error != nil {
			return output.Error
		}
		if output.RowsAffected < 1 {
			return apperrors.ErrNotFound
		}
		err := tx.Model(&group).Association("RecipientEmails").Replace(recipientEmails)

		return err
	})
	return err

}
