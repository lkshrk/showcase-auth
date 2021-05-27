package app

import (
	"net/http"
)

func (s *Server) Routes() {
	http.HandleFunc("/login", s.Login)
	// http.HandleFunc("/refresh", Refresh)
	http.HandleFunc("/register", s.CreateUser)
}
