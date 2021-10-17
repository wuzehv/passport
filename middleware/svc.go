package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/wuzehv/passport/model"
	"github.com/wuzehv/passport/service/rdb"
	"github.com/wuzehv/passport/util/common"
	"github.com/wuzehv/passport/util/static"
	"net/url"
)

// Svc svc调用入口，校验token
func Svc() gin.HandlerFunc {
	return func(c *gin.Context) {
		var res static.SvcRequest
		if c.ShouldBind(&res) != nil {
			c.AbortWithStatusJSON(static.ParamsError.Msg(nil))
			return
		}

		var u model.User
		if rdb.GetJson(res.Token, &u) {
			c.AbortWithStatusJSON(static.Success.Msg(u))
			return
		}

		domain, err := url.QueryUnescape(res.Domain)
		if err != nil {
			c.AbortWithStatusJSON(static.SystemError.Msg(err))
			return
		}

		var cl model.Client
		err = cl.GetByDomain(domain)
		if err != nil {
			c.AbortWithStatusJSON(static.SystemError.Msg(err))
			return
		}

		if cl.Id == 0 || cl.Status != model.StatusNormal {
			c.AbortWithStatusJSON(static.ClientDisabled.Msg(nil))
			return
		}

		m := make(map[string]string)
		m[static.Token] = res.Token
		m[static.Timestamp] = res.Timestamp
		m[static.Domain] = res.Domain
		if common.GenSign(m, cl.Secret) != res.Sign {
			c.AbortWithStatusJSON(static.SignatureError.Msg(nil))
			return
		}

		c.Set(static.Token, res.Token)
	}
}
