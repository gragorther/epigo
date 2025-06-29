package db

import (
	"github.com/gragorther/epigo/models"
	"github.com/gragorther/epigo/types"
)

type Users interface {
	UpdateUserInterval(userID uint, cron string) error
	GetUserIntervals() ([]types.UserIntervalsOutput, error)
	CheckIfUserExistsByUsernameAndEmail(username string, email string) (bool, error)
	CheckIfUserExistsByUsername(username string) (bool, error)
	CreateUser(*models.User) error
	GetUserByUsername(username string) (*models.User, error)
	SaveUserData(*models.User) error
}
