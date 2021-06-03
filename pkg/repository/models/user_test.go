package models_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"harke.me/showcase-auth/pkg/repository/models"
)

func TestCheckPassword(t *testing.T) {

	const password = "securePassword"
	var cut = models.User{
		Password: password,
	}
	cut.HashPassword()

	t.Run("password matches", func(t *testing.T) {

		err := cut.CheckPassword(password)
		assert.Nil(t, err)

	})

	t.Run("password does not match", func(t *testing.T) {

		err := cut.CheckPassword("failPw")
		assert.NotNil(t, err)

	})
}

func TestHashPassword(t *testing.T) {

	const password = "secure123"
	var cut = models.User{
		Password: password,
	}
	cut.HashPassword()

	assert.NotEqual(t, password, cut.Password)
}
