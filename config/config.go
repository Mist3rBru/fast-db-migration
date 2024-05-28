package config

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var (
	Postgres *gorm.DB
	Mongo    *mongo.Database
)

func Init() error {
	var err error

	Postgres, err = initPostgres()
	if err != nil {
		return fmt.Errorf("error initializing postgres: %v", err)
	}

	Mongo, err = initMongo()
	if err != nil {
		return fmt.Errorf("error initializing mongo: %v", err)
	}

	return nil
}
