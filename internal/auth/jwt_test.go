package auth_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jfernstad/habitz/web/internal/auth"
	"github.com/stretchr/testify/assert"
)

const (
	testSecret = "00000000000000000000000000000000"
)

var (
	testClaims = &auth.HabitzJWTClaims{
		Firstname: "Tester McTestFace",
		StandardClaims: jwt.StandardClaims{
			Subject: "0123456789",
		},
	}
)

func TestTokenSigning(t *testing.T) {

	expiration := time.Now().Add(1 * time.Second)
	signer := auth.NewJWTService([]byte(testSecret))
	token, err := signer.NewToken(testClaims, &expiration)
	assert.Nil(t, err)
	assert.NotEmpty(t, token)
}

func TestInvalidSigningKey(t *testing.T) {
	expiration := time.Now().Add(1 * time.Second)
	signer := auth.NewJWTService([]byte("")) // << Not correct length
	token, err := signer.NewToken(testClaims, &expiration)
	assert.NotNil(t, err)
	assert.Empty(t, token)
}

func TestVerifyToken(t *testing.T) {
	// Create & Sign
	expiration := time.Now().Add(1 * time.Second).UTC()
	signer := auth.NewJWTService([]byte(testSecret))
	token, err := signer.NewToken(testClaims, &expiration)
	assert.Nil(t, err)
	assert.NotEmpty(t, token)

	// Verify
	ok, claims, err := signer.VerifyToken(token)
	assert.True(t, ok)
	assert.NotNil(t, claims)
	assert.Nil(t, err)
}

func TestEmptyValidationKey(t *testing.T) {
	// Create & Sign
	expiration := time.Now().Add(1 * time.Second).UTC()
	signer := auth.NewJWTService([]byte(testSecret))
	token, err := signer.NewToken(testClaims, &expiration)

	// Verify
	verifyer := auth.NewJWTService([]byte(""))
	ok, claims, err := verifyer.VerifyToken(token)
	assert.False(t, ok)
	assert.Nil(t, claims)
	assert.NotNil(t, err)
}

func TestWrongValidationKey(t *testing.T) {
	// Create & Sign
	expiration := time.Now().Add(1 * time.Second).UTC()
	signer := auth.NewJWTService([]byte(testSecret))
	token, err := signer.NewToken(testClaims, &expiration)

	// Verify with wrong key
	verifyer := auth.NewJWTService([]byte("11111111111111111111111111111111"))
	ok, claims, err := verifyer.VerifyToken(token)
	assert.False(t, ok)
	assert.Nil(t, claims)
	assert.NotNil(t, err)
}

func TestExpiredToken(t *testing.T) {
	// Create & Sign
	expiration := time.Now().Add(time.Duration(-10) * time.Second).UTC()
	signer := auth.NewJWTService([]byte(testSecret))
	token, err := signer.NewToken(testClaims, &expiration)
	assert.Nil(t, err)
	assert.NotEmpty(t, token)

	// Verify
	ok, claims, err := signer.VerifyToken(token)
	assert.False(t, ok)
	assert.Nil(t, claims)
	assert.Equal(t, err.(*jwt.ValidationError).Errors, jwt.ValidationErrorExpired)
}
