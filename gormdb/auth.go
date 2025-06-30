package gormdb

import (
	"github.com/gragorther/epigo/models"
	"gorm.io/gorm"
)

type AuthDB struct {
	DB *gorm.DB
}

func (a *AuthDB) CheckUserAuthorizationForGroup(groupIDs []uint, userID uint) (bool, error) {
	var authorizedGroups int64
	if err := a.DB.Model(&models.Group{}).Where("user_id = ?", userID).
		Where("id IN ?", groupIDs).
		Count(&authorizedGroups).Error; err != nil {

		return false, err
	}
	if int(authorizedGroups) != len(groupIDs) {
		return false, nil
	}
	return true, nil
}
func (a *AuthDB) CheckUserAuthorizationForLastMessage(messageID uint, userID uint) (bool, error) {
	var authorizedCount int64
	res := a.DB.Model(&models.LastMessage{}).Where("id = ?", messageID).Where("user_id = ?", userID).Count(&authorizedCount)
	if authorizedCount != 1 {
		return false, res.Error
	}
	return true, res.Error
}
