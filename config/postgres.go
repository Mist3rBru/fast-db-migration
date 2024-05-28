package config

import (
	"fast-db-migration/internal/domain"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func initPostgres() (*gorm.DB, error) {
	_logger := NewLogger("postgres")

	dns := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASS"), os.Getenv("POSTGRES_DBNAME"), os.Getenv("POSTGRES_PORT"))
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		_logger.Error("opening error: ", err)
		return nil, err
	}

	err = db.AutoMigrate(&domain.User{})
	if err != nil {
		_logger.Error("automigration error: ", err)
		return nil, err
	}

	_logger.Info("connected to postgres")

	return db, nil
}
