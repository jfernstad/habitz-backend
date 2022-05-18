package auth

import "github.com/golang-jwt/jwt"

type HabitzJWTClaims struct {
	jwt.StandardClaims
	Firstname string `json:"firstname"`
	UserID    string `json:"user_id"`
}
