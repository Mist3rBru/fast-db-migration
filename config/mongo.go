package config

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initMongo() (*mongo.Database, error) {
	logger := NewLogger("mongo")

	clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logger.Error("connection error: ", err)
		return nil, err
	}
	logger.Info("connected to mongo")

	db := client.Database(os.Getenv("MONGODB_DBNAME"))
	if err := db.CreateCollection(context.TODO(), "users", &options.CreateCollectionOptions{}); err != nil {
		logger.Error("create users collection error: ", err)
	}

	return db, nil
}
