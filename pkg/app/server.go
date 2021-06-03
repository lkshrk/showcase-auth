package app

import (
	"log"
	"net/http"

	"harke.me/showcase-auth/pkg/api"
	"harke.me/showcase-auth/pkg/utils"
)

type Server struct {
	userService      api.UserService
	UserRouteHandler UserRouteHandler
	jwtWrapper       utils.JwtWrapper
}

func NewServer(userService api.UserService, userRouteHandler UserRouteHandler, jwtWrapper utils.JwtWrapper) *Server {
	return &Server{
		userService,
		userRouteHandler,
		jwtWrapper,
	}
}

func (s *Server) Run() error {

	mux := SetupHandlers(s.UserRouteHandler)

	err := http.ListenAndServe(":8000", mux)

	if err != nil {
		log.Printf("Server - there was an error starting the server: %v", err)
		return err
	}

	return nil
}
