package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/wuzehv/passport/util"
	"net/http"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Allow-Methods", "*")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatusJSON(http.StatusNoContent, util.Success.Msg(nil))
			return
		}
	}
}
