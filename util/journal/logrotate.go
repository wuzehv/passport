// 日志切割
package journal

import (
	"time"
)

func logrotate() {
	go func() {
		d := func() time.Duration {
			cur := time.Now()
			zero := time.Date(cur.Year(), cur.Month(), cur.Day()+1, 0, 0, 0, 0, cur.Location())

			Debug("logrotate_date", cur.String()+" ~ "+zero.String())
			return zero.Sub(cur)
		}

		for t := time.NewTimer(d()); ; t.Reset(d()) {
			select {
			case i := <-t.C:
				filename, err := initLogWriter(i)
				if err != nil {
					Error("logrotate", err)
					continue
				}

				Info("logrotate", filename)
			}
		}
	}()
}
