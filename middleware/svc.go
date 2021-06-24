package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/wuzehv/passport/model"
	"github.com/wuzehv/passport/service/db"
	"github.com/wuzehv/passport/service/rdb"
	"github.com/wuzehv/passport/util"
	"net/http"
	"net/url"
	"time"
)

// Svc svc调用入口，校验token
func Svc() gin.HandlerFunc {
	return func(c *gin.Context) {
		var res util.SvcRequest
		if c.ShouldBind(&res) != nil {
			c.AbortWithStatusJSON(http.StatusOK, util.ParamsError.Msg(nil))
			return
		}

		var u model.User
		if rdb.GetJson(res.Token, &u) {
			c.AbortWithStatusJSON(http.StatusOK, util.Success.Msg(u))
			return
		}

		domain := res.Domain
		domain, _ = url.QueryUnescape(domain)

		var cl model.Client
		err := cl.GetByDomain(domain)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, util.SystemError.Msg(nil))
			return
		}

		if cl.Id == 0 || cl.Status != model.StatusNormal {
			c.AbortWithStatusJSON(http.StatusOK, util.ClientDisabled.Msg(nil))
			return
		}

		m := make(map[string]string)
		m[util.Token] = res.Token
		m[util.Timestamp] = res.Timestamp
		m[util.Domain] = res.Domain
		if util.GenSign(m, cl.Secret) != res.Sign {
			c.AbortWithStatusJSON(http.StatusOK, util.SignatureError.Msg(nil))
			return
		}

		t := res.Token

		var s model.Session
		err = s.GetByToken(t)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, util.SystemError.Msg(nil))
			return
		}

		if s.Id == 0 {
			c.AbortWithStatusJSON(http.StatusOK, util.TokenNotExists.Msg(nil))
			return
		}

		db.Db.First(&u, s.UserId)

		if u.Id == 0 || u.Status != model.StatusNormal {
			c.AbortWithStatusJSON(http.StatusOK, util.UserDisabled.Msg(nil))
			return
		}

		// 客户端和session不匹配
		if cl.Id != s.ClientId {
			c.AbortWithStatusJSON(http.StatusOK, util.SystemError.Msg(nil))
			return
		}

		// 过期检测
		if time.Now().After(s.ExpireTime) {
			c.AbortWithStatusJSON(http.StatusOK, util.SessionExpired.Msg(nil))
			return
		}

		c.Set(util.Session, &s)
		c.Set(util.User, &u)
	}
}
