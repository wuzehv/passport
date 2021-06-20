package action

import (
	"github.com/gin-gonic/gin"
	"github.com/wuzehv/passport/model/action"
	"github.com/wuzehv/passport/model/base"
	"github.com/wuzehv/passport/util"
	"net/http"
)

func Index(c *gin.Context) {
	var t action.Action
	res, err := base.Paginate2(c, &base.Param{Table: &t})
	if err != nil {
		c.JSON(http.StatusOK, util.SystemError.Msg(nil))
		return
	}

	c.JSON(http.StatusOK, util.Success.Msg(res))
}
