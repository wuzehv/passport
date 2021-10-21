package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wuzehv/passport/app/api/v1/client"
	"github.com/wuzehv/passport/app/api/v1/index"
	"github.com/wuzehv/passport/app/api/v1/record"
	"github.com/wuzehv/passport/app/api/v1/user"
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
	constructNoLogin(router)
}

// sso主页
func constructSso(router *gin.Engine) {
	r := router.Group("/")
	r.Use(middleware.Base(), middleware.Cors(), middleware.Sso())
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
	r.Use(middleware.Base(), middleware.Rate(), middleware.Cors(), middleware.Api())
	{
		r.GET("/index", index.Index)

		r.GET("/users", user.Index)
		r.POST("/users", user.Add)
		r.PUT("/users/:id", user.Update)
		r.POST("/users/:id/toggle-status", user.ToggleStatus)
		r.POST("/users/reset-password", user.ResetPassword)

		r.GET("/clients", client.Index)
		r.POST("/clients", client.Add)
		r.PUT("/clients/:id", client.Update)
		r.POST("/clients/:id/toggle-status", client.ToggleStatus)
		r.GET("/clients/check-callback", client.CheckCallback)

		r.GET("/records", record.Index)

		r.GET("/records", client.Index)
	}
}

// 对外接口
func constructSvc(router *gin.Engine) {
	r := router.Group("/svc")
	r.Use(middleware.Base(), middleware.Svc())
	{
		r.POST("/userinfo", svc.Userinfo)
		r.POST("/session", svc.Session)
	}
}

// 不需要登录的接口
func constructNoLogin(router *gin.Engine) {
	r := router.Group("/common")
	r.Use(middleware.Base(), middleware.Rate(), middleware.Cors())
	{
		r.GET("/reset-password", user.ResetPasswordPage)
		r.POST("/reset-password", user.DoResetPassword)
	}
}
