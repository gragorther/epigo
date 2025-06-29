package db

import "github.com/gragorther/epigo/models"

type Users interface {
	UpdateUserInterval(userID uint, cron string) error
	GetUserIntervals() ([]UserIntervalsOutput, error)
	CheckIfUserExistsByUsernameAndEmail(username string, email string) (bool, error)
	CheckIfUserExistsByUsername(username string) (bool, error)
	CreateUser(*models.User) error
	GetUserByUsername(username string) (*models.User, error)
	SaveUserData(*models.User) error
}
type UserIntervalsOutput struct {
	ID        uint   `gorm:"primarykey"`
	Email     string `json:"email" gorm:"unique"`
	EmailCron string `json:"emailCron"`
}
