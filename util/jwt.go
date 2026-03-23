package util

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var JwtSecret = []byte("myshop-secret-key")

type MyClaims struct { //
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, username string) (string, error) {
	claims := MyClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "myshop",
			Subject:   username,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //HS256对称签名
	return token.SignedString(JwtSecret)
}
func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return JwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*MyClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
