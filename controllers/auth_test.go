package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"harke.me/showcase-auth/database"
	"harke.me/showcase-auth/models"
)

// ToDo: Fix Db tests
// func TestSignUp(t *testing.T) {

// 	user := models.User{
// 		Name:     "Test User",
// 		Role:     "usr",
// 		Password: "password",
// 	}

// 	payload, _ := json.Marshal(&user)

// 	request, _ := http.NewRequest(http.MethodPost, "", bytes.NewBuffer(payload))
// 	w := httptest.NewRecorder()

// 	_ = database.InitDatabase()
// 	database.DB.AutoMigrate(&models.User{})

// 	Register(w, request)

// 	assert.Equal(t, 201, w.Code)

// 	//ToDo: Validate data is written to db
// }

func TestRegisterInvalidJSON(t *testing.T) {
	user := "test"

	payload, _ := json.Marshal(&user)

	request, _ := http.NewRequest(http.MethodPost, "", bytes.NewBuffer(payload))
	w := httptest.NewRecorder()

	Register(w, request)

	assert.Equal(t, 400, w.Code)
}

// ToDo: Fix Db tests
// func TestLogin(t *testing.T) {
// 	user := Credentials{
// 		Username: "jondoe",
// 		Password: "password",
// 	}

// 	payload, _ := json.Marshal(&user)

// 	request, _ := http.NewRequest(http.MethodPost, "", bytes.NewBuffer(payload))

// 	w := httptest.NewRecorder()

// 	_ = database.InitDatabase()

// 	database.DB.AutoMigrate(&models.User{})

// 	Login(w, request)

// 	assert.Equal(t, 200, w.Code)

// }

func TestLoginInvalidRequest(t *testing.T) {
	user := "trash"

	payload, _ := json.Marshal(&user)

	request, _ := http.NewRequest(http.MethodPost, "", bytes.NewBuffer(payload))

	w := httptest.NewRecorder()

	_ = database.InitDatabase()

	database.DB.AutoMigrate(&models.User{})

	Login(w, request)

	assert.Equal(t, 400, w.Code)

}

func TestLoginInvalidCreds(t *testing.T) {
	user := Credentials{
		Username: "jondoe",
		Password: "secret",
	}

	payload, _ := json.Marshal(&user)

	request, _ := http.NewRequest(http.MethodPost, "", bytes.NewBuffer(payload))

	w := httptest.NewRecorder()

	_ = database.InitDatabase()

	database.DB.AutoMigrate(&models.User{})

	Login(w, request)

	assert.Equal(t, 401, w.Code)

}
