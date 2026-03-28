package utils

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func ValidateToken(toeknStr string) (*JWTClaims, error) {
	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.ParseWithClaims(toeknStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("Invalid Token")
	}

	return claims, nil
}
