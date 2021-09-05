package goemail

import (
	"github.com/jordan-wright/email"
	"github.com/wuzehv/passport/util/config"
	"github.com/wuzehv/passport/util/journal"
	"net/smtp"
	"sync"
	"time"
)

func Send(e *email.Email) error {
	return e.Send(config.Email.Address, smtp.PlainAuth("", config.Email.UserName, config.Email.Password, config.Email.Host))
}

func SendPool(ch <-chan *email.Email) error {
	p, err := email.NewPool(config.Email.Address, config.Email.PoolSize, smtp.PlainAuth("", config.Email.UserName, config.Email.Password, config.Email.Host))
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for e := range ch {
		wg.Add(1)
		go func(e *email.Email) {
			err := p.Send(e, 10*time.Second)
			if err != nil {
				journal.Error("send main error", err)
			}
			wg.Done()
		}(e)
	}

	wg.Wait()
	return nil
}
