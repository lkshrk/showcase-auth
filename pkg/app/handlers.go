package app

import (
	"encoding/json"
	"net/http"
	"strings"

	"harke.me/showcase-auth/pkg/api"
)

type userHandler struct {
	userService api.UserService
}

func NewUserRouteHandler(userService api.UserService) UserRouteHandler {
	return &userHandler{
		userService: userService,
	}
}

func (u *userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !u.extractAndValidateToken(r, "admin") {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var newUserRequest api.NewUserRequest

	err := json.NewDecoder(r.Body).Decode(&newUserRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = u.userService.New(newUserRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func (u *userHandler) Login(w http.ResponseWriter, r *http.Request) {

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

	signedToken, err := u.userService.Login(creds)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(api.TokenResponse{Token: signedToken})

}

func (u *userHandler) extractAndValidateToken(r *http.Request, role string) bool {
	const authHeader = "Authorization"
	const tokenPrefix = "Bearer "

	bearerToken := r.Header.Get(authHeader)
	splitArr := strings.Split(bearerToken, tokenPrefix)
	if len(splitArr) != 2 {
		return false
	}

	return u.userService.ValidateTokenAndRole(splitArr[1], role)

}
