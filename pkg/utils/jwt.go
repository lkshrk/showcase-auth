package utils

//go:generate mockgen -destination=../mocks/mock_jwtWrapper.go -package=mocks harke.me/showcase-auth/pkg/utils JwtWrapper

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JwtWrapper interface {
	GenerateToken(name string, role string) (string, error)
	ValidateToken(signedToken string) (*JwtClaim, error)
}

type jwtWrapper struct {
	secretKey       string
	issuer          string
	expirationHours int64
}

type NonParseableError struct{}

func (n *NonParseableError) Error() string {
	return "couldn't parse claims"
}

func newNonParseableError() *NonParseableError {
	return &NonParseableError{}
}

type TokenExpiredError struct{}

func (n *TokenExpiredError) Error() string {
	return "couldn't parse claims"
}

func newTokenExpiredError() *TokenExpiredError {
	return &TokenExpiredError{}
}

func NewJwtWrapper(config AuthConfig) JwtWrapper {
	return &jwtWrapper{
		config.Secret,
		config.Issuer,
		config.ExpirationHours,
	}
}

type JwtClaim struct {
	Role string `json:"role"`
	jwt.StandardClaims
}

func (j *jwtWrapper) GenerateToken(name string, role string) (string, error) {
	claims := &JwtClaim{
		Role: role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(j.expirationHours)).Unix(),
			Issuer:    j.issuer,
			Subject:   name,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.secretKey))
}

func (j *jwtWrapper) ValidateToken(signedToken string) (*JwtClaim, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JwtClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.secretKey), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JwtClaim)
	if !ok {
		err = newNonParseableError()
		return nil, err
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = newTokenExpiredError()
		return nil, err
	}

	return claims, nil

}
