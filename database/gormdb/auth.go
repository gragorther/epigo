package gormdb

import (
	"github.com/gragorther/epigo/models"
)

func (a *GormDB) CheckUserAuthorizationForGroup(groupIDs []uint, userID uint) (bool, error) {
	var authorizedGroups int64
	if err := a.db.Model(&models.Group{}).Where("user_id = ?", userID).
		Where("id IN ?", groupIDs).
		Count(&authorizedGroups).Error; err != nil {

		return false, err
	}
	if int(authorizedGroups) != len(groupIDs) {
		return false, nil
	}
	return true, nil
}
func (a *GormDB) CheckUserAuthorizationForLastMessage(messageID uint, userID uint) (bool, error) {
	var authorizedCount int64
	res := a.db.Model(&models.LastMessage{}).Where("id = ?", messageID).Where("user_id = ?", userID).Count(&authorizedCount)
	if authorizedCount != 1 {
		return false, res.Error
	}
	return true, res.Error
}
