package db

import (
	"github.com/gragorther/epigo/models"
	"gorm.io/gorm"
)

type DBHandler struct {
	DB *gorm.DB
}

func NewDBHandler(db *gorm.DB) *DBHandler {
	return &DBHandler{DB: db}
}

func (h *DBHandler) CheckUserAuthorizationForGroup(groupID uint, userID uint) (bool, error) {

	var authorizedGroup int64
	if err := h.DB.Model(&models.Group{}).
		Where("id = ?", groupID).
		Where("user_id = ?", userID).
		Count(&authorizedGroup).Error; err != nil {
		return false, err
	}
	return (authorizedGroup == 1), nil
}
