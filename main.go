package main

import (
	"github.com/wuzehv/passport/router"
	"github.com/wuzehv/passport/util/config"
	"log"
)

// @title 单点登录系统
// @version 1.0
// @description 单点登录系统api文档
// @contact.name wuzehui
// @contact.email
// @host http://sso.com:8099
// @BasePath /
func main() {
	r := router.InitRouter()
	if err := r.Run(config.App.Port); err != nil {
		log.Fatalf("server run error: %v\n", err)
	}
}
