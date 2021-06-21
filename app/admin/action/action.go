package action

import (
	"github.com/gin-gonic/gin"
	"github.com/wuzehv/passport/model/action"
	"github.com/wuzehv/passport/model/base"
	"github.com/wuzehv/passport/util"
	"net/http"
)

// @description 接口列表
// @tags 后台系统
// @accept application/x-www-form-urlencoded
// @produce application/json
// @param page query int false "页码"
// @param page_size query int false "每页数量"
// @success 200 {object} util.Response
// @failure 200 {object} util.Response
// @router /admin/action/index [GET]
func Index(c *gin.Context) {
	var t action.Action
	res, err := base.Paginate2(c, &base.Param{Table: &t})
	if err != nil {
		c.JSON(http.StatusOK, util.SystemError.Msg(nil))
		return
	}

	c.JSON(http.StatusOK, util.Success.Msg(res))
}

// @description 添加接口
// @tags 后台系统
// @Accept application/x-www-form-urlencoded
// @produce application/json
// @param url query string true "路由"
// @param remark query string true "备注"
// @success 200 {object} util.Response
// @failure 200 {object} util.Response
// @router /admin/action/add [POST]
func Add(c *gin.Context) {

}
