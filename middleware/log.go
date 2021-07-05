package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wuzehv/passport/util/config"
	"io"
	"log"
	"os"
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

func init() {
	gin.DisableConsoleColor()

	logDir := config.Log.Dir
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err = os.MkdirAll(logDir, 0777); err != nil {
			fmt.Println(logDir)
			log.Fatalf("log dir create error: %v\n", err)
		}
	}

	finalFile := logDir + "/gin.log"

	f, err := os.OpenFile(finalFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("log file create error: %v\n", err)
	}

	writer := []io.Writer{
		f,
	}

	if config.IsDev() {
		writer = append(writer, os.Stdout)
	}

	gin.DefaultWriter = io.MultiWriter(writer...)
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
