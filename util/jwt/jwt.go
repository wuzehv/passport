package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"time"
)

type Claims struct {
	Data interface{}
	jwt.StandardClaims
}

const Issuer = "passport"

func GenToken(data interface{}, secret string, exp time.Duration) (string, error) {
	v := Claims{
		data,
		jwt.StandardClaims{
			Issuer:    Issuer,
			ExpiresAt: time.Now().Unix() + int64(exp),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, v)
	return token.SignedString([]byte(secret))
}

func ValidToken(tokenString string, secret string) (*Claims, error) {
	v := new(Claims)
	token, err := jwt.ParseWithClaims(tokenString, v, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if res, ok := token.Claims.(*Claims); ok {
		if err = res.Valid(); err != nil {
			return nil, err
		}

		return res, nil
	}

	return nil, errors.New("jwt token valid error")
}
