package action

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wuzehv/passport/model"
	"github.com/wuzehv/passport/service/db"
	"github.com/wuzehv/passport/util"
	"net/http"
)

type Form struct {
	Url    string `form:"url" valid:"Required"`
	Remark string `form:"remark" valid:"Required"`
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
// @Param url body string true "路由"
// @Param remark body string true "备注"
// @Success 200 {object} util.Response
// @Failure 200 {object} util.Response
// @Router /api/v1/actions [POST]
func Add(c *gin.Context) {
	var json Form
	if err := c.ShouldBind(&json); err != nil {
		fmt.Println(json)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var action model.Action
	action.Url, action.Remark = json.Url, json.Remark
	db.Db.Save(&action)

	c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
}

// @Description 添加接口
// @Tags 后台系统
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param id query int true "ID"
// @Param url query string true "路由"
// @Param remark query string true "备注"
// @Success 200 {object} util.Response
// @Failure 200 {object} util.Response
// @Router /api/v1/actions/{id} [PUT]
func Update(c *gin.Context) {

}
