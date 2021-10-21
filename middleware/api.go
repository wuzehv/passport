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

		var u model.User
		defer func() {
			if u.Id > 0 && u.Status == model.StatusNormal {
				return
			}

			f()
			c.Redirect(http.StatusTemporaryRedirect, "/")
			c.Abort()
		}()

		token, err := c.Cookie(static.CookieFlag)
		if err != nil {
			return
		}

		uid, err := strconv.Atoi(token[32:])
		if err != nil {
			return
		}

		if err = db.Db.First(&u, uid).Error; err != nil {
			return
		}

		// 判断登录是否过期
		if u.Token != token || time.Now().After(u.ExpireTime) {
			return
		}

		c.Set(static.User, &u)
	}
}
