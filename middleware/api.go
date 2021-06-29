package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/wuzehv/passport/model"
	"github.com/wuzehv/passport/service/db"
	"github.com/wuzehv/passport/util"
	"github.com/wuzehv/passport/util/config"
	"net/http"
	"strconv"
	"time"
)

// Api admin接口
func Api() gin.HandlerFunc {
	return func(c *gin.Context) {
		f := func() {
			// 显式的删除cookie
			c.SetCookie(util.CookieFlag, "false", -1, "/", "", !config.IsDev(), true)
		}

		token, err := c.Cookie(util.CookieFlag)
		if err != nil {
			f()
			c.AbortWithStatusJSON(http.StatusTemporaryRedirect, util.UserNotLogin.Msg(nil))
			return
		}

		uid, err := strconv.Atoi(token[32:])
		if err != nil {
			f()
			c.AbortWithStatusJSON(http.StatusTemporaryRedirect, util.UserNotLogin.Msg(nil))
			return
		}

		var u model.User
		db.Db.First(&u, uid)
		if u.Id == 0 || u.Status != model.StatusNormal {
			c.AbortWithStatusJSON(http.StatusForbidden, util.UserDisabled.Msg(nil))
			return
		}

		// 判断登录是否过期
		if u.Token != token || time.Now().After(u.ExpireTime) {
			f()
			c.AbortWithStatusJSON(http.StatusTemporaryRedirect, util.UserNotLogin.Msg(nil))
			return
		}

		c.Set(util.User, &u)
	}
}
