package main

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"harke.me/showcase-auth/pkg/api"
	"harke.me/showcase-auth/pkg/app"
	"harke.me/showcase-auth/pkg/repository"
	"harke.me/showcase-auth/pkg/repository/models"
	"harke.me/showcase-auth/pkg/utils"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "this is the startup error: %s\\n", err)
		os.Exit(1)
	}
}

func run() error {

	config, err := utils.LoadConfig()
	if err != nil {
		return err
	}

	storage, err := setupStorage(config.Database)
	if err != nil {
		return err
	}

	jwtWrapper := utils.NewJwtWrapper(config.Auth)

	userService := api.NewUserService(storage, jwtWrapper)

	userRouteHandler := app.NewUserRouteHandler(userService)

	server := app.NewServer(userService, userRouteHandler, jwtWrapper)

	err = server.Run()

	if err != nil {
		return err
	}

	return nil

}

func setupStorage(config utils.DatabaseConfig) (repository.Storage, error) {
	var db *gorm.DB
	var err error

	if config.InMemory {
		db, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	} else {
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.Hostname, config.Port, config.User, config.Password, config.DB)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}
	if err != nil {
		return nil, err
	}

	storage := repository.NewStorage(db)

	err = storage.RunMigrations()
	if err != nil {
		return nil, err
	}

	var users []models.User
	db.Find(&users)
	if len(users) == 0 {
		defaultUser := models.User{
			Username: config.DefaultUser,
			Password: config.DefaultPassword,
			Role:     "admin",
		}
		defaultUser.HashPassword()
		storage.CreateUser(defaultUser)
	}

	return storage, err
}
