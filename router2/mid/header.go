package Router2Mid

import (
	"errors"
	"fmt"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	BasePedometer "github.com/fotomxq/weeekj_core/v5/base/pedometer"
	BaseSafe "github.com/fotomxq/weeekj_core/v5/base/safe"
	BaseToken2 "github.com/fotomxq/weeekj_core/v5/base/token2"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreRPCX "github.com/fotomxq/weeekj_core/v5/core/rpcx"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/gin-gonic/gin"
)

// headerBaseData 设定无头信息的请求
func headerBaseData(c *gin.Context) {
	//对ip进行安全检查
	SafetyIPON, err := BaseConfig.GetDataBool("SafetyIPON")
	if err != nil {
		SafetyIPON = true
	}
	ipAddr := c.ClientIP()
	if !CoreFilter.CheckIP(ipAddr) {
		reportGin(c, true, 0, nil, "ip:"+ipAddr, false, "err_ip", 0, nil)
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
		reportGin(c, false, 0, nil, "", false, "err_ip_ban", 0, nil)
		return
	}
	//继续
	c.Next()
}

// headerLoginBefore 登陆会话之前的请求，包含未登陆会话
func headerLoginBefore(c *gin.Context) {
	//采用api基础验证
	if err := checkAPIAndToken(c, c.Request.URL.Path); err != nil {
		reportGin(c, Router2SystemConfig.GlobConfig.Router.NeedTokenLog, 0, err, "token not api type", false, "token_error", 0, nil)
		return
	}
	//对token进行安全事件检查
	tokenInfo := GetTokenInfo(c)
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
		reportGin(c, false, 0, nil, "", false, "err_token_ban", 0, nil)
		return
	}
	//继续
	c.Next()
}

// checkAPIAndToken 获取token并验证
func checkAPIAndToken(c *gin.Context, urlAction string) error {
	//从form获取数据
	action := c.GetHeader("action")
	timestamp := c.GetHeader("timestamp")
	nonce := c.GetHeader("nonce")
	secretID := c.GetHeader("secret_id")
	signatureKey := c.GetHeader("signature_key")
	signatureMethod := c.GetHeader("signature_method")
	//检查api
	tokenID, err := BaseToken2.Check(action, urlAction, timestamp, nonce, secretID, signatureKey, signatureMethod)
	if err != nil || tokenID < 1 {
		return errors.New(fmt.Sprint("check token api, ", err))
	}
	//保存token数据
	c.Set("tokenID", tokenID)
	//反馈
	return nil
}
