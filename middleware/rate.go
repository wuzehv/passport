package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/wuzehv/passport/util/config"
	"github.com/wuzehv/passport/util/static"
	"golang.org/x/time/rate"
)

var limiter *rate.Limiter

func init() {
	limiter = rate.NewLimiter(rate.Limit(config.Rate.Limit), config.Rate.Period)
}

// 访问频率控制
func Rate() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(static.Forbidden.Msg(nil))
			return
		}
	}
}
