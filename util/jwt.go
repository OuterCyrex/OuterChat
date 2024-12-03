package util

import (
	"OuterChat/config"
	"github.com/golang-jwt/jwt"
	"time"
)

var jwtKey = []byte(config.JwtKey)

type Claims struct {
	UID uint `json:"uid"`
	jwt.StandardClaims
}

func CreateToken(uid uint) (string, error) {
	expireTime := time.Now().Add(time.Hour * 24 * 30).Unix()
	claims := Claims{
		UID: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime,
			IssuedAt:  time.Now().Unix(),
			Issuer:    "Whatever Issuer",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims.(*Claims), err
}
