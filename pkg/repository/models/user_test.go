package models_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"harke.me/showcase-auth/pkg/repository/models"
)

func TestCheckPassword(t *testing.T) {

	const password = "securePassword"
	var sampleUser = models.User{
		Password: password,
	}
	sampleUser.HashPassword()

	t.Run("password matches", func(t *testing.T) {

		err := sampleUser.CheckPassword(password)
		assert.Nil(t, err)

	})

	t.Run("password does not match", func(t *testing.T) {

		err := sampleUser.CheckPassword("failPw")
		assert.NotNil(t, err)

	})
}

func TestHashPassword(t *testing.T) {

	const password = "secure123"
	var sampleUser = models.User{
		Password: password,
	}
	sampleUser.HashPassword()

	assert.NotEqual(t, password, sampleUser.Password)
}
