package svc

import (
	"errors"
	"github.com/wuzehv/passport/model"
	"github.com/wuzehv/passport/service/db"
	"github.com/wuzehv/passport/util/journal"
	"github.com/wuzehv/passport/util/static"
	"gorm.io/gorm"
	"time"
)

type mysql struct {
	data
}

func (_ *mysql) Generate(userId, clientId uint) (string, error) {
	return model.NewSession(userId, clientId)
}

func (_ *mysql) Confirm(token string) error {
	var s model.Session

	if err := s.GetByToken(token); err != nil {
		return err
	}

	if s.Status != model.StatusInit {
		journal.Error("confirm_token", "session状态不合法")
		return static.SystemError
	}

	// 更新session状态
	s.Status = model.StatusLogin
	if err := db.Db.Save(&s).Error; err != nil {
		return err
	}

	return nil
}

func (_ *mysql) Valid(token string, user *model.User) error {
	var s model.Session
	err := s.GetByToken(token)
	if err != nil {
		return err
	}

	if s.Id == 0 {
		return static.SessionNotExists
	}

	if s.Status != model.StatusLogin {
		return static.SessionExpired
	}

	// 过期检测
	if time.Now().After(s.ExpireTime) {
		return static.SessionExpired
	}

	if err = db.Db.First(user, s.UserId).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return static.SystemError
	}

	if user.Id == 0 || user.Status != model.StatusNormal {
		return static.UserDisabled
	}

	// 客户端和session不匹配
	//if cl.Id != s.ClientId {
	//	c.AbortWithStatusJSON(static.SystemError.Msg("session与客户端不匹配"))
	//}

	return nil
}

func (_ *mysql) Destroy(userId uint) error {
	return model.LogoutAll(userId)
}
