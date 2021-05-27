package helper

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JwtWrapper struct {
	secretKey       string
	issuer          string
	expirationHours int64
}

func NewJwtWrapper(secretKey string, issuer string, expirationHours int64) JwtWrapper {
	return JwtWrapper{
		secretKey,
		issuer,
		expirationHours,
	}
}

type JwtClaim struct {
	Role string `json:"role"`
	jwt.StandardClaims
}

func (j *JwtWrapper) GenerateToken(name string, role string) (string, error) {
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

func (j *JwtWrapper) ValidateToken(signedToken string) (*JwtClaim, error) {
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
		err = errors.New("couldn't parse claims")
		return nil, err
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("JWT is expired")
		return nil, err
	}

	return claims, nil

}
