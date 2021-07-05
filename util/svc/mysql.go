package svc

import (
	"errors"
	"github.com/wuzehv/passport/model"
	"github.com/wuzehv/passport/service/db"
	"github.com/wuzehv/passport/util/static"
	"gorm.io/gorm"
	"time"
)

type mysql struct {
	data
}

func (j *mysql) GenToken(userId, clientId uint) (string, error) {
	return model.NewSession(userId, clientId)
}

func (j *mysql) ValidToken(token string, user *model.User) error {
	var s model.Session
	err := s.GetByToken(token)
	if err != nil {
		return err
	}

	if s.Id == 0 {
		return static.SessionNotExists
	}

	if err = db.Db.First(user, s.UserId).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return static.SystemError
	}

	if user.Id == 0 || user.Status != model.StatusNormal {
		return static.UserDisabled
	}

	// 客户端和session不匹配
	//if cl.Id != s.ClientId {
	//	c.AbortWithStatusJSON(http.StatusOK, static.SystemError.Msg("session与客户端不匹配"))
	//}

	// 过期检测
	if time.Now().After(s.ExpireTime) {
		return static.SessionExpired
	}

	return nil
}
