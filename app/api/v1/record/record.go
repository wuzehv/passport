package record

import (
	"github.com/gin-gonic/gin"
	"github.com/wuzehv/passport/model"
	"github.com/wuzehv/passport/service/db"
	"github.com/wuzehv/passport/util/static"
	"strings"
	"time"
)

type Result struct {
	Id         uint
	Email      string
	ClientName string
	Realname   string
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
	userEmail := c.Query("email")
	client := c.Query("client")

	j := db.Db.Table("login_records r").Select("r.id, u.email, u.realname, c.name client_name, r.created_at," +
		" r.ip_addr, r.user_agent, r.type").
		Joins("inner join users u on r.user_id = u.id").
		Joins("left join clients c on r.client_id = c.id")

	p := model.Param{
		JoinObj:    j,
		Bind:       make(map[string]interface{}),
		OrderField: "r.id",
		OrderType:  c.Query(static.OrderType),
	}

	var where strings.Builder
	where.WriteString("1")

	if strings.TrimSpace(client) != "" {
		p.Bind["client_name"] = "%" + client + "%"
		where.WriteString(" and c.name like @client_name")
	}

	if strings.TrimSpace(userEmail) != "" {
		p.Bind["email"] = "%" + userEmail + "%"
		where.WriteString(" and u.email like @email")
	}

	p.Where = where.String()

	res, err := model.PaginateContext(c, &p)
	if err != nil {
		c.JSON(static.SystemError.Msg(err))
		return
	}

	c.JSON(static.Success.Msg(res))
}
