package BaseWeixinApp

import (
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	OrgCore "github.com/fotomxq/weeekj_core/v5/org/core"
)

// 获取商户的APPID等配置
func getAppConfig(orgID int64) (appID string, appKey string) {
	if orgID > 0 {
		androidKeys, _ := OrgCore.GetSystem(&OrgCore.ArgsGetSystem{
			OrgID:      orgID,
			SystemMark: "android",
		})
		iosKeys, _ := OrgCore.GetSystem(&OrgCore.ArgsGetSystem{
			OrgID:      orgID,
			SystemMark: "ios",
		})
		appID = androidKeys.Mark
		appKey = androidKeys.Params.GetValNoBool("key")
		if appID == "" {
			appID = iosKeys.Mark
			appKey = iosKeys.Params.GetValNoBool("key")
		}
	} else {
		appID, _ = BaseConfig.GetDataString("WeixinAppAppID")
		appKey, _ = BaseConfig.GetDataString("WeixinAppAppKey")
	}
	return
}
