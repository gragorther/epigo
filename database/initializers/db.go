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

	return db, err
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{}, &models.LastMessage{}, &models.Group{}, &models.Recipient{}, &models.Profile{})
}
