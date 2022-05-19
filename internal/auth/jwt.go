package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	HabitzJWTIssuer   = "https://api.myhabitz.app"
	HabitzJWTAudience = "https://api.myhabitz.app"

	SigningKeyLength = 32
)

type HabitzJWTClaims struct {
	jwt.StandardClaims
	Firstname string `json:"firstname"`
	// UserID    string `json:"user_id"`
}

// Lets start with symmetric encryption key
// We're the only one that care about this JWT token anyway

type JWTServicer interface {
	NewToken(claims *HabitzJWTClaims, expiration *time.Time) (token string, err error)
	VerifyToken(tokenString string) (bool, *HabitzJWTClaims, error)
}

type jwtService struct {
	signingKey []byte
}

func NewJWTService(signingKey []byte) JWTServicer {
	return &jwtService{
		signingKey: signingKey,
	}
}

func (j *jwtService) NewToken(claims *HabitzJWTClaims, expiration *time.Time) (string, error) {

	// Fill in missing data
	claims.ExpiresAt = expiration.Unix()
	claims.Issuer = HabitzJWTIssuer
	claims.Audience = HabitzJWTAudience

	// Make sure we have a good key
	if len(j.signingKey) != SigningKeyLength {
		return "", errors.New(fmt.Sprintf("invalid signing key length %d", len(j.signingKey)))
	}

	// Add Basic information
	claims.Issuer = HabitzJWTIssuer

	// HS256 - symmetric key used for signing and verifying
	habitzToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := habitzToken.SignedString(j.signingKey)
	if err != nil {
		return "", err
	}

	return tokenString, err
}

func (j *jwtService) VerifyToken(tokenString string) (bool, *HabitzJWTClaims, error) {
	// Make sure we have a good key
	if len(j.signingKey) != SigningKeyLength {
		return false, nil, errors.New(fmt.Sprintf("invalid signing key length %d", len(j.signingKey)))
	}

	claims := &HabitzJWTClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return j.signingKey, nil
		},
	)
	if err != nil {
		return false, nil, err
	}

	var ok bool
	claims, ok = token.Claims.(*HabitzJWTClaims)
	if !ok {
		return false, nil, errors.New("invalid Habitz JWT")
	}

	if claims.Issuer != HabitzJWTIssuer {
		return false, nil, errors.New("incorrect issuer")
	}

	if claims.Audience != HabitzJWTAudience {
		return false, nil, errors.New("incorrect audience")
	}

	return true, claims, nil
}
