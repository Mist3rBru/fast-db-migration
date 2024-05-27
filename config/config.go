package config

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var (
	pg *gorm.DB
	mg *mongo.Database
)

func Init() error {
	var err error

	pg, err = initPostgres()
	if err != nil {
		return fmt.Errorf("error initializing postgres: %v", err)
	}

	mg, err = initMongo()
	if err != nil {
		return fmt.Errorf("error initializing mongo: %v", err)
	}

	return nil
}

func GetPostgres() *gorm.DB {
	return pg
}

func GetMongo() *mongo.Database {
	return mg
}
