package Router2VCode

import (
	BaseVcode "github.com/fotomxq/weeekj_core/v5/base/vcode"
	Router2Mid "github.com/fotomxq/weeekj_core/v5/router2/mid"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// CheckImage 检查图形验证码
func CheckImage(c any, value string) bool {
	ctx, _ := getContextData(c)
	//如果启动debug，则自动忽略
	if Router2SystemConfig.Debug {
		return true
	}
	//获取token
	tokenInfo := Router2Mid.GetTokenInfo(ctx)
	//验证验证码是否正确
	if !BaseVcode.Check(&BaseVcode.ArgsCheck{
		Token: tokenInfo.ID,
		Value: value,
	}) {
		Router2Mid.ReportWarnLog(c, "check img vcode", nil, "err_img_vcode")
		return false
	}
	return true
}
