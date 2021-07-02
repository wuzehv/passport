package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt"
)

type Claims struct {
	Data interface{}
	jwt.StandardClaims
}

func GenToken(data interface{}, secret []byte) (string, error) {
	v := Claims{
		data,
		jwt.StandardClaims{},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, v)
	return token.SignedString(secret)
}

func ParseToken(tokenString string, secret []byte) (*Claims, error) {
	v := new(Claims)
	token, err := jwt.ParseWithClaims(tokenString, v, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if res, ok := token.Claims.(*Claims); ok {
		return res, nil
	}

	return nil, errors.New("jwt token解析错误")
}
