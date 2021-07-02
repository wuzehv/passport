package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/wuzehv/passport/model"
	"github.com/wuzehv/passport/service/db"
	"github.com/wuzehv/passport/service/rdb"
	"github.com/wuzehv/passport/util/common"
	"github.com/wuzehv/passport/util/static"
	"net/http"
	"net/url"
	"time"
)

// Svc svc调用入口，校验token
func Svc() gin.HandlerFunc {
	return func(c *gin.Context) {
		var res static.SvcRequest
		if c.ShouldBind(&res) != nil {
			c.AbortWithStatusJSON(http.StatusOK, static.ParamsError.Msg(nil))
			return
		}

		var u model.User
		if rdb.GetJson(res.Token, &u) {
			c.AbortWithStatusJSON(http.StatusOK, static.Success.Msg(u))
			return
		}

		domain := res.Domain
		domain, _ = url.QueryUnescape(domain)

		var cl model.Client
		err := cl.GetByDomain(domain)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, static.SystemError.Msg(nil))
			return
		}

		if cl.Id == 0 || cl.Status != model.StatusNormal {
			c.AbortWithStatusJSON(http.StatusOK, static.ClientDisabled.Msg(nil))
			return
		}

		m := make(map[string]string)
		m[static.Token] = res.Token
		m[static.Timestamp] = res.Timestamp
		m[static.Domain] = res.Domain
		if common.GenSign(m, cl.Secret) != res.Sign {
			c.AbortWithStatusJSON(http.StatusOK, static.SignatureError.Msg(nil))
			return
		}

		t := res.Token

		var s model.Session
		err = s.GetByToken(t)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, static.SystemError.Msg(nil))
			return
		}

		if s.Id == 0 {
			c.AbortWithStatusJSON(http.StatusOK, static.TokenNotExists.Msg(nil))
			return
		}

		db.Db.First(&u, s.UserId)

		if u.Id == 0 || u.Status != model.StatusNormal {
			c.AbortWithStatusJSON(http.StatusOK, static.UserDisabled.Msg(nil))
			return
		}

		// 客户端和session不匹配
		if cl.Id != s.ClientId {
			c.AbortWithStatusJSON(http.StatusOK, static.SystemError.Msg(nil))
			return
		}

		// 过期检测
		if time.Now().After(s.ExpireTime) {
			c.AbortWithStatusJSON(http.StatusOK, static.SessionExpired.Msg(nil))
			return
		}

		c.Set(static.Session, &s)
		c.Set(static.User, &u)
	}
}
