// 日志切割
package journal

import (
	"github.com/gin-gonic/gin"
	"github.com/wuzehv/passport/util/config"
	"github.com/wuzehv/passport/util/file"
	"io"
	"os"
	"time"
)

func init() {
	go func() {
		d := func() time.Duration {
			cur := time.Now()
			zero := time.Date(cur.Year(), cur.Month(), cur.Day()+1, 0, 0, 0, 0, cur.Location())
			return zero.Sub(cur)
		}

		for t := time.NewTimer(d()); ; t.Reset(d()) {
			select {
			case i := <-t.C:
				logDir := config.Log.Dir
				filename := i.Format(config.Log.Filename) + ".log"

				f, err := file.AppendOpenFile(logDir, filename)
				if err != nil {
					Error("logrotate", err)
					continue
				}

				Info("logrotate", filename)

				writer := []io.Writer{
					f,
				}

				if config.IsDev() {
					writer = append(writer, os.Stdout)
				}

				gin.DefaultWriter = io.MultiWriter(writer...)
			}
		}
	}()
}
