package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"harke.me/showcase-auth/pkg/repository"
	"harke.me/showcase-auth/pkg/repository/models"
)

func TestCreateUser(t *testing.T) {

	t.Run("create user", func(t *testing.T) {

		db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
		if err != nil {
			return
		}

		storage := repository.NewStorage(db)

		err = storage.RunMigrations()
		assert.Nil(t, err)

	})

	t.Run("unique constraint violated", func(t *testing.T) {

		storage, err := setupStorage()
		if err != nil {
			return
		}

		user := models.User{
			Username: "user",
			Password: "password",
			Role:     "role",
		}

		err = storage.CreateUser(user)
		assert.Nil(t, err)

		err = storage.CreateUser(user)
		assert.Error(t, err, "UNIQUE constraint failed: users.username")

	})

}

func TestFindUser(t *testing.T) {

	storage, err := setupStorage()
	if err != nil {
		return
	}

	t.Run("user not found", func(t *testing.T) {

		_, err = storage.FindUser("peter")
		assert.Error(t, err)

	})

	t.Run("user found", func(t *testing.T) {

		user := models.User{
			Username: "user",
			Password: "password",
			Role:     "role",
		}

		storage.CreateUser(user)

		actual, err := storage.FindUser(user.Username)
		if err != nil {
			return
		}

		assert.Equal(t, user.Username, actual.Username)
		assert.Equal(t, user.Password, actual.Password)
		assert.Equal(t, user.Role, actual.Role)

	})
}

func setupStorage() (repository.Storage, error) {

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	storage := repository.NewStorage(db)

	err = storage.RunMigrations()
	if err != nil {
		return nil, err
	}

	return storage, nil

}
