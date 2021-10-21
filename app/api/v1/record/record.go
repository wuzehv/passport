package record

import (
	"github.com/gin-gonic/gin"
	"github.com/wuzehv/passport/model"
	"github.com/wuzehv/passport/util/static"
	"strings"
)

func Index(c *gin.Context) {
	var t model.LoginRecord

	client := c.Query("client")
	userEmail := c.Query("email")

	p := model.Param{
		Table:      &t,
		Bind:       make(map[string]interface{}),
	}

	var where strings.Builder
	where.WriteString("1")

	if strings.TrimSpace(client) != "" {
		p.Bind["id"] = id
		where.WriteString(" and id = @id")
	}

	if strings.TrimSpace(userEmail) != "" {
		p.Bind["email"] = "%" + userEmail + "%"
		where.WriteString(" and email like @email")
	}

	p.Where = where.String()

	res, err := model.PaginateContext(c, &p)
	if err != nil {
		c.JSON(static.SystemError.Msg(err))
		return
	}

	c.JSON(static.Success.Msg(res))
}
