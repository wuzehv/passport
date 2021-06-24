package index

import (
	"github.com/gin-gonic/gin"
	"github.com/wuzehv/passport/model"
	"github.com/wuzehv/passport/service/db"
	"net/http"
)

func Index(c *gin.Context) {
	var u []model.User
	db.Db.Find(&u)

	var cl []model.Client
	db.Db.Find(&cl)
	c.HTML(http.StatusOK, "api/index/index", gin.H{
		"users":   u,
		"clients": cl,
		"login":   true,
	})
}
