// 初始化数据库和表结构
package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/wuzehv/passport/model"
	"github.com/wuzehv/passport/util"
	"github.com/wuzehv/passport/util/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var dropDb = `
DROP DATABASE IF EXISTS passport
`

var passportDb = `
CREATE DATABASE IF NOT EXISTS passport DEFAULT CHARACTER SET utf8
`

func main() {
	u := config.Db.User
	passwd := config.Db.Passwd
	host := config.Db.Host
	dbName := config.Db.DbName

	dsn := u + ":" + passwd + "@tcp(" + host + ")/?parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err.Error())
	}

	db.Exec(dropDb)
	log.Println("drop database done")

	db.Exec(passportDb)
	log.Println("create database done")

	db.Exec(passportDb)
	log.Println("create database done")

	dsn = u + ":" + passwd + "@tcp(" + host + ")/" + dbName + "?parseTime=true"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = db.AutoMigrate(model.User{})
	if err != nil {
		panic(err)
	}
	log.Println("create users table done")

	err = db.AutoMigrate(model.Session{})
	if err != nil {
		panic(err)
	}
	log.Println("create tokens table done")

	err = db.AutoMigrate(model.Client{})
	if err != nil {
		panic(err)
	}
	log.Println("create clients table done")

	err = db.AutoMigrate(model.LoginRecord{})
	if err != nil {
		panic(err)
	}
	log.Println("create records table done")

	err = db.AutoMigrate(model.Action{})
	if err != nil {
		panic(err)
	}
	log.Println("create actions table done")

	db.Create(&model.Client{Domain: "client.one.com:8081", Callback: "http://client.one.com:8081/callback", Secret: "123456", Status: model.StatusNormal})
	db.Create(&model.Client{Domain: "client.two.com:8082", Callback: "http://client.two.com:8082/callback", Secret: "123456", Status: model.StatusNormal})

	log.Println("initialize client done")

	u = "api@gmail.com"
	up := "api"
	p := util.GenPassword(up)
	db.Create(&model.User{Email: u, Password: p, Status: model.StatusNormal})
	log.Printf("initialize user: %s, password: %s\n", u, up)
}
