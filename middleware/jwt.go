package middleware

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

var jwtSecret = []byte("xiaobaicai")

type Claims struct {
	UserId int64
	jwt.RegisteredClaims
}

func GetToken(id int64) (string, error) {
	claims := Claims{UserId: id, RegisteredClaims: jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(3 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
	}}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("获取token失败")
	}
	return tokenString, nil
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, _ := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	return token, claims, nil
}
