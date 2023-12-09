package RouterOrgAuth

import (
	"fmt"
	ClassConfig "gitee.com/weeekj/weeekj_core/v5/class/config"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	OrgCore "gitee.com/weeekj/weeekj_core/v5/org/core"
	RouterMidOrg "gitee.com/weeekj/weeekj_core/v5/router/mid/org"
	RouterReport "gitee.com/weeekj/weeekj_core/v5/router/report"
	"github.com/gin-gonic/gin"
	"strings"
)

// CheckAuth 处理商户授权
// orgID 授权的商户
// configMark 授权配置，位于授权商户名下的配置
func CheckAuth(c *gin.Context, orgID int64, configMark string) (b bool) {
	//获取组织
	orgData := RouterMidOrg.GetOrg(c)
	//获取参数商户
	authConfig, err := OrgCore.Config.GetConfigVal(&ClassConfig.ArgsGetConfig{
		BindID:    orgID,
		Mark:      configMark,
		VisitType: "admin",
	})
	if err != nil {
		RouterReport.WarnLog(c, fmt.Sprint("org manager get org auth data failed, no auth config, org id: ", orgID, ", need auth org: ", orgData.ID), nil, "no_auth", "无授权")
		return
	}
	//检查授权
	if !checkAuthConfig(orgData.ID, authConfig) {
		//检查全局授权
		authConfig, err = OrgCore.Config.GetConfigVal(&ClassConfig.ArgsGetConfig{
			BindID:    orgID,
			Mark:      "AuthorizationAll",
			VisitType: "admin",
		})
		if err != nil {
			RouterReport.WarnLog(c, fmt.Sprint("org manager get org auth data failed, no auth config, org id: ", orgID, ", need auth org: ", orgData.ID), nil, "no_auth", "无授权")
			return
		}
		if !checkAuthConfig(orgData.ID, authConfig) {
			RouterReport.WarnLog(c, fmt.Sprint("org manager get org auth data failed, no auth, org id: ", orgID, ", need auth org: ", orgData.ID), nil, "no_auth", "无授权")
			return
		}
	}
	return true
}

func checkAuthConfig(needOrgID int64, authConfig string) bool {
	if authConfig == "" {
		return false
	}
	authConfigs := strings.Split(authConfig, ",")
	for _, v := range authConfigs {
		vInt64, err := CoreFilter.GetInt64ByString(v)
		if err != nil {
			continue
		}
		if needOrgID == vInt64 {
			return true
		}
	}
	return false
}
