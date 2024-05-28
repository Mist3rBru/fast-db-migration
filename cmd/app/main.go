package main

import (
	"context"
	"fast-db-migration/config"
	"fast-db-migration/internal/domain"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	logger     *config.Logger
	pageLength = int64(4000)
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
	count := 1000000
	countDocs, err := config.Mongo.Collection("users").CountDocuments(context.TODO(), bson.D{})
	if err != nil || int(countDocs) < count {
		logger.Error("count users error: ", err)

		var users []interface{}

		for i := range count {
			user := domain.NewUser(fmt.Sprintf("User%d", i), fmt.Sprintf("user%d@example.com", i), generateRandomPassword(10))
			users = append(users, user)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		_, err := config.Mongo.Collection("users").InsertMany(ctx, users)
		if err != nil {
			logger.Error("insert users error: ", err)
			return err
		}

		logger.Info("mocked data on mongo")
	}

	logger.Infof("data mocked %v", count)

	return nil
}

func getPagedUsers(page int) ([]domain.User, error) {
	skip := int64(page) * pageLength
	findOptions := options.Find().SetSkip(skip).SetLimit(pageLength)
	cursor, err := config.Mongo.Collection("users").Find(context.TODO(), bson.D{}, findOptions)
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

	logger.Info("starting app")

	err := config.Init()
	if err != nil {
		logger.Error("config init error: ", err)
		return
	}

	logger.Info("mocking mongo users...")
	if err := mockData(); err != nil {
		logger.Error("mock data error: ", err)
		return
	}

	logger.Info("deleting postgres users...")
	config.Postgres.Exec("DELETE FROM users")

	mongoUsers, err := config.Mongo.Collection("users").CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		logger.Error("mock data error: ", err)
		return
	}
	pages := int((mongoUsers / pageLength) + 1)

	wg := sync.WaitGroup{}
	for page := range pages {
		wg.Add(1)
		go func(page int) {
			defer wg.Done()
			
			users, err := getPagedUsers(page)
			if err != nil {
				logger.Error("page load error: ", err)
				return
			}
			config.Postgres.Create(&users)
			logger.Infof("migraded %d of %d pages", page, pages)
		}(page)
	}

	wg.Wait()

	var users []domain.User
	config.Postgres.Table("users").Find(&users)

	logger.Infof("migraded %d users", len(users))
}
