// 平台内部接口

package sso

import (
	"github.com/gin-gonic/gin"
	"github.com/wuzehv/passport/model"
	"github.com/wuzehv/passport/service/db"
	"github.com/wuzehv/passport/util"
	"github.com/wuzehv/passport/util/config"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func Index(c *gin.Context) {
	tmp, _ := c.Get(util.Client)
	cl := tmp.(*model.Client)

	jump := c.GetString(util.Jump)

	uid := c.GetInt(util.Uid)

	if uid == 0 {
		c.HTML(http.StatusOK, "sso/login", gin.H{
			"domain": cl.Domain,
			"jump":   jump,
		})
		return
	}

	commonDeal(c, cl, uint(uid), jump)
}

func Login(c *gin.Context) {
	jump := c.GetString(util.Jump)

	name := c.PostForm("username")
	passwd := c.PostForm("password")

	// 校验密码
	var u model.User
	err := u.GetByEmail(name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, util.SystemError.Msg(nil))
		return
	}

	tmp, _ := c.Get(util.Client)
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
		db.Db.Save(&r)
		c.AbortWithStatusJSON(http.StatusOK, util.UsernamePasswdFailNumOut.Msg(nil))
		return
	}

	if !util.VerifyPassword(u.Password, passwd) {
		r.Type = model.TypeFail
		db.Db.Save(&r)
		c.AbortWithStatusJSON(http.StatusOK, util.UsernamePasswdNotMatch.Msg(nil))
		return
	}

	// 初始化token
	token := util.GenToken() + strconv.FormatUint(uint64(u.Id), 10)
	u.Token = token
	exp, _ := time.Parse("2006-01-02 15:04:05", time.Now().Add(model.ExpireTime).Format("2006-01-02")+" 04:00:00")
	u.ExpireTime = exp
	db.Db.Save(&u)
	// 设置会话为浏览器关闭即失效
	c.SetCookie(util.CookieFlag, token, 0, "/", "", !config.IsDev(), true)

	// 重置所有客户端session状态
	err = model.LogoutAll(u.Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, util.SystemError.Msg(nil))
		return
	}

	r.Type = model.TypeSuccess
	db.Db.Save(&r)

	commonDeal(c, cl, u.Id, jump)
}

func commonDeal(c *gin.Context, cl *model.Client, userId uint, jump string) {
	callbackUrl, _ := url.Parse(cl.Callback)

	// 持久化
	s := model.NewSession(userId, cl.Id)

	callbackParams := url.Values{}
	callbackParams.Add(util.Token, s.Token)
	callbackParams.Add(util.Jump, jump)

	callbackUrl.RawQuery = callbackParams.Encode()

	isSso := c.GetBool(util.Sso)

	if isSso {
		c.HTML(http.StatusOK, "sso/redirect", gin.H{
			"callback": callbackUrl,
		})
	} else {
		// 如果不是sso，跳转到首页
		c.Redirect(http.StatusMovedPermanently, "/api/v1/index")
	}
}

func Logout(c *gin.Context) {
	uid := c.GetInt(util.Uid)
	model.LogoutAll(uint(uid))

	c.SetCookie(util.CookieFlag, "false", -1, "/", "", !config.IsDev(), true)
	c.HTML(http.StatusOK, "sso/logout", gin.H{})
}
