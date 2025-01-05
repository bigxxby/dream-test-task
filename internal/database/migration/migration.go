package migration

import (
	"github.com/bigxxby/dream-test-task/internal/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&models.User{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&models.ShortLink{})
	if err != nil {
		return err
	}
	return nil
}
