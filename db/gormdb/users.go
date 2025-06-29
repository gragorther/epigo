package db

import (
	"log"

	"github.com/gragorther/epigo/db"
	"github.com/gragorther/epigo/models"
	"gorm.io/gorm"
)

type userDB struct {
	DB *gorm.DB
}

func (u *userDB) UpdateUserInterval(userID uint, cron string) error {
	res := u.DB.Model(&models.User{}).Where("id = ?", userID).Update("email_cron", cron)
	return res.Error
}

type userIntervalsOutput struct {
	ID        uint   `gorm:"primarykey"`
	Email     string `json:"email" gorm:"unique"`
	EmailCron string `json:"emailCron"`
}

func (u *userDB) GetUserIntervals() ([]db.UserIntervalsOutput, error) {
	var intervals []db.UserIntervalsOutput
	res := u.DB.Model(&models.User{}).Find(&intervals)
	return intervals, res.Error
}

// true if user exists, false if they don't exist
func (u *userDB) CheckIfUserExistsByUsernameAndEmail(username string, email string) (bool, error) {
	var foundUsers int64
	res := u.DB.Model(&models.User{}).
		Where("username = ? OR email = ?", username, email).Count(&foundUsers)

	if res.Error != nil {
		log.Printf("Couldn't check if user exists: %v", res.Error)
		return true, res.Error
	}

	if foundUsers > 0 {
		log.Printf("User %v already exists", username)
		return true, nil
	}
	return false, nil
}
func (u *userDB) CheckIfUserExistsByUsername(username string) (bool, error) {
	var userFound int64

	res := u.DB.Where("username=?", username).Count(&userFound)
	if res.Error != nil {
		return false, res.Error
	}

	if userFound == 0 {
		return false, nil
	}
	return true, nil
}

func (u *userDB) CreateUser(user *models.User) error {
	res := u.DB.Create(user)
	return res.Error
}
func (u *userDB) GetUserByUsername(username string) (*models.User, error) {
	var userFound models.User
	res := u.DB.Where("username = ?", username).Find(&userFound)
	return &userFound, res.Error
}
func (u *userDB) SaveUserData(user *models.User) error {
	res := u.DB.Save(user)
	return res.Error
}
