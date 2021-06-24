package action

import (
	"github.com/gin-gonic/gin"
	"github.com/wuzehv/passport/model"
	"github.com/wuzehv/passport/util"
	"net/http"
)

// @Description 接口列表
// @Tags 后台系统
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} util.Response
// @Failure 200 {object} util.Response
// @Router /admin/action/index [GET]
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
// @Param url query string true "路由"
// @Param remark query string true "备注"
// @Success 200 {object} util.Response
// @Failure 200 {object} util.Response
// @Router /admin/action/add [POST]
func Add(c *gin.Context) {

}
