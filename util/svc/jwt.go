package svc

import (
	"encoding/json"
	jwtBase "github.com/golang-jwt/jwt"
	"github.com/wuzehv/passport/model"
	"github.com/wuzehv/passport/service/db"
	"github.com/wuzehv/passport/util/journal"
	"github.com/wuzehv/passport/util/jwt"
	"github.com/wuzehv/passport/util/static"
)

type j struct {
	*model.User
}

func (d *j) GenToken(userId, clientId uint) (string, error) {
	var u model.User
	if err := db.Db.First(&u).Error; err != nil {
		return "", err
	}

	return jwt.GenToken(j{&u}, "")
}

func (d *j) ValidToken(token string, user *model.User) error {
	x, err := jwt.ValidToken(token, "")

	if err != nil {
		switch err.(*jwtBase.ValidationError).Errors {
		case jwtBase.ValidationErrorExpired:
			return static.SessionExpired
		default:
			journal.Error("jwt_svc", err.Error())
			return static.SystemError
		}
	}

	// 解析出用户信息
	c, err := json.Marshal(x.Data)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(c, user); err != nil {
		return err
	}

	return nil
}
