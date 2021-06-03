package app

import (
	"encoding/json"
	"net/http"
	"strings"

	"harke.me/showcase-auth/pkg/api"
)

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !s.verifyTokenAndRole(r, "admin") {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var newUserRequest api.NewUserRequest

	err := json.NewDecoder(r.Body).Decode(&newUserRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.userService.New(newUserRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var creds api.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	signedToken, err := s.userService.Login(creds)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(api.TokenResponse{Token: signedToken})

}

func (s *Server) verifyTokenAndRole(r *http.Request, role string) bool {
	const authHeader = "Authorization"
	const tokenPrefix = "Bearer "

	bearerToken := r.Header.Get(authHeader)
	splitArr := strings.Split(bearerToken, tokenPrefix)
	if len(splitArr) != 2 {
		return false
	}

	claims, err := s.jwtWrapper.ValidateToken(splitArr[1])
	if err != nil {
		return false
	}
	return claims.Role == role
}
