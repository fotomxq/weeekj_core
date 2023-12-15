package RouterMidAPI

import (
	"fmt"

	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	BasePedometer "github.com/fotomxq/weeekj_core/v5/base/pedometer"
	BaseSafe "github.com/fotomxq/weeekj_core/v5/base/safe"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreRPCX "github.com/fotomxq/weeekj_core/v5/core/rpcx"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	RouterMidCore "github.com/fotomxq/weeekj_core/v5/router/mid/core"
	RouterReport "github.com/fotomxq/weeekj_core/v5/router/report"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/gin-gonic/gin"
)

// HeaderBaseData 设定无头信息的请求
func HeaderBaseData(c *gin.Context) {
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
		RouterReport.BaseError(c, "ip_ban", "禁止访问")
		return
	}
	//继续
	c.Next()
}

// HeaderLoginBefore 登陆会话之前的请求，包含未登陆会话
func HeaderLoginBefore(c *gin.Context) {
	//设置origin
	RouterMidCore.SetOriginConfig(c)
	//采用api基础验证
	if err := checkAPIAndToken(c, c.Request.URL.Path); err != nil {
		if Router2SystemConfig.GlobConfig.Router.NeedTokenLog {
			RouterReport.WarnLog(c, "token not api type, ", err, "token_error", "API错误")
		} else {
			RouterReport.BaseError(c, "token_error", "API错误")
		}
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
		RouterReport.BaseError(c, "token_ban", "会话被禁用")
		return
	}
	//继续
	c.Next()
}

// HeaderLoggedUser 登陆会话之后的请求，包含已登陆会话
func HeaderLoggedUser(c *gin.Context) {
	//设置origin
	RouterMidCore.SetOriginConfig(c)
	//采用api基础验证
	if err := checkAPIAndToken(c, c.Request.URL.Path); err != nil {
		if Router2SystemConfig.GlobConfig.Router.NeedTokenLog {
			RouterReport.WarnLog(c, "token not api type, ", err, "token_error", "API错误")
		} else {
			RouterReport.BaseError(c, "token_error", "API错误")
		}
		return
	}
	//对token进行安全事件检查
	tokenInfo := getTokenInfo(c)
	//是否启动安全预警
	SafetyTokenON, err := BaseConfig.GetDataBool("SafetyTokenON")
	if err != nil {
		SafetyTokenON = true
	}
	if SafetyTokenON && BasePedometer.CheckData(&CoreRPCX.ArgsFrom{From: CoreSQLFrom.FieldsFrom{System: "safe-token", ID: tokenInfo.ID}}) {
		BaseSafe.CreateLog(&BaseSafe.ArgsCreateLog{
			System: "api.token_ban",
			Level:  1,
			IP:     c.ClientIP(),
			UserID: 0,
			OrgID:  0,
			Des:    fmt.Sprint("会话[", tokenInfo.ID, "]被禁用,会话组织ID:", tokenInfo.OrgID, ",用户ID:", tokenInfo.UserID, ",设备ID:", tokenInfo.DeviceID, ",但尝试访问URL:", c.Request.URL),
		})
		RouterReport.BaseError(c, "token_ban", "会话被禁用")
		return
	}
	//获取用户结构
	userData, err := GetUserDataByToken(c)
	if err != nil {
		//RouterReport.WarnLog(c, "token cannot get user data, ", err, "token_error", "无效用户")
		RouterReport.BaseError(c, "token_error", "无效用户")
		return
	}
	//对user进行安全事件检查
	SafetyUserON, err := BaseConfig.GetDataBool("SafetyUserON")
	if err != nil {
		SafetyUserON = true
	}
	if SafetyUserON && BasePedometer.CheckData(&CoreRPCX.ArgsFrom{From: CoreSQLFrom.FieldsFrom{System: "safe_user", ID: userData.Info.ID}}) {
		BaseSafe.CreateLog(&BaseSafe.ArgsCreateLog{
			System: "api.token_ban",
			Level:  1,
			IP:     c.ClientIP(),
			UserID: 0,
			OrgID:  0,
			Des:    fmt.Sprint("用户[", userData.Info.ID, "]被禁用,但尝试访问URL:", c.Request.URL),
		})
		RouterReport.BaseError(c, "user_ban", "用户被禁用")
		return
	}
	//继续
	c.Next()
}
