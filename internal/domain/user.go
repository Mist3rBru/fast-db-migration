package domain

import (
	"github.com/google/uuid"
)

type User struct {
	Id       string `gorm:"primarykey" bson:"_id"`
	Name     string
	Email    string
	Password string
}

func NewUser(name string, email string, password string) User {
	return User{
		Id:       uuid.NewString(),
		Name:     name,
		Email:    email,
		Password: password,
	}
}
