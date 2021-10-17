package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/wuzehv/passport/model"
	"github.com/wuzehv/passport/util/static"
	"strconv"
)

// Sso sso中心页面入口
func Sso() gin.HandlerFunc {
	return func(c *gin.Context) {
		domain := c.Query(static.Domain)
		var cl model.Client
		err := cl.GetByDomain(domain)
		if err != nil {
			c.AbortWithStatusJSON(static.SystemError.Msg(err))
			return
		}

		if cl.Id > 0 && cl.Status != model.StatusNormal {
			c.AbortWithStatusJSON(static.ClientDisabled.Msg(nil))
			return
		}

		c.Set(static.Sso, cl.Id != 0)
		c.Set(static.Client, &cl)
		c.Set(static.Jump, c.Query(static.Jump))
		c.Set(static.Uid, 0)

		// 根据token解析出用户信息
		token, err := c.Cookie(static.CookieFlag)
		if err != nil {
			return
		}

		uid, err := strconv.Atoi(token[32:])
		if err != nil {
			c.AbortWithStatusJSON(static.TokenParseError.Msg(nil))
			return
		}

		c.Set(static.Uid, uid)
	}
}
