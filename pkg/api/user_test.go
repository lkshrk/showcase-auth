package api_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"harke.me/showcase-auth/pkg/api"
	"harke.me/showcase-auth/pkg/mocks"
	"harke.me/showcase-auth/pkg/repository/models"
)

type eqUserRequestUserModelMatcher struct {
	user models.User
}

func (e eqUserRequestUserModelMatcher) Matches(x interface{}) bool {
	user, ok := x.(models.User)
	if !ok {
		return false
	}

	err := user.CheckPassword(e.user.Password)
	if err != nil {
		return false
	}

	return e.user.Username == user.Username && e.user.Role == user.Role
}

func (e eqUserRequestUserModelMatcher) String() string {
	return fmt.Sprintf("matches user %#v", e.user)
}

func EqCreateUserParams(user models.User) gomock.Matcher {
	return eqUserRequestUserModelMatcher{user: user}
}

func TestNewUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepo := mocks.NewMockUserRepository(mockCtrl)
	jwtWrapper := mocks.NewMockJwtWrapper(mockCtrl)

	cut := api.NewUserService(mockRepo, jwtWrapper)

	t.Run("Create fails", func(t *testing.T) {

		userRequest := api.NewUserRequest{
			Username: "user",
			Password: "password",
			Role:     "role",
		}

		expectedErr := errors.New("an err")

		mockRepo.EXPECT().CreateUser(gomock.Any()).Return(expectedErr)

		err := cut.New(userRequest)
		assert.Equal(t, expectedErr, err)

	})

	t.Run("Create success", func(t *testing.T) {
		userRequest := api.NewUserRequest{
			Username: "user",
			Password: "password",
			Role:     "role",
		}

		expectedUser := models.User{
			Username: userRequest.Username,
			Password: userRequest.Password,
			Role:     userRequest.Role,
		}

		mockRepo.EXPECT().CreateUser(EqCreateUserParams(expectedUser)).Return(nil)

		err := cut.New(userRequest)
		assert.Nil(t, err)
	})

}

func TestLogin(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepo := mocks.NewMockUserRepository(mockCtrl)
	mockJwtWrapper := mocks.NewMockJwtWrapper(mockCtrl)

	cut := api.NewUserService(mockRepo, mockJwtWrapper)

	credentials := api.LoginRequest{
		Username: "username",
		Password: "securepassword",
	}

	t.Run("user not found", func(t *testing.T) {

		expected := errors.New("some err")

		mockRepo.EXPECT().FindUser(gomock.Eq(credentials.Username)).Return(nil, expected)

		_, err := cut.Login(credentials)

		assert.Equal(t, expected, err)
	})

	t.Run("password invalid", func(t *testing.T) {

		user := models.User{
			Username: credentials.Username,
			Password: "anotherPassword",
		}
		user.HashPassword()

		mockRepo.EXPECT().FindUser(gomock.Eq(credentials.Username)).Return(&user, nil)

		_, err := cut.Login(credentials)
		assert.NotNil(t, err)

	})

	t.Run("succesful login", func(t *testing.T) {

		const validToken = "I_AM_A_VALID_TOKEN"
		user := models.User{
			Username: credentials.Username,
			Password: credentials.Password,
			Role:     "someRole",
		}
		user.HashPassword()

		mockRepo.EXPECT().FindUser(gomock.Eq(credentials.Username)).Return(&user, nil)
		mockJwtWrapper.EXPECT().GenerateToken(gomock.Eq(user.Username), gomock.Eq(user.Role)).Return(validToken, nil)

		token, err := cut.Login(credentials)

		assert.Nil(t, err)
		assert.Equal(t, validToken, token)

	})
}
