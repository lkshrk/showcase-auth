package api

import (
	"harke.me/showcase-auth/pkg/helper"
	"harke.me/showcase-auth/pkg/repository/models"
)

type UserService interface {
	New(user NewUserRequest) error
	Login(credentials LoginRequest) (string, error)
}

type UserRepository interface {
	CreateUser(models.User) error
	FindUser(string) (*models.User, error)
}

type userService struct {
	storage    UserRepository
	jwtWrapper helper.JwtWrapper
}

func NewUserService(userRepo UserRepository, jwtWrapper helper.JwtWrapper) UserService {
	return &userService{
		storage:    userRepo,
		jwtWrapper: jwtWrapper,
	}
}

func (u *userService) New(userRequest NewUserRequest) error {

	user := models.User{
		Username: userRequest.Username,
		Role:     userRequest.Role,
		Password: userRequest.Password,
	}
	user.HashPassword()
	err := u.storage.CreateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (u *userService) Login(credentials LoginRequest) (string, error) {

	user, err := u.storage.FindUser(credentials.Username)
	if err != nil {
		return "", err
	}

	err = user.CheckPassword(credentials.Password)
	if err != nil {
		return "", err
	}

	return u.jwtWrapper.GenerateToken(user.Username, user.Role)

}
