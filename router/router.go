package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wuzehv/passport/app/admin/action"
	"github.com/wuzehv/passport/app/admin/index"
	"github.com/wuzehv/passport/app/sso"
	"github.com/wuzehv/passport/app/svc"
	"github.com/wuzehv/passport/middleware"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/wuzehv/passport/docs"
)

func construct(router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	constructSso(router)
	constructAdmin(router)
	constructSvc(router)
}

// sso主页
func constructSso(router *gin.Engine) {
	r := router.Group("/")
	r.Use(middleware.Sso())
	{
		r.GET("/", sso.Index)
		r.GET("/sso/index", sso.Index)

		r.POST("/sso/login", sso.Login)
		r.GET("/sso/logout", sso.Logout)
	}
}

// admin内部
func constructAdmin(router *gin.Engine) {
	r := router.Group("/admin")
	r.Use(middleware.Admin())
	{
		r.GET("/index/index", index.Index)
		r.GET("/index/test", index.Test)

		r.POST("/action/index", action.Index)
	}
}

// 对外接口
func constructSvc(router *gin.Engine) {
	r := router.Group("/svc")
	r.Use(middleware.Svc())
	{
		r.POST("/userinfo", svc.Userinfo)
		r.POST("/session", svc.Session)
	}
}
