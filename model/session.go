package model

import (
	"errors"
	"github.com/wuzehv/passport/service/db"
	"github.com/wuzehv/passport/util/common"
	"gorm.io/gorm"
	"time"
)

type Session struct {
	Model
	Token      string    `gorm:"unique;not null"`
	UserId     uint      `gorm:"index;not null"`
	ClientId   uint      `gorm:"index;not null"`
	ExpireTime time.Time `gorm:"not null;"`
	Status     uint      `gorm:"not null;type:tinyint unsigned" json:"-"`
}

const (
	// 初始化状态
	StatusInit = iota + 1
	// 已登录
	StatusLogin
	// 已退出
	StatusLogout
)

const ExpireTime = 24 * time.Hour

func (s *Session) Base() {}

func NewSession(userId, clientId uint) Session {
	s := Session{
		Token:      common.GenToken(),
		UserId:     userId,
		ClientId:   clientId,
		Status:     StatusInit,
		ExpireTime: time.Now().Add(ExpireTime),
	}
	db.Db.Create(&s)

	return s
}

func (s *Session) GetByToken(t string) error {
	if err := db.Db.Where("token = ?", t).First(&s).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return nil
}

func LogoutAll(userId uint) error {
	return db.Db.Model(Session{}).Where("user_id = ? and status = ?", userId, StatusLogin).Updates(Session{Status: StatusLogout}).Error
}
