package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID       uint `json:"uid"`
	TokenVersion int  `json:"tv"`
	jwt.RegisteredClaims
}

func SignJWT(secret string, uid uint, tokenVersion int, expiresInSec int64) (string, time.Time, error) {
	exp := time.Now().Add(time.Duration(expiresInSec) * time.Second)
	claims := &Claims{
		UserID:       uid,
		TokenVersion: tokenVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	return signed, exp, err
}
