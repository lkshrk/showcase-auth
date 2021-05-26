package main

import (
	"log"
	"net/http"

	"harke.me/showcase-auth/controllers"
	"harke.me/showcase-auth/database"
	"harke.me/showcase-auth/models"
)

func main() {

	err := database.InitDatabase()
	if err != nil {
		log.Fatalln("could not create database", err)
	}

	database.DB.AutoMigrate(&models.User{})

	http.HandleFunc("/login", controllers.Login)
	// http.HandleFunc("/refresh", Refresh)
	http.HandleFunc("/register", controllers.Register)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
