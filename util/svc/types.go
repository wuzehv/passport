package svc

import (
	"github.com/wuzehv/passport/model"
	"github.com/wuzehv/passport/util/journal"
)

type Adapter interface {
	Generate(userId, clientId uint) (string, error)
	Confirm(token string) error
	Valid(token string, user *model.User) error
	Destroy(userId uint) error
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
