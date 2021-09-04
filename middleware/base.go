package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/wuzehv/passport/util/journal"
	"time"
)

type Access struct {
	LogFormat
	Duration string `json:"duration"`
}

// Base 统计执行时间
func Base() gin.HandlerFunc {
	return func(c *gin.Context) {
		s := time.Now()
		c.Next()
		l := Access{
			LogFormat{
				Path: c.Request.URL.Path,
			},
			time.Now().Sub(s).String(),
		}

		journal.Info("duration", l)
	}
}
