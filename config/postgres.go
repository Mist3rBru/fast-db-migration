package config

import (
	"fast-db-migration/internal/domain"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initPostgres() (*gorm.DB, error) {
	logger := NewLogger("postgres")

	dns := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASS"), os.Getenv("POSTGRES_DBNAME"), os.Getenv("POSTGRES_PORT"))
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		logger.Error("opening error: ", err)
		return nil, err
	}
	// createDBSQL := fmt.Sprintf("CREATE DATABASE %s;", os.Getenv("POSTGRES_DBNAME"))
	// if err := db.Exec(createDBSQL).Error; err != nil {
	// 	logger.Error("create database error: ", err)
	// 	return nil, err
	// }
	err = db.AutoMigrate(&domain.User{})
	if err != nil {
		logger.Error("automigration error: ", err)
		return nil, err
	}

	logger.Info("connected to postgres")

	return db, nil
}
