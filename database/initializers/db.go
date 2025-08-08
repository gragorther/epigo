package initializers

import (
	"github.com/gragorther/epigo/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{TranslateError: true, FullSaveAssociations: true})

	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&models.User{}, &models.LastMessage{}, &models.Group{}, &models.Recipient{}, models.Profile{})
	return db, nil
}
