package jwt

import (
	"time"

	"github.com/golang-jwt/jwt"
)

const secretKey = "secret"

func GenerateJWT(userID string, ttl time.Duration) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(ttl).Unix(),
		Issuer:    userID,
	})

	return claims.SignedString([]byte(secretKey))
}

func ParseJWT(cookie string) (userID string, err error) {
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		return "", err
	}

	claims := token.Claims.(*jwt.StandardClaims)

	return claims.Issuer, nil
}
