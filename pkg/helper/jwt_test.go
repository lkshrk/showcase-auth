package helper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"harke.me/showcase-auth/pkg/helper"
)

func TestValidateToken(t *testing.T) {

	cut := helper.NewJwtWrapper("secretKey123", "issuer", 24)

	t.Run("non parseable token", func(t *testing.T) {

		_, err := cut.ValidateToken("abc")
		assert.Error(t, err)
	})

	t.Run("expired token", func(t *testing.T) {

		const EXPIRED_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYWRtaW4iLCJleHAiOjEzMjIzMDIwMjAsImlzcyI6Imlzc3VlciIsInN1YiI6InVzciJ9.DLO6cd2bBSCfA8kLOUtaPi6hsXc1_6ysEPpisNwoEF8"

		_, err := cut.ValidateToken(EXPIRED_TOKEN)
		assert.Error(t, err)
	})

	t.Run("valid token", func(t *testing.T) {

		const VALID_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYWRtaW4iLCJleHAiOjI2MjIzMDIwMjAsImlzcyI6Imlzc3VlciIsInN1YiI6InVzciJ9.Qgjo95kPB3iLaZMGhnrBr89HgQeUNNbdmW1ghkPy98c"
		_, err := cut.ValidateToken(VALID_TOKEN)
		assert.Nil(t, err)
	})

}

func TestJwtWrapper(t *testing.T) {

	const issuer = "example.com"
	const name = "peter"
	const role = "user"

	cut := helper.NewJwtWrapper("secretKey123", issuer, 24)

	t.Run("create and validate token", func(t *testing.T) {
		token, err := cut.GenerateToken(name, role)
		if err != nil {
			return
		}

		claim, err := cut.ValidateToken(token)
		if err != nil {
			return
		}

		assert.Equal(t, name, claim.Subject)
		assert.Equal(t, role, claim.Role)
		assert.Equal(t, issuer, claim.Issuer)
	})

}
