package initializers

import (
	"context"

	"github.com/gragorther/epigo/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(ctx context.Context, dsn string) (*gorm.DB, error) {

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true, FullSaveAssociations: true})

	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&models.User{}, &models.LastMessage{}, &models.Group{}, &models.Recipient{}, &models.Profile{})

	return db, err
}
