package app_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"harke.me/showcase-auth/pkg/api"
	"harke.me/showcase-auth/pkg/app"
	"harke.me/showcase-auth/pkg/mocks"
)

func TestCreateUserTokenValidation(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserService := mocks.NewMockUserService(mockCtrl)

	cut := app.NewUserRouteHandler(mockUserService)

	t.Run("No POST request", func(t *testing.T) {

		req, err := http.NewRequest("GET", "", nil)
		if err != nil {
			return
		}
		recorder := httptest.NewRecorder()

		cut.CreateUser(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode)
	})

	t.Run("Unauthorized with no token", func(t *testing.T) {

		req, err := http.NewRequest("POST", "", nil)
		if err != nil {
			return
		}
		recorder := httptest.NewRecorder()

		cut.CreateUser(recorder, req)

		assert.Equal(t, http.StatusUnauthorized, recorder.Result().StatusCode)
	})

	t.Run("Unauthorized with malformated token", func(t *testing.T) {

		req, err := http.NewRequest("POST", "", nil)
		if err != nil {
			return
		}
		req.Header.Add("Authorization", "123abc")
		recorder := httptest.NewRecorder()

		cut.CreateUser(recorder, req)

		assert.Equal(t, http.StatusUnauthorized, recorder.Result().StatusCode)
	})

	t.Run("Unauthorized with invalid token", func(t *testing.T) {

		const bearer = "123abc"

		req, err := http.NewRequest("POST", "", nil)
		if err != nil {
			return
		}
		req.Header.Add("Authorization", "Bearer "+bearer)
		recorder := httptest.NewRecorder()

		mockUserService.EXPECT().ValidateTokenAndRole(gomock.Eq(bearer), "admin").Return(false)

		cut.CreateUser(recorder, req)

		assert.Equal(t, http.StatusUnauthorized, recorder.Result().StatusCode)
	})

	t.Run("Authorized with bad request", func(t *testing.T) {

		const bearer = "123abc"

		req, err := http.NewRequest("POST", "", strings.NewReader(""))
		if err != nil {
			return
		}
		req.Header.Add("Authorization", "Bearer "+bearer)
		recorder := httptest.NewRecorder()

		mockUserService.EXPECT().ValidateTokenAndRole(gomock.Eq(bearer), "admin").Return(true)

		cut.CreateUser(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode)
	})

}

func TestCreateNewUser(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserService := mocks.NewMockUserService(mockCtrl)

	cut := app.NewUserRouteHandler(mockUserService)

	const bearer = "123abc"

	userRequest := api.NewUserRequest{
		Username: "user123",
		Password: "securepw",
		Role:     "some",
	}
	requestJson := fmt.Sprintf("{\"username\":\"%s\",\"password\":\"%s\",\"role\":\"%s\"}", userRequest.Username, userRequest.Password, userRequest.Role)

	mockUserService.EXPECT().ValidateTokenAndRole(gomock.Eq(bearer), "admin").AnyTimes().Return(true)

	t.Run("Server error with valid request", func(t *testing.T) {

		req, err := http.NewRequest("POST", "", strings.NewReader(requestJson))
		if err != nil {
			return
		}
		req.Header.Add("Authorization", "Bearer "+bearer)

		recorder := httptest.NewRecorder()

		mockUserService.EXPECT().New(gomock.Eq(userRequest)).Return(errors.New("err"))

		cut.CreateUser(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Result().StatusCode)
	})

	t.Run("Valid request", func(t *testing.T) {

		req, err := http.NewRequest("POST", "", strings.NewReader(requestJson))
		if err != nil {
			return
		}
		req.Header.Add("Authorization", "Bearer "+bearer)

		recorder := httptest.NewRecorder()

		mockUserService.EXPECT().New(gomock.Eq(userRequest)).Return(nil)

		cut.CreateUser(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Result().StatusCode)
	})

}

func TestLogin(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserService := mocks.NewMockUserService(mockCtrl)

	cut := app.NewUserRouteHandler(mockUserService)

	t.Run("No POST request", func(t *testing.T) {

		req, err := http.NewRequest("GET", "", nil)
		if err != nil {
			return
		}
		recorder := httptest.NewRecorder()

		cut.Login(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode)
	})

	t.Run("invalid body", func(t *testing.T) {

		req, err := http.NewRequest("POST", "", strings.NewReader(""))
		if err != nil {
			return
		}
		recorder := httptest.NewRecorder()

		cut.Login(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode)
	})

	t.Run("invalid credentials", func(t *testing.T) {

		loginRequest := api.LoginRequest{
			Username: "user123",
			Password: "securepw",
		}
		requestJson := fmt.Sprintf("{\"username\":\"%s\",\"password\":\"%s\"}", loginRequest.Username, loginRequest.Password)

		req, err := http.NewRequest("POST", "", strings.NewReader(requestJson))
		if err != nil {
			return
		}
		recorder := httptest.NewRecorder()

		mockUserService.EXPECT().Login(gomock.Eq(loginRequest)).Return("", errors.New("err"))

		cut.Login(recorder, req)

		assert.Equal(t, http.StatusUnauthorized, recorder.Result().StatusCode)
	})

	t.Run("login successful", func(t *testing.T) {

		const tokenString = "someToken"

		loginRequest := api.LoginRequest{
			Username: "user123",
			Password: "securepw",
		}
		requestJson := fmt.Sprintf("{\"username\":\"%s\",\"password\":\"%s\"}", loginRequest.Username, loginRequest.Password)

		req, err := http.NewRequest("POST", "", strings.NewReader(requestJson))
		if err != nil {
			return
		}
		recorder := httptest.NewRecorder()

		mockUserService.EXPECT().Login(gomock.Eq(loginRequest)).Return(tokenString, nil)

		cut.Login(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Result().StatusCode)
		assert.Equal(t, fmt.Sprintf("{\"token\":\"%s\"}\n", tokenString), recorder.Body.String())
	})

}
