package Router2Mid

import (
	"fmt"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	BasePedometer "github.com/fotomxq/weeekj_core/v5/base/pedometer"
	BaseSafe "github.com/fotomxq/weeekj_core/v5/base/safe"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/gin-gonic/gin"
)

// headerLoggedUser 登陆会话之后的请求，包含已登陆会话
func headerLoggedUser(c *gin.Context) {
	//采用api基础验证
	if err := checkAPIAndToken(c, c.Request.URL.Path); err != nil {
		reportGin(c, Router2SystemConfig.GlobConfig.Router.NeedTokenLog, 0, err, "token not api type", false, "token_error", 0, nil)
		return
	}
	//对token进行安全事件检查
	tokenInfo := GetTokenInfo(c)
	//是否启动安全预警
	SafetyTokenON, err := BaseConfig.GetDataBool("SafetyTokenON")
	if err != nil {
		SafetyTokenON = true
	}
	if SafetyTokenON && BasePedometer.CheckData(CoreSQLFrom.FieldsFrom{System: "safe-token", ID: tokenInfo.ID}) {
		BaseSafe.CreateLog(&BaseSafe.ArgsCreateLog{
			System: "api.token_ban",
			Level:  1,
			IP:     c.ClientIP(),
			UserID: 0,
			OrgID:  0,
			Des:    fmt.Sprint("会话[", tokenInfo.ID, "]被禁用,会话组织ID:", tokenInfo.OrgID, ",用户ID:", tokenInfo.UserID, ",设备ID:", tokenInfo.DeviceID, ",但尝试访问URL:", c.Request.URL),
		})
		reportGin(c, false, 0, err, "", false, "err_token_ban", 0, nil)
		return
	}
	//获取用户结构
	userID, err := getUserDataByToken(c)
	if err != nil {
		reportGin(c, false, 0, err, "", false, "err_user", 0, nil)
		return
	}
	//对user进行安全事件检查
	SafetyUserON, err := BaseConfig.GetDataBool("SafetyUserON")
	if err != nil {
		SafetyUserON = true
	}
	if SafetyUserON && BasePedometer.CheckData(CoreSQLFrom.FieldsFrom{System: "safe_user", ID: userID}) {
		BaseSafe.CreateLog(&BaseSafe.ArgsCreateLog{
			System: "api.token_ban",
			Level:  1,
			IP:     c.ClientIP(),
			UserID: 0,
			OrgID:  0,
			Des:    fmt.Sprint("用户[", userID, "]被禁用,但尝试访问URL:", c.Request.URL),
		})
		reportGin(c, false, 0, err, "", false, "err_user_ban", 0, nil)
		return
	}
	//继续
	c.Next()
}
