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
	Time  time.Time   `json:"time"`
	Type  string      `json:"type"`
	Data  interface{} `json:"data"`
	Stack interface{} `json:"stack"`
}

const (
	TypeError = "error"
	TypeInfo  = "info"
	TypeDebug = "debug"
)

func New(lt string, data interface{}) *logger {
	return &logger{
		Time: time.Now(),
		Type: lt,
		Data: data,
	}
}

func (l *logger) String() string {
	s, _ := json.Marshal(l)
	return string(s)
}

func (l *logger) log() {
	if _, err := fmt.Fprintf(gin.DefaultWriter, "%s\n", l); err != nil {
		log.Fatalf("journal error: %v\n", err)
	}
}

func Info(data interface{}) {
	New(TypeInfo, data).log()
}

func Error(data interface{}) {
	New(TypeError, data).log()
}

func Debug(data interface{}) {
	New(TypeDebug, data).log()
}
