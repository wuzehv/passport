// 外部客户端接口

package svc

import (
	"github.com/gin-gonic/gin"
	"github.com/wuzehv/passport/model"
	"github.com/wuzehv/passport/service/db"
	"github.com/wuzehv/passport/service/rdb"
	"github.com/wuzehv/passport/util"
	"net/http"
	"time"
)

// @Description 客户端回调确认接口，更新session状态为已登录
// @Tags Svc接口
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param _ formData util.SvcRequest false "_"
// @Success 200 {object} util.Response
// @Failure 200 {object} util.Response
// @Router /svc/session [POST]
func Session(c *gin.Context) {
	tmp, _ := c.Get(util.Session)
	s := tmp.(*model.Session)

	if s.Status != model.StatusInit {
		c.AbortWithStatusJSON(http.StatusOK, util.SystemError.Msg(nil))
		return
	}

	// 更新session状态
	s.Status = model.StatusLogin
	db.Db.Save(&s)

	c.JSON(http.StatusOK, util.Success.Msg(nil))
}

// @Description 客户端业务代码执行之前，调用该接口获取用户信息
// @Tags Svc接口
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param _ formData util.SvcRequest false "_"
// @Success 200 {object} util.Response{data=model.User}
// @Failure 200 {object} util.Response
// @Router /svc/userinfo [POST]
func Userinfo(c *gin.Context) {
	tmp, _ := c.Get(util.Session)
	s := tmp.(*model.Session)

	// 登录状态
	if s.Status != model.StatusLogin {
		c.AbortWithStatusJSON(http.StatusOK, util.SessionStatusNotLogin.Msg(nil))
		return
	}

	tmp, _ = c.Get(util.User)
	u := tmp.(*model.User)

	rdb.SetJson(s.Token, u, time.Minute)

	c.JSON(http.StatusOK, util.Success.Msg(u))
}
