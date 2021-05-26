package auth

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	jwtWrapper := JwtWrapper{
		SecretKey:       "verysecretkey",
		Issuer:          "AuthService",
		ExpirationHours: 24,
	}

	generatedToken, err := jwtWrapper.GenerateToken("SampleName", "admin")
	assert.NoError(t, err)

	os.Setenv("testToken", generatedToken)
}

func TestValidateToken(t *testing.T) {

	jwtWrapper := JwtWrapper{
		SecretKey:       "verysecretkey",
		Issuer:          "AuthService",
		ExpirationHours: 24,
	}

	encodedToken, _ := jwtWrapper.GenerateToken("SampleName", "admin")

	claims, _ := jwtWrapper.ValidateToken(encodedToken)

	assert.Equal(t, "SampleName", claims.Subject)
	assert.Equal(t, "admin", claims.Role)
	assert.Equal(t, "AuthService", claims.Issuer)
}
