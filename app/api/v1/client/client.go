package client

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/wuzehv/passport/model"
	"github.com/wuzehv/passport/service/db"
	"github.com/wuzehv/passport/util"
	"gorm.io/gorm"
	"net/http"
)

type Form struct {
	Domain   string `form:"domain" valid:"Required" binding:"required" minLength:"1" maxLength:"255"`   // 域名
	Callback string `form:"callback" valid:"Required" binding:"required" minLength:"1" maxLength:"255"` // 回调地址
	Secret   string `form:"secret" valid:"Required" binding:"required" minLength:"1" maxLength:"255"`   // 用来签名校验的密钥
}

// @Description 客户端列表
// @Tags 客户端管理
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param _ query validator.Pager false "_"
// @Success 200 {object} util.Response{data=model.Client}
// @Failure 500 {object} util.Response
// @Router /api/v1/clients [GET]
func Index(c *gin.Context) {
	var t model.Client
	res, err := model.PaginateContext(c, &model.Param{Table: &t})
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.SystemError.Msg(nil))
		return
	}

	c.JSON(http.StatusOK, util.Success.Msg(res))
}

// @Description 添加客户端
// @Tags 客户端管理
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param _ formData Form false "_"
// @Success 200 {object} util.Response
// @Failure 400 {object} util.Response
// @Failure 404 {object} util.Response
// @Failure 500 {object} util.Response
// @Router /api/v1/clients [POST]
func Add(c *gin.Context) {
	var data Form
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, util.ParamsError.Msg(err.Error()))
		return
	}

	d := model.Client{
		Domain:   data.Domain,
		Callback: data.Callback,
		Secret:   data.Secret,
		Status:   model.StatusNormal,
	}
	if err := db.Db.Save(&d).Error; err != nil {
		c.JSON(http.StatusInternalServerError, util.SystemError.Msg(nil))
		return
	}

	c.JSON(http.StatusOK, util.Success.Msg(nil))
}

// @Description 客户端更新
// @Tags 客户端管理
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param id path int true "ID"
// @Param _ formData Form false "_"
// @Success 200 {object} util.Response
// @Failure 400 {object} util.Response
// @Failure 404 {object} util.Response
// @Failure 500 {object} util.Response
// @Router /api/v1/clients/{id} [PUT]
func Update(c *gin.Context) {
	id := c.Param("id")
	var data Form
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, util.ParamsError.Msg(err.Error()))
		return
	}

	var d model.Client
	err := db.Db.First(&d, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, util.ParamsError.Msg(nil))
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, util.SystemError.Msg(nil))
		return
	}

	d.Domain = data.Domain
	d.Callback = data.Callback
	d.Secret = data.Secret
	if err := db.Db.Save(&d).Error; err != nil {
		c.JSON(http.StatusInternalServerError, util.SystemError.Msg(nil))
		return
	}

	c.JSON(http.StatusOK, util.Success.Msg(nil))
}

// @Description 客户端启用/禁用
// @Tags 客户端管理
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param id path int true "ID"
// @Param disabled query bool true "启用：false 禁用：true"
// @Success 200 {object} util.Response
// @Failure 400 {object} util.Response
// @Failure 404 {object} util.Response
// @Failure 500 {object} util.Response
// @Router /api/v1/clients/{id} [PATCH]
func Disable(c *gin.Context) {
	id := c.Param("id")

	var d model.Client
	m := db.Db.First(&d, id)
	err := m.Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, util.ParamsError.Msg(nil))
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, util.SystemError.Msg(nil))
		return
	}

	disabled := c.Query("disabled")
	status := model.StatusDisabled
	if disabled == "false" {
		status = model.StatusNormal
	}

	if uint(status) == d.Status {
		c.JSON(http.StatusOK, util.Success.Msg(nil))
		return
	}

	if err := m.Update("status", status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, util.SystemError.Msg(nil))
		return
	}

	c.JSON(http.StatusOK, util.Success.Msg(nil))
}
