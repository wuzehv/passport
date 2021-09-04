// 平台内部接口

package sso

import (
	"github.com/gin-gonic/gin"
	"github.com/wuzehv/passport/model"
	"github.com/wuzehv/passport/service/db"
	"github.com/wuzehv/passport/util/common"
	"github.com/wuzehv/passport/util/config"
	"github.com/wuzehv/passport/util/static"
	"github.com/wuzehv/passport/util/svc"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Form struct {
	Username string `form:"username" binding:"required,email"`         // 用户邮箱
	Password string `form:"password" binding:"required,gte=6,lte=255"` // 密码
}

func Index(c *gin.Context) {
	tmp, _ := c.Get(static.Client)
	cl := tmp.(*model.Client)

	jump := c.GetString(static.Jump)
	uid := c.GetInt(static.Uid)

	if uid == 0 {
		c.HTML(http.StatusOK, "sso/login", gin.H{
			"domain": cl.Domain,
			"jump":   jump,
		})
		return
	}

	commonDeal(c, cl, uint(uid), jump)
}

// @Description 登录
// @Tags Sso入口
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param domain query string false "客户端标识"
// @Param _ formData Form false "_"
// @Success 200 {object} static.Response
// @Failure 400 {object} static.Response
// @Failure 403 {object} static.Response
// @Failure 500 {object} static.Response
// @Router /sso/login [POST]
func Login(c *gin.Context) {
	jump := c.GetString(static.Jump)

	var data Form
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, static.ParamsError.Msg(err))
		return
	}

	// 校验密码
	var u model.User
	err := u.GetByEmail(data.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, static.SystemError.Msg(err))
		return
	}

	tmp, _ := c.Get(static.Client)
	cl := tmp.(*model.Client)

	// 初始化登录信息
	r := model.LoginRecord{
		UserId:    u.Id,
		ClientId:  cl.Id,
		IpAddr:    c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
	}

	if model.FailNumOut() {
		r.Type = model.TypeOther
		if err = db.Db.Save(&r).Error; err != nil {
			c.JSON(http.StatusInternalServerError, static.SystemError.Msg(err))
			return
		}
		c.JSON(http.StatusForbidden, static.UsernamePasswdFailNumOut.Msg(nil))
		return
	}

	if !common.VerifyPassword(u.Password, data.Password) {
		r.Type = model.TypeFail
		if err = db.Db.Save(&r).Error; err != nil {
			c.JSON(http.StatusInternalServerError, static.SystemError.Msg(err))
			return
		}
		c.JSON(http.StatusForbidden, static.UsernamePasswdNotMatch.Msg(nil))
		return
	}

	// 初始化token
	token := common.GenToken() + strconv.FormatUint(uint64(u.Id), 10)
	u.Token = token
	exp, _ := time.Parse("2006-01-02 15:04:05", time.Now().Add(config.Svc.ExpireTime).Format("2006-01-02")+" 04:00:00")
	u.ExpireTime = exp
	if err = db.Db.Save(&u).Error; err != nil {
		c.JSON(http.StatusInternalServerError, static.SystemError.Msg(err))
		return
	}

	// 设置会话为浏览器关闭即失效
	c.SetCookie(static.CookieFlag, token, 0, "/", "", !config.IsDev(), true)

	// 重置所有客户端session状态
	err = model.LogoutAll(u.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, static.SystemError.Msg(err))
		return
	}

	r.Type = model.TypeSuccess
	if err = db.Db.Save(&r).Error; err != nil {
		c.JSON(http.StatusInternalServerError, static.SystemError.Msg(err))
		return
	}

	commonDeal(c, cl, u.Id, jump)
}

func commonDeal(c *gin.Context, cl *model.Client, userId uint, jump string) {
	callbackUrl, _ := url.Parse(cl.Callback)

	// 持久化
	adp := svc.New(config.Svc.Adapter)
	token, err := adp.Generate(userId, cl.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, static.SystemError.Msg(err))
		return
	}

	callbackParams := url.Values{}
	callbackParams.Add(static.Token, token)
	callbackParams.Add(static.Jump, jump)

	callbackUrl.RawQuery = callbackParams.Encode()

	isClient := c.GetBool(static.Sso)
	if isClient {
		c.HTML(http.StatusOK, "sso/redirect", gin.H{
			"callback": callbackUrl,
		})
	} else {
		// 如果是sso中心登录，跳转到首页
		c.Redirect(http.StatusMovedPermanently, "/api/v1/index")
		// todo 直接返回json
		//c.JSON(http.StatusOK, static.Success.Msg(token))
	}
}

// @Description 退出
// @Tags Sso入口
// @Produce application/json
// @Success 200 {object} static.Response
// @Failure 500 {object} static.Response
// @Router /sso/logout [POST]
func Logout(c *gin.Context) {
	uid := c.GetInt(static.Uid)

	adp := svc.New(config.Svc.Adapter)
	if err := adp.Destroy(uint(uid)); err != nil {
		c.JSON(http.StatusInternalServerError, static.SystemError.Msg(err))
		return
	}

	c.SetCookie(static.CookieFlag, "false", -1, "/", "", !config.IsDev(), true)
	c.HTML(http.StatusOK, "sso/logout", nil)
}
