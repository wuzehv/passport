package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/wuzehv/passport/model"
	"github.com/wuzehv/passport/util"
	"net/http"
	"strconv"
)

// Sso sso中心页面入口
func Sso() gin.HandlerFunc {
	return func(c *gin.Context) {
		domain := c.Query(util.Domain)
		var cl model.Client
		err := cl.GetByDomain(domain)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, util.SystemError)
			return
		}

		if cl.Id > 0 && cl.Status != model.StatusNormal {
			c.AbortWithError(http.StatusForbidden, util.ClientDisabled)
			return
		}

		c.Set(util.Sso, cl.Id != 0)
		c.Set(util.Client, &cl)
		c.Set(util.Jump, c.Query(util.Jump))
		c.Set(util.Uid, 0)

		// 根据token解析出用户信息
		token, err := c.Cookie(util.CookieFlag)
		if err != nil {
			return
		}

		uid, err := strconv.Atoi(token[32:])
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, util.TokenParseError)
			return
		}

		c.Set(util.Uid, uid)
	}
}
