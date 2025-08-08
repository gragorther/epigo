package initializers

import (
	"context"
	"fmt"

	"github.com/gragorther/epigo/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresConfig struct {
	User     string
	Password string
	DBName   string
	Port     string
	TimeZone string
	Host     string
}

func ConnectDB(ctx context.Context, config PostgresConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=%v host=%v", config.User, config.Password, config.DBName, config.Port, config.TimeZone, config.Host)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true, FullSaveAssociations: true})

	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&models.User{}, &models.LastMessage{}, &models.Group{}, &models.Recipient{}, models.Profile{})
	return db, nil
}
