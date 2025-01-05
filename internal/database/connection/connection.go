package connection

import (
	"fmt"

	"github.com/bigxxby/dream-test-task/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDB(cfg *config.Config) (*gorm.DB, error) {
	// Формируем строку подключения
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)

	// Подключаемся к базе данных через GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	return db, nil
}
