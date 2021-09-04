package main

import (
	"github.com/wuzehv/passport/router"
	"github.com/wuzehv/passport/util/config"
	_ "github.com/wuzehv/passport/util/journal"
	"log"
)

// @Title 单点登录系统
// @Version 1.0
// @Description 单点登录系统api文档
// @Contact.name wuzehui
// @Contact.email
// @Host sso.com:8099
// @BasePath /
func main() {
	r := router.InitRouter()
	if err := r.Run(config.App.Port); err != nil {
		log.Fatalf("server run error: %v\n", err)
	}
}
