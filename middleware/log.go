package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"time"
)

type LogFormat struct {
	ClientIp   string      `json:"client_ip"`
	Timestamp  string      `json:"timestamp"`
	Method     string      `json:"method"`
	Path       string      `json:"path"`
	Code       int         `json:"code"`
	UserAgent  string      `json:"user_agent"`
	Message    string      `json:"message"`
	BodySize   int         `json:"body_size"`
	Form       interface{} `json:"form"`
	Host       string      `json:"host"`
	RemoteAddr string      `json:"remote_addr"`
}

func Log() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		l := &LogFormat{
			ClientIp:   param.ClientIP,
			Timestamp:  param.TimeStamp.Format(time.RFC1123),
			Method:     param.Method,
			Path:       param.Path,
			Code:       param.StatusCode,
			UserAgent:  param.Request.UserAgent(),
			Host:       param.Request.Host,
			RemoteAddr: param.Request.RemoteAddr,
			BodySize:   param.BodySize,
			Form:       param.Request.PostForm,
			Message:    param.ErrorMessage,
		}

		res, _ := json.Marshal(l)
		return string(res) + "\n"
	})
}
