package RouterMidWeb

import (
	"fmt"
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	BasePedometer "gitee.com/weeekj/weeekj_core/v5/base/pedometer"
	BaseSafe "gitee.com/weeekj/weeekj_core/v5/base/safe"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreRPCX "gitee.com/weeekj/weeekj_core/v5/core/rpcx"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	RouterMidAPI "gitee.com/weeekj/weeekj_core/v5/router/mid/api"
	RouterReport "gitee.com/weeekj/weeekj_core/v5/router/report"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/gin-gonic/gin"
)

// HeaderBaseData 设定无头信息的请求
// 用于一些特殊页面，如错误页面等
func HeaderBaseData(c *gin.Context) {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("router mid, ", r)
		}
	}()
	//总的中间件
	RouterMid(c)
	//对ip进行安全检查
	SafetyIPON, err := BaseConfig.GetDataBool("SafetyIPON")
	if err != nil {
		SafetyIPON = true
	}
	ipAddr := c.ClientIP()
	if !CoreFilter.CheckIP(ipAddr) {
		RouterReport.WarnLog(c, "ip is error, "+ipAddr, nil, "ip_error", "IP错误")
		return
	}
	//如果关闭debug / 启动了IP安全检查 / 计数器超出限制
	if !Router2SystemConfig.Debug && SafetyIPON && BasePedometer.CheckData(&CoreRPCX.ArgsFrom{
		From: CoreSQLFrom.FieldsFrom{System: "safe-ip", Mark: ipAddr},
	}) {
		BaseSafe.CreateLog(&BaseSafe.ArgsCreateLog{
			System: "api.token_ban",
			Level:  1,
			IP:     c.ClientIP(),
			UserID: 0,
			OrgID:  0,
			Des:    fmt.Sprint("IP[", ipAddr, "]被禁用,但尝试访问URL:", c.Request.URL),
		})
		return
	}
	//继续
	c.Next()
}

// HeaderLoginBefore 登陆会话之前的请求，包含未登陆会话
func HeaderLoginBefore(c *gin.Context) {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("router mid, ", r)
		}
	}()
	//总的中间件
	RouterMid(c)
	//尝试获取token
	err := RouterMidAPI.TokenGetCookie(c)
	if err != nil {
		return
	}
	//对token进行安全事件检查
	tokenInfo := getTokenInfo(c)
	SafetyTokenON, err := BaseConfig.GetDataBool("SafetyTokenON")
	if err != nil {
		SafetyTokenON = true
	}
	if SafetyTokenON && BasePedometer.CheckData(&CoreRPCX.ArgsFrom{From: CoreSQLFrom.FieldsFrom{System: "safe_token", ID: tokenInfo.ID}}) {
		BaseSafe.CreateLog(&BaseSafe.ArgsCreateLog{
			System: "api.token_ban",
			Level:  1,
			IP:     c.ClientIP(),
			UserID: 0,
			OrgID:  0,
			Des:    fmt.Sprint("会话[", tokenInfo.ID, "]被禁用,会话组织ID:", tokenInfo.OrgID, ",用户ID:", tokenInfo.UserID, ",设备ID:", tokenInfo.DeviceID, ",但尝试访问URL:", c.Request.URL),
		})
		c.Redirect(200, "/501")
		return
	}
	//尝试获取用户数据
	if b, _ := getUserData(c); !b {
		return
	}
	//对token进行续约
	if err := RouterMidAPI.UpdateTokenCookie(c); err != nil {
		RouterReport.WarnLog(c, "token is expire, ", err, "token_error", "无法更新会话时间")
		return
	}
	//继续
	c.Next()
}

// HeaderLoggedUser 登陆会话之后的请求，包含已登陆会话
func HeaderLoggedUser(c *gin.Context) {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("router mid, ", r)
		}
	}()
	//总的中间件
	RouterMid(c)
	//尝试获取token
	err := RouterMidAPI.TokenGetCookie(c)
	if err != nil {
		c.Redirect(200, "/login")
		return
	}
	//对token进行安全事件检查
	tokenInfo := getTokenInfo(c)
	SafetyTokenON, err := BaseConfig.GetDataBool("SafetyTokenON")
	if err != nil {
		SafetyTokenON = true
	}
	if SafetyTokenON && BasePedometer.CheckData(&CoreRPCX.ArgsFrom{From: CoreSQLFrom.FieldsFrom{System: "safe_token", ID: tokenInfo.ID}}) {
		BaseSafe.CreateLog(&BaseSafe.ArgsCreateLog{
			System: "api.token_ban",
			Level:  1,
			IP:     c.ClientIP(),
			UserID: 0,
			OrgID:  0,
			Des:    fmt.Sprint("会话[", tokenInfo.ID, "]被禁用,会话组织ID:", tokenInfo.OrgID, ",用户ID:", tokenInfo.UserID, ",设备ID:", tokenInfo.DeviceID, ",但尝试访问URL:", c.Request.URL),
		})
		c.Redirect(200, "/501")
		return
	}
	//尝试获取用户数据
	if b, err := getUserData(c); err != nil {
		if b {
			c.Redirect(200, "/login")
			return
		}
	}
	//对token进行续约
	if err := RouterMidAPI.UpdateTokenCookie(c); err != nil {
		RouterReport.WarnLog(c, "token is expire, ", err, "token_error", "无法更新会话时间")
		c.Redirect(200, "/login")
		return
	}
	//继续
	c.Next()
}
