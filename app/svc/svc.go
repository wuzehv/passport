// 外部客户端接口

package svc

import (
	"github.com/gin-gonic/gin"
	"github.com/wuzehv/passport/model"
	"github.com/wuzehv/passport/util/config"
	"github.com/wuzehv/passport/util/static"
	"github.com/wuzehv/passport/util/svc"
	"net/http"
)

// @Description 客户端回调确认接口，更新session状态为已登录
// @Tags Svc接口
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param _ formData static.SvcRequest false "_"
// @Success 200 {object} static.Response
// @Failure 200 {object} static.Response
// @Router /svc/session [POST]
func Session(c *gin.Context) {
	//tmp, _ := c.Get(static.Session)
	//s := tmp.(*model.Session)
	//
	//if s.Status != model.StatusInit {
	//	c.AbortWithStatusJSON(http.StatusOK, static.SystemError.Msg(nil))
	//	return
	//}

	// 更新session状态
	//s.Status = model.StatusLogin
	//db.Db.Save(&s)

	c.JSON(http.StatusOK, static.Success.Msg(nil))
}

// @Description 客户端业务代码执行之前，调用该接口获取用户信息
// @Tags Svc接口
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param _ formData static.SvcRequest false "_"
// @Success 200 {object} static.Response{data=model.User}
// @Failure 200 {object} static.Response
// @Router /svc/userinfo [POST]
func Userinfo(c *gin.Context) {
	u := new(model.User)
	token := c.GetString(static.Token)

	adp := svc.New(config.Svc.Adapter)
	err := adp.ValidToken(token, u)
	if err == nil {
		c.JSON(http.StatusOK, static.Success.Msg(u))
		return
	}

	switch err.(static.Code) {
	case static.SessionNotExists:
		c.JSON(http.StatusNotFound, static.SessionNotExists.Msg(nil))
	case static.SessionExpired:
		c.JSON(http.StatusForbidden, static.SessionExpired.Msg(nil))
	default:
		c.JSON(http.StatusInternalServerError, static.SystemError.Msg(err))
	}
}
