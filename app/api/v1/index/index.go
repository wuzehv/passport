package index

import (
	"github.com/gin-gonic/gin"
	"github.com/wuzehv/passport/model"
	"github.com/wuzehv/passport/service/db"
	"github.com/wuzehv/passport/util/static"
	"net/http"
)

func Index(c *gin.Context) {
	var userNum, clientNum, sessionNum, recodeNum int64

	db.Db.Model(&model.User{}).Count(&userNum)
	db.Db.Model(&model.Client{}).Count(&clientNum)
	db.Db.Model(&model.Session{}).Where("status = ?", model.StatusLogin).Count(&sessionNum)
	db.Db.Model(&model.LoginRecord{}).Count(&recodeNum)

	u, _ := c.Get(static.User)
	user := u.(*model.User)

	c.HTML(http.StatusOK, "index", gin.H{
		"user_num":    userNum,
		"client_num":  clientNum,
		"session_num": sessionNum,
		"record_num":  recodeNum,
		"user_name":   user.Realname,
	})
}
