package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wuzehv/passport/app/api/v1/action"
	"github.com/wuzehv/passport/app/api/v1/client"
	"github.com/wuzehv/passport/app/api/v1/index"
	"github.com/wuzehv/passport/app/sso"
	"github.com/wuzehv/passport/app/svc"
	"github.com/wuzehv/passport/middleware"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/wuzehv/passport/doc"
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
	r := router.Group("/api/v1")
	r.Use(middleware.Api())
	{
		r.GET("/index", index.Index)

		r.GET("/actions", action.Index)
		r.POST("/actions", action.Add)
		r.PUT("/actions/:id", action.Update)

		r.GET("/clients", client.Index)
		r.POST("/clients", client.Add)
		r.PUT("/clients/:id", client.Update)
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
