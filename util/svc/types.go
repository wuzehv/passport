package svc

import (
	"github.com/wuzehv/passport/model"
	"github.com/wuzehv/passport/util/journal"
)

type Adapter interface {
	GenToken(userId, clientId uint) (string, error)
	ValidToken(token string, user *model.User) error
}

type data struct {
	userId, clientId uint
}

func New(adp string) Adapter {
	switch adp {
	case "jwt":
		return &j{}
	case "mysql":
		return &mysql{}
	default:
		journal.Error("svc_adapter", "adapter config error, use default mysql")
		return &mysql{}
	}
}
