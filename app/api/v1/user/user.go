package user

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
	Email    string `form:"email" valid:"Required" binding:"required" minLength:"1" maxLength:"255"`                    // 用户邮箱
	RealName string `form:"realname" json:"realname" valid:"Required" binding:"required" minLength:"1" maxLength:"255"` // 真实姓名
	Gender   int    `form:"gender" valid:"Required" binding:"required" minimum:"1" maximum:"2" default:"1"`             // 性别
	Mobile   string `form:"mobile" valid:"Required" binding:"required" minLength:"1" maxLength:"255"`                   // 手机号
	Password string `form:"password" valid:"Required" binding:"required" minLength:"1" maxLength:"255"`                 // 密码
}

// @Description 用户列表
// @Tags 用户管理
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param _ query validator.Pager false "_"
// @Success 200 {object} util.Response{data=model.User}
// @Failure 500 {object} util.Response
// @Router /api/v1/users [GET]
func Index(c *gin.Context) {
	var t model.User
	res, err := model.PaginateContext(c, &model.Param{Table: &t})
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.SystemError.Msg(nil))
		return
	}

	c.JSON(http.StatusOK, util.Success.Msg(res))
}

// @Description 添加用户
// @Tags 用户管理
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param _ formData Form false "_"
// @Success 200 {object} util.Response
// @Failure 400 {object} util.Response
// @Failure 404 {object} util.Response
// @Failure 500 {object} util.Response
// @Router /api/v1/users [POST]
func Add(c *gin.Context) {
	var data Form
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, util.ParamsError.Msg(err.Error()))
		return
	}

	d := model.User{
		Email:    data.Email,
		Password: util.GenPassword(data.Password),
		Gender:   data.Gender,
		Mobile:   data.Mobile,
		Realname: data.RealName,
		Status:   model.StatusNormal,
	}
	if err := db.Db.Save(&d).Error; err != nil {
		c.JSON(http.StatusInternalServerError, util.SystemError.Msg(nil))
		return
	}

	c.JSON(http.StatusOK, util.Success.Msg(nil))
}

// @Description 用户更新
// @Tags 用户管理
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param id path int true "ID"
// @Param _ formData Form false "_"
// @Success 200 {object} util.Response
// @Failure 400 {object} util.Response
// @Failure 404 {object} util.Response
// @Failure 500 {object} util.Response
// @Router /api/v1/users/{id} [PUT]
func Update(c *gin.Context) {
	id := c.Param("id")
	var data Form
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, util.ParamsError.Msg(err.Error()))
		return
	}

	var d model.User
	err := db.Db.First(&d, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, util.ParamsError.Msg(nil))
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, util.SystemError.Msg(nil))
		return
	}

	d.Email = data.Email
	d.Realname = data.RealName
	d.Mobile = data.Mobile
	d.Gender = data.Gender
	d.Password = util.GenPassword(data.Password)

	if err := db.Db.Save(&d).Error; err != nil {
		c.JSON(http.StatusInternalServerError, util.SystemError.Msg(nil))
		return
	}

	c.JSON(http.StatusOK, util.Success.Msg(nil))
}

// @Description 用户启用/禁用
// @Tags 用户管理
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param id path int true "ID"
// @Param disabled query bool true "启用：false 禁用：true"
// @Success 200 {object} util.Response
// @Failure 400 {object} util.Response
// @Failure 404 {object} util.Response
// @Failure 500 {object} util.Response
// @Router /api/v1/users/{id} [PATCH]
func Disable(c *gin.Context) {
	id := c.Param("id")

	var d model.User
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
