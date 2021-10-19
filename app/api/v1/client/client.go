package client

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/wuzehv/passport/model"
	"github.com/wuzehv/passport/service/db"
	"github.com/wuzehv/passport/util/static"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

type Form struct {
	Name     string `form:"name" binding:"required,gte=3,lte=255"`   // 域名
	Domain   string `form:"domain" binding:"required,unix_addr"`     // 域名
	Callback string `form:"callback" binding:"required,url"`         // 回调地址
	Secret   string `form:"secret" binding:"required,gte=6,lte=255"` // 用来签名校验的密钥
}

// @Description 客户端列表
// @Tags 客户端管理
// @Produce application/json
// @Param _ query validator.Pager false "_"
// @Success 200 {object} static.Response{data=model.Client}
// @Failure 500 {object} static.Response
// @Router /api/v1/clients [GET]
func Index(c *gin.Context) {
	var t model.Client

	id := c.Query("id")
	domain := c.Query("domain")

	p := model.Param{
		Table:      &t,
		Bind:       make(map[string]interface{}),
		OrderField: c.Query(static.OrderField),
		OrderType:  c.Query(static.OrderType),
	}

	var where strings.Builder
	where.WriteString("1")

	if strings.TrimSpace(id) != "" {
		p.Bind["id"] = id
		where.WriteString(" and id = @id")
	}

	if strings.TrimSpace(domain) != "" {
		p.Bind["domain"] = "%" + domain + "%"
		where.WriteString(" and domain like @domain")
	}

	p.Where = where.String()

	res, err := model.PaginateContext(c, &p)
	if err != nil {
		c.JSON(static.SystemError.Msg(err))
		return
	}

	c.JSON(static.Success.Msg(res))
}

// @Description 添加客户端
// @Tags 客户端管理
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param _ formData Form false "_"
// @Success 200 {object} static.Response
// @Failure 400 {object} static.Response
// @Failure 404 {object} static.Response
// @Failure 500 {object} static.Response
// @Router /api/v1/clients [POST]
func Add(c *gin.Context) {
	var data Form
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(static.ParamsError.Msg(err))
		return
	}

	d := model.Client{
		Name:     data.Name,
		Domain:   data.Domain,
		Callback: data.Callback,
		Secret:   data.Secret,
		Status:   model.StatusNormal,
	}
	if err := db.Db.Save(&d).Error; err != nil {
		c.JSON(static.SystemError.Msg(err))
		return
	}

	c.JSON(static.Success.Msg(nil))
}

// @Description 客户端更新
// @Tags 客户端管理
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param id path int true "ID"
// @Param _ formData Form false "_"
// @Success 200 {object} static.Response
// @Failure 400 {object} static.Response
// @Failure 404 {object} static.Response
// @Failure 500 {object} static.Response
// @Router /api/v1/clients/{id} [PUT]
func Update(c *gin.Context) {
	id := c.Param("id")
	var data Form
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(static.ParamsError.Msg(err))
		return
	}

	var d model.Client
	err := db.Db.First(&d, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(static.ParamsError.Msg(nil))
		return
	}

	if err != nil {
		c.JSON(static.SystemError.Msg(err))
		return
	}

	d.Name = data.Name
	d.Domain = data.Domain
	d.Callback = data.Callback
	d.Secret = data.Secret
	if err := db.Db.Save(&d).Error; err != nil {
		c.JSON(static.SystemError.Msg(err))
		return
	}

	c.JSON(static.Success.Msg(nil))
}

// @Description 客户端启用/禁用
// @Tags 客户端管理
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param id path int true "ID"
// @Success 200 {object} static.Response
// @Failure 400 {object} static.Response
// @Failure 404 {object} static.Response
// @Failure 500 {object} static.Response
// @Router /api/v1/clients/{id}/toggle-status [POST]
func ToggleStatus(c *gin.Context) {
	id := c.Param("id")

	var d model.Client
	m := db.Db.First(&d, id)
	err := m.Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(static.ParamsError.Msg(nil))
		return
	}

	if err != nil {
		c.JSON(static.SystemError.Msg(err))
		return
	}

	status := model.StatusNormal
	if uint(status) == d.Status {
		status = model.StatusDisabled
	}

	if err := m.Update("status", status).Error; err != nil {
		c.JSON(static.SystemError.Msg(err))
		return
	}

	c.JSON(static.Success.Msg(nil))
}

// @Description 检测客户端地址
// @Tags 客户端管理
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param url path string true "回调地址"
// @Success 200 {object} static.Response
// @Failure 400 {object} static.Response
// @Failure 404 {object} static.Response
// @Failure 500 {object} static.Response
// @Router /api/v1/clients/check-callback [HEAD]
func CheckCallback(c *gin.Context) {
	url := c.Query("url")
	res, err := http.Head(url)
	if err != nil {
		c.JSON(static.ParamsError.Msg(err))
		return
	}

	defer res.Body.Close()

	c.JSON(static.Success.Msg(nil))
}
