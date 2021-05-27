package repository

import (
	"gorm.io/gorm"
	"harke.me/showcase-auth/pkg/repository/models"
)

type Storage interface {
	RunMigrations() error
	CreateUser(models.User) error
	FindUser(string) (*models.User, error)
}

func (s *storage) RunMigrations() error {
	return s.db.AutoMigrate(&models.User{})
}

type storage struct {
	db *gorm.DB
}

func NewStorage(db *gorm.DB) Storage {
	return &storage{
		db: db,
	}
}

func (s *storage) CreateUser(user models.User) error {

	result := s.db.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *storage) FindUser(username string) (*models.User, error) {
	var user models.User
	result := s.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
