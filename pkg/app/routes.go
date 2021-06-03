package app

//go:generate mockgen -destination=../mocks/mock_userRouteHandler.go -package=mocks harke.me/showcase-auth/pkg/app UserRouteHandler

import (
	"net/http"
)

type UserRouteHandler interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}

func SetupHandlers(userHandler UserRouteHandler) *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("/login", userHandler.Login)
	mux.HandleFunc("/register", userHandler.CreateUser)

	return mux
}
