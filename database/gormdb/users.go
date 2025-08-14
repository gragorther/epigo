package gormdb

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/gragorther/epigo/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (g *GormDB) UpdateUserInterval(userID uint, cron string) error {
	res := g.db.Model(&models.User{}).Where("id = ?", userID).UpdateColumn("cron", cron)
	return res.Error
}

type UserInterval struct {
	ID        uint   `gorm:"primarykey"`
	Email     string `json:"email" gorm:"unique"`
	EmailCron string `json:"emailCron"`
}

func (g *GormDB) GetUserIntervals() ([]UserInterval, error) {
	var intervals []UserInterval
	res := g.db.Model(&models.User{}).Find(&intervals)
	return intervals, res.Error
}

// true if user exists, false if they don't exist
func (g *GormDB) CheckIfUserExistsByUsernameAndEmail(username string, email string) (bool, error) {
	var foundUsers int64
	res := g.db.Model(&models.User{}).
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
func (g *GormDB) CheckIfUserExistsByUsername(username string) (bool, error) {
	var userFound int64

	res := g.db.Model(&models.User{}).Where("username=?", username).Count(&userFound)
	if res.Error != nil {
		return false, res.Error
	}

	if userFound == 0 {
		return false, nil
	}
	return true, nil
}

func (g *GormDB) CreateUser(user *models.User) error {
	res := g.db.Create(user)
	return res.Error
}
func (g *GormDB) GetUserByUsername(username string) (*models.User, error) {
	var userFound models.User
	res := g.db.Model(&models.User{}).Where("username = ?", username).Find(&userFound)
	return &userFound, res.Error
}

func (g *GormDB) CheckIfUserExistsByID(ctx context.Context, ID uint) (bool, error) {

	userFound, err := gorm.G[models.User](g.db).Where("id = ?", ID).Count(ctx, "id")

	if err != nil {
		return false, err
	}

	if userFound == 0 {
		return false, nil
	}
	return true, nil
}
func (g *GormDB) GetUserByID(ID uint) (*models.User, error) {
	var user models.User
	res := g.db.Model(&models.User{ID: ID}).Find(&user)
	return &user, res.Error
}

func (g *GormDB) DeleteUser(ctx context.Context, ID uint) error {
	rowsaffected, err := gorm.G[models.User](g.db).Where("id = ?", ID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	if rowsaffected < 1 {
		return ErrNoRowsAffected
	}
	return nil

}

func (g *GormDB) EditUser(ctx context.Context, user models.User) error {
	rowsAffected, err := gorm.G[models.User](g.db).Where("id = ?", user.ID).Updates(ctx, user)
	if err != nil {
		return err
	}
	if rowsAffected < 1 {
		return ErrNoRowsAffected
	}
	return nil

}

func (g *GormDB) DeleteUserAndAllAssociations(ID uint) error {
	res := g.db.Select(clause.Associations).Delete(&models.User{ID: ID})
	return res.Error
}
func (g *GormDB) CreateProfile(ctx context.Context, newProfile *models.Profile) error {
	err := gorm.G[models.Profile](g.db).Create(ctx, newProfile)
	return err
}

var ErrNoRowsAffected error = errors.New("no rows affected")

func (g *GormDB) UpdateProfile(ctx context.Context, profile models.Profile) error {
	rowsAffected, err := gorm.G[models.Profile](g.db).Where("id = ?", profile.ID).Updates(ctx, profile)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNoRowsAffected
	}
	return nil
}
