package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/wuzehv/passport/model"
	"github.com/wuzehv/passport/service/db"
	"github.com/wuzehv/passport/util/config"
	"github.com/wuzehv/passport/util/static"
	"net/http"
	"strconv"
	"time"
)

// Api admin接口
func Api() gin.HandlerFunc {
	return func(c *gin.Context) {
		f := func() {
			// 显式的删除cookie
			c.SetCookie(static.CookieFlag, "false", -1, "/", "", !config.IsDev(), true)
		}

		token, err := c.Cookie(static.CookieFlag)
		if err != nil {
			f()
			//c.AbortWithStatusJSON(http.StatusTemporaryRedirect, static.UserNotLogin.Msg(nil))
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}

		uid, err := strconv.Atoi(token[32:])
		if err != nil {
			f()
			c.Redirect(http.StatusTemporaryRedirect, "/")
			//c.AbortWithStatusJSON(http.StatusTemporaryRedirect, static.UserNotLogin.Msg(nil))
			return
		}

		var u model.User
		if err = db.Db.First(&u, uid).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, static.SystemError.Msg(err))
			return
		}

		// 判断登录是否过期
		if u.Token != token || time.Now().After(u.ExpireTime) {
			f()
			c.AbortWithStatusJSON(http.StatusTemporaryRedirect, static.UserNotLogin.Msg(nil))
			return
		}

		c.Set(static.User, &u)
	}
}
