package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/wuzehv/passport/util/static"
	"golang.org/x/time/rate"
	"net/http"
)

var limiter *rate.Limiter

func init() {
	limiter = rate.NewLimiter(10, 5)
}

// 访问频率控制
func Rate() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusForbidden, static.Success.Msg(nil))
			return
		}
	}
}
