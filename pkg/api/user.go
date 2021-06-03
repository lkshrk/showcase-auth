package api

import (
	"harke.me/showcase-auth/pkg/repository/models"
	"harke.me/showcase-auth/pkg/utils"
)

//go:generate mockgen -destination=../mocks/mock_userRepository.go -package=mocks harke.me/showcase-auth/pkg/api UserRepository
//go:generate mockgen -destination=../mocks/mock_userService.go -package=mocks harke.me/showcase-auth/pkg/api UserService

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
	jwtWrapper utils.JwtWrapper
}

func NewUserService(userRepo UserRepository, jwtWrapper utils.JwtWrapper) UserService {
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
