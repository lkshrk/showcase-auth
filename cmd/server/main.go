package main

import (
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"harke.me/showcase-auth/pkg/api"
	"harke.me/showcase-auth/pkg/app"
	"harke.me/showcase-auth/pkg/helper"
	"harke.me/showcase-auth/pkg/repository"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "this is the startup error: %s\\n", err)
		os.Exit(1)
	}
}

func run() error {

	db, err := setupDatabase()
	if err != nil {
		return err
	}

	storage := repository.NewStorage(db)

	err = storage.RunMigrations()

	if err != nil {
		return err
	}

	jwtWrapper := helper.NewJwtWrapper("verysecretkey", "showcase-auth", 24)

	userService := api.NewUserService(storage, jwtWrapper)

	server := app.NewServer(userService, jwtWrapper)

	err = server.Run()

	if err != nil {
		return err
	}

	return nil

}

func setupDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, err
}
