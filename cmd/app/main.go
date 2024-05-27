package main

import (
	"context"
	"fast-db-migration/config"
	"fast-db-migration/internal/domain"
	"fmt"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

var (
	logger *config.Logger
	mg     *mongo.Database
	pg     *gorm.DB
)

func generateRandomPassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func mockData() error {
	count, err := mg.Collection("users").CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		logger.Error("count users err: ", err)

		var users []interface{}
		count = int64(100000)

		for i := 0; i < int(count); i++ {
			user := domain.NewUser(fmt.Sprintf("User%d", i), fmt.Sprintf("user%d@example.com", i), generateRandomPassword(10))
			users = append(users, user)
		}

		// Insert the users into the collection
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		_, err := mg.Collection("users").InsertMany(ctx, users)
		if err != nil {
			return err
		}

		logger.Info("mocked data on mongo")
	}

	logger.Infof("data mocked %d", count)

	return nil
}

func getPagedUsers(page int) ([]domain.User, error) {
	pageLength := int64(4000)
	skip := int64(page) * pageLength
	findOptions := options.Find().SetSkip(skip).SetLimit(pageLength)
	cursor, err := mg.Collection("users").Find(context.TODO(), bson.D{}, findOptions)
	if err != nil {
		return nil, err
	}

	var users []domain.User
	if err := cursor.All(context.TODO(), &users); err != nil {
		return nil, err
	}

	return users, nil
}

func main() {
	logger = config.NewLogger("main")

	err := config.Init()
	if err != nil {
		logger.Error("config init error: ", err)
		return
	}
	mg = config.GetMongo()
	pg = config.GetPostgres()

	if err := mockData(); err != nil {
		logger.Error("mock data error: ", err)
		return
	}

	users, err := getPagedUsers(0)
	if err != nil {
		logger.Error("page load error: ", err)
		return
	}

	logger.Info(users)
}
