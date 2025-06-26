package db

import "github.com/gragorther/epigo/models"

func (h *DBHandler) UpdateUserInterval(userID uint, cron string) {
	h.DB.Model(&models.User{}).Where("id = ?", userID).Update("email_cron", cron)
}
