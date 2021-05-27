package app

import (
	"log"
	"net/http"

	"harke.me/showcase-auth/pkg/api"
	"harke.me/showcase-auth/pkg/helper"
)

type Server struct {
	userService api.UserService
	jwtWrapper  helper.JwtWrapper
}

func NewServer(userService api.UserService, jwtWrapper helper.JwtWrapper) *Server {
	return &Server{
		userService,
		jwtWrapper,
	}
}

func (s *Server) Run() error {

	s.Routes()

	err := http.ListenAndServe(":8000", nil)

	if err != nil {
		log.Printf("Server - there was an error starting the server: %v", err)
		return err
	}

	return nil
}
