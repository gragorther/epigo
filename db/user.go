package db

import (
	"github.com/gragorther/epigo/models"
)

func (h *DBHandler) UpdateUserInterval(userID uint, cron string) error {
	res := h.DB.Model(&models.User{}).Where("id = ?", userID).Update("email_cron", cron)
	return res.Error
}

type userIntervalsOutput struct {
	ID        uint   `gorm:"primarykey"`
	Email     string `json:"email" gorm:"unique"`
	EmailCron string `json:"emailCron"`
}

func (h *DBHandler) GetUserIntervals() ([]userIntervalsOutput, error) {
	var intervals []userIntervalsOutput
	res := h.DB.Model(&models.User{}).Find(&intervals)
	return intervals, res.Error
}
