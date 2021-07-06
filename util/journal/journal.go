// journal 自定义日志方法
// log、logger这两个包名太容易与其他包冲突，所以包名定义为journal
package journal

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

type logger struct {
	Time    time.Time   `json:"time"`
	Type    string      `json:"type"`
	Keyword string      `json:"keyword"`
	Data    interface{} `json:"data"`
	Stack   interface{} `json:"stack"`
}

const (
	TypeError = "error"
	TypeInfo  = "info"
	TypeDebug = "debug"
)

func New(logType, keyword string, data interface{}) *logger {
	return &logger{
		Time:    time.Now(),
		Type:    logType,
		Keyword: keyword,
		Data:    data,
	}
}

func (l *logger) String() string {
	s, _ := json.Marshal(l)
	return string(s)
}

func (l *logger) log() {
	if _, err := fmt.Fprintf(gin.DefaultWriter, "%s\n", l); err != nil {
		log.Printf("journal error: %v\n", err)
	}
}

func Info(keyword string, data interface{}) {
	New(TypeInfo, keyword, data).log()
}

func Error(keyword string, data interface{}) {
	New(TypeError, keyword, data).log()
}

func Debug(keyword string, data interface{}) {
	New(TypeDebug, keyword, data).log()
}
