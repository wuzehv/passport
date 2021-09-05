package db

import (
	"fmt"
	"github.com/wuzehv/passport/util/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var Db *gorm.DB

func init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=true", config.Db.User, config.Db.Password, config.Db.Host, config.Db.DbName, config.Db.Charset)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Warn,
			Colorful:      false,
		},
	)

	var err error
	Db, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                      dsn,
		DisableDatetimePrecision: true,
	}), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Fatalf("mysql init error: %v\n", err)
	}

	sqlDB, err := Db.DB()
	if err != nil {
		log.Fatalf("get sql db error: %v\n", err)
	}

	// 初始化连接池
	sqlDB.SetMaxIdleConns(config.Db.MaxIdleConn)
	sqlDB.SetMaxOpenConns(config.Db.MaxActiveConn)
	sqlDB.SetConnMaxLifetime(config.Db.MaxConnIdleTimeout)
}
