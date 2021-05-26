package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"gorm.io/gorm"
	"harke.me/showcase-auth/database"
	"harke.me/showcase-auth/database/auth"
	"harke.me/showcase-auth/models"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result := database.DB.Where("username = ?", creds.Username).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = user.CheckPassword(creds.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	JwtWrapper := auth.JwtWrapper{
		SecretKey:       "verysecretkey",
		Issuer:          "showcase-auth",
		ExpirationHours: 24,
	}

	signedToken, err := JwtWrapper.GenerateToken(user.Username, user.Role)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(TokenResponse{Token: signedToken})

}

func Register(w http.ResponseWriter, r *http.Request) {

	if !verifyTokenAndRole(r, "admin") {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = user.HashPassword(user.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = user.CreateUser()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func verifyTokenAndRole(r *http.Request, role string) bool {
	bearerToken := r.Header.Get("Authorization")
	splitArr := strings.Split(bearerToken, "Bearer ")
	if len(splitArr) != 2 {
		return false
	}

	JwtWrapper := auth.JwtWrapper{
		SecretKey: "verysecretkey",
		Issuer:    "AuthService",
	}

	claims, err := JwtWrapper.ValidateToken(splitArr[1])
	if err != nil {
		return false
	}
	return claims.Role == role
}
