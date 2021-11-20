package record

import (
	"github.com/gin-gonic/gin"
	"github.com/wuzehv/passport/model"
	"github.com/wuzehv/passport/service/db"
	"github.com/wuzehv/passport/util/static"
	"time"
)

type Result struct {
	Email      string
	ClientName string
	CreatedAt  time.Time
	IpAddr     string
	UserAgent  string
}

// @Description 登录日志
// @Tags 信息查看
// @Produce application/json
// @Param _ query validator.Pager false "_"
// @Success 200 {object} static.Response{data=record.Result}
// @Failure 200 {object} static.Response
// @Router /api/v1/records [GET]
func Index(c *gin.Context) {
	j := db.Db.Table("login_records r").Select("u.email, c.name client_name, r.created_at," +
		" r.ip_addr, r.user_agent, r.type").
		Joins("inner join users u on r.user_id = u.id").
		Joins("left join clients c on r.client_id = c.id")

	p := model.Param{
		JoinObj:    j,
		Bind:       make(map[string]interface{}),
		OrderField: "r.id",
	}

	res, err := model.PaginateContext(c, &p)
	if err != nil {
		c.JSON(static.SystemError.Msg(err))
		return
	}

	c.JSON(static.Success.Msg(res))
}
