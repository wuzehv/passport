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

type AddForm struct {
	Url    string `form:"url" valid:"Required"`
	Remark string `form:"remark" valid:"Required"`
}

type UpdateForm struct {
	Id int `form:"id" valid:"Required"`
	AddForm
}

// @Description 接口列表
// @Tags 后台系统
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} util.Response
// @Failure 200 {object} util.Response
// @Router /api/v1/actions [GET]
func Index(c *gin.Context) {
	var t model.Action
	res, err := model.PaginateContext(c, &model.Param{Table: &t})
	if err != nil {
		c.JSON(http.StatusOK, util.SystemError.Msg(nil))
		return
	}

	c.JSON(http.StatusOK, util.Success.Msg(res))
}

// @Description 添加接口
// @Tags 后台系统
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param url formData string true "路由"
// @Param remark formData string true "备注"
// @Success 200 {object} util.Response
// @Failure 400 {object} util.Response
// @Failure 404 {object} util.Response
// @Failure 500 {object} util.Response
// @Router /api/v1/actions [POST]
func Add(c *gin.Context) {
	var data AddForm
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, util.ParamsError.Msg(err.Error()))
		return
	}

	var action model.Action
	action.Url, action.Remark = data.Url, data.Remark
	if err := db.Db.Save(&action).Error; err != nil {
		c.JSON(http.StatusInternalServerError, util.SystemError.Msg(nil))
		return
	}

	c.JSON(http.StatusOK, util.Success.Msg(nil))
}

// @Description 更新接口
// @Tags 后台系统
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param id formData int true "ID"
// @Param url formData string true "路由"
// @Param remark formData string true "备注"
// @Success 200 {object} util.Response
// @Failure 400 {object} util.Response
// @Failure 404 {object} util.Response
// @Failure 500 {object} util.Response
// @Router /api/v1/actions/{id} [PUT]
func Update(c *gin.Context) {
	var data UpdateForm
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, util.ParamsError.Msg(err.Error()))
		return
	}

	var action model.Action
	err := db.Db.First(&action, data.Id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, util.ParamsError.Msg(nil))
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, util.SystemError.Msg(nil))
		return
	}

	action.Url, action.Remark = data.Url, data.Remark
	if err := db.Db.Save(&action).Error; err != nil {
		c.JSON(http.StatusInternalServerError, util.SystemError.Msg(nil))
		return
	}

	c.JSON(http.StatusOK, util.Success.Msg(nil))
}
