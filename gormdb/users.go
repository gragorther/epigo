package gormdb

import (
	"log"

	"github.com/gragorther/epigo/models"
	"github.com/gragorther/epigo/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserDB struct {
	DB *gorm.DB
}

func (u *UserDB) UpdateUserInterval(userID uint, cron string) error {
	res := u.DB.Model(&models.User{}).Where("id = ?", userID).Update("email_cron", cron)
	return res.Error
}

func (u *UserDB) GetUserIntervals() ([]types.UserIntervalsOutput, error) {
	var intervals []types.UserIntervalsOutput
	res := u.DB.Model(&models.User{}).Find(&intervals)
	return intervals, res.Error
}

// true if user exists, false if they don't exist
func (u *UserDB) CheckIfUserExistsByUsernameAndEmail(username string, email string) (bool, error) {
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
func (u *UserDB) CheckIfUserExistsByUsername(username string) (bool, error) {
	var userFound int64

	res := u.DB.Model(&models.User{}).Where("username=?", username).Count(&userFound)
	if res.Error != nil {
		return false, res.Error
	}

	if userFound == 0 {
		return false, nil
	}
	return true, nil
}

func (u *UserDB) CreateUser(user *models.User) error {
	res := u.DB.Create(user)
	return res.Error
}
func (u *UserDB) GetUserByUsername(username string) (*models.User, error) {
	var userFound models.User
	res := u.DB.Model(&models.User{}).Where("username = ?", username).Find(&userFound)
	return &userFound, res.Error
}

func (u *UserDB) CheckIfUserExistsByID(ID uint) (bool, error) {
	var userFound int64

	res := u.DB.Model(&models.User{}).Where("id=?", ID).Count(&userFound)
	if res.Error != nil {
		return false, res.Error
	}

	if userFound == 0 {
		return false, nil
	}
	return true, nil
}
func (u *UserDB) GetUserByID(ID uint) (*models.User, error) {
	var user models.User
	res := u.DB.Model(&models.User{ID: ID}).Find(&user)
	return &user, res.Error
}
func (u *UserDB) SaveUserData(user *models.User) error {
	res := u.DB.Save(user)
	return res.Error
}

func (u *UserDB) DeleteUser(ID uint) error {
	res := u.DB.Delete(&models.User{}, ID)
	return res.Error
}

func (u *UserDB) EditUser(user *models.User) error {
	res := u.DB.Model(&models.User{ID: user.ID}).Updates(user)
	return res.Error
}

func (u *UserDB) DeleteUserAndAllAssociations(ID uint) error {
	res := u.DB.Select(clause.Associations).Delete(&models.User{ID: ID})
	return res.Error
}
