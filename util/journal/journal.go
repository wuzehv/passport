// journal 自定义日志方法
// log、logger这两个包名太容易与其他包冲突，所以包名定义为journal
package journal

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/wuzehv/passport/util/config"
	"github.com/wuzehv/passport/util/file"
	"io"
	"log"
	"os"
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

// 内部log句柄，使用内置的log包，避免潜在的并发风险
var innerLog *log.Logger

func init() {
	gin.DisableConsoleColor()
	if _, err := initLogWriter(time.Now()); err != nil {
		log.Fatalf("log init error: %v\n", err)
	}
	logrotate()
}

func initLogWriter(t time.Time) (string, error) {
	logDir := config.Log.Dir
	filename := t.Format(config.Log.Filename) + ".log"

	f, err := file.AppendOpenFile(logDir, filename)
	if err != nil {
		return "", err
	}

	writer := []io.Writer{
		f,
	}

	if config.IsDev() {
		writer = append(writer, os.Stdout)
	}

	gin.DefaultWriter = io.MultiWriter(writer...)

	innerLog = log.New(gin.DefaultWriter, "", 0)

	return filename, nil
}

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
	if e, ok := l.Data.(error); ok {
		l.Data = e.Error()
	}

	innerLog.Println(l)
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
