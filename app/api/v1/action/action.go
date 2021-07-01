package action

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
	Url    string `form:"url" valid:"Required" binding:"required" minLength:"1" maxLength:"255"`    // 路径
	Remark string `form:"remark" valid:"Required" binding:"required" minLength:"1" maxLength:"255"` // 备注
}

// @Description 接口列表
// @Tags 接口管理
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param _ query validator.Pager false "_"
// @Success 200 {object} util.Response{data=model.Action}
// @Failure 500 {object} util.Response
// @Router /api/v1/actions [GET]
func Index(c *gin.Context) {
	var t model.Action
	res, err := model.PaginateContext(c, &model.Param{Table: &t})
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.SystemError.Msg(nil))
		return
	}

	c.JSON(http.StatusOK, util.Success.Msg(res))
}

// @Description 添加接口
// @Tags 接口管理
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param _ formData Form false "_"
// @Success 200 {object} util.Response
// @Failure 400 {object} util.Response
// @Failure 404 {object} util.Response
// @Failure 500 {object} util.Response
// @Router /api/v1/actions [POST]
func Add(c *gin.Context) {
	var data Form
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, util.ParamsError.Msg(err.Error()))
		return
	}

	var d model.Action
	d.Url, d.Remark = data.Url, data.Remark
	if err := db.Db.Save(&d).Error; err != nil {
		c.JSON(http.StatusInternalServerError, util.SystemError.Msg(nil))
		return
	}

	c.JSON(http.StatusOK, util.Success.Msg(nil))
}

// @Description 更新接口
// @Tags 接口管理
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param id path int true "ID"
// @Param _ formData Form false "_"
// @Success 200 {object} util.Response
// @Failure 400 {object} util.Response
// @Failure 404 {object} util.Response
// @Failure 500 {object} util.Response
// @Router /api/v1/actions/{id} [PUT]
func Update(c *gin.Context) {
	id := c.Param("id")
	var data Form
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, util.ParamsError.Msg(err.Error()))
		return
	}

	var d model.Action
	err := db.Db.First(&d, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, util.ParamsError.Msg(nil))
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, util.SystemError.Msg(nil))
		return
	}

	d.Url, d.Remark = data.Url, data.Remark
	if err := db.Db.Save(&d).Error; err != nil {
		c.JSON(http.StatusInternalServerError, util.SystemError.Msg(nil))
		return
	}

	c.JSON(http.StatusOK, util.Success.Msg(nil))
}
