package models

import (
	"harke.me/showcase-auth/database"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username"`
	Role     string `json:"role"`
	Password string `json:"password"`
}

func (user *User) CreateUser() error {
	result := database.DB.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}

	user.Password = string(bytes)

	return nil
}

func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}

	return nil
}
