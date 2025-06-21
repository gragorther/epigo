package initializers

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{TranslateError: true, FullSaveAssociations: true})

	if err != nil {
		return nil, err
	}
	return db, nil
}
