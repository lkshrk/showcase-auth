package models

import (
	"os"
	"testing"

	"harke.me/showcase-auth/database"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	user := User{
		Password: "secret",
	}

	err := user.HashPassword(user.Password)
	assert.NoError(t, err)

	os.Setenv("passwordHash", user.Password)
}

func TestCreateUserRecord(t *testing.T) {
	var userResult User

	err := database.InitDatabase()
	if err != nil {
		t.Error(err)
	}

	err = database.DB.AutoMigrate(&User{})
	assert.NoError(t, err)

	user := User{
		Username: "Test User",
		Role:     "RoleAdmin",
		Password: os.Getenv("passwordHash"),
	}

	err = user.CreateUser()
	assert.NoError(t, err)

	database.DB.Where("name = ?", user.Username).Find(&userResult)

	database.DB.Unscoped().Delete(&user)

	assert.Equal(t, "Test User", userResult.Username)
	assert.Equal(t, "RoleAdmin", userResult.Role)

}

func TestCheckPassword(t *testing.T) {
	hash := os.Getenv("passwordHash")

	user := User{
		Password: hash,
	}

	err := user.CheckPassword("secret")
	assert.NoError(t, err)
}
