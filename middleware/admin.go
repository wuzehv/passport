package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wuzehv/passport/model"
	"github.com/wuzehv/passport/service/db"
	"github.com/wuzehv/passport/util"
	"github.com/wuzehv/passport/util/config"
	"net/http"
	"strconv"
	"time"
)

// Admin admin页面
func Admin() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(util.CookieFlag)
		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}

		uid, err := strconv.Atoi(token[32:])
		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}

		var u model.User
		db.Db.First(&u, uid)
		if u.Id == 0 || u.Status != model.StatusNormal {
			fmt.Fprintln(c.Writer, util.UserDisabled.Msg(nil))
			return
		}

		// 判断登录是否过期
		if u.Token != token || time.Now().After(u.ExpireTime) {
			// 显式的删除cookie
			c.SetCookie(util.CookieFlag, "false", -1, "/", "", !config.IsDev(), true)

			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}

		c.Set(util.User, &u)
	}
}
