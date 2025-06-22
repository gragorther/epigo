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

func (h *DBHandler) CheckUserAuthorizationForGroup(groupIDs []uint, userID uint) (bool, error) {
	var authorizedGroups int64
	if err := h.DB.Model(&models.Group{}).Where("user_id = ?", userID).
		Where("id IN ?", groupIDs).
		Count(&authorizedGroups).Error; err != nil {

		return false, err
	}
	if int(authorizedGroups) != len(groupIDs) {
		return false, nil
	}
	return true, nil
}
func (h *DBHandler) CheckUserAuthorizationForLastMessage(messageID uint, userID uint) (bool, error) {
	var authorizedCount int64
	res := h.DB.Model(&models.LastMessage{}).Where("id = ?", messageID).Where("user_id = ?", userID).Count(&authorizedCount)
	if authorizedCount != 1 {
		return false, res.Error
	}
	return true, res.Error
}
