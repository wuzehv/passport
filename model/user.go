package model

import (
	"errors"
	"github.com/wuzehv/passport/service/db"
	"gorm.io/gorm"
	"time"
)

type User struct {
	Model
	Email      string    `gorm:"unique" json:"email"`
	Realname   string    `gorm:"not null;type:varchar(255)" json:"realname"`
	Gender     int       `gorm:"not null;type:tinyint" json:"gender"`
	Mobile     string    `gorm:"unique;not null;type:varchar(255)" json:"mobile"`
	Password   string    `gorm:"not null;type:varchar(255)" json:"-"`
	Token      string    `gorm:"index;not null;default:''" json:"-"`
	ExpireTime time.Time `gorm:"default:null" json:"-"`
	Status     uint      `gorm:"not null;type:tinyint unsigned" json:"-"`
}

func (u *User) Base() {}

func (u *User) GetByEmail(email string) error {
	if err := db.Db.Where("email = ?", email).First(u).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// sql错误
		return err
	}

	return nil
}
