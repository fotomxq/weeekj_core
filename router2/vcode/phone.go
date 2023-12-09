package Router2VCode

import (
	BaseSMS "gitee.com/weeekj/weeekj_core/v5/base/sms"
	Router2Mid "gitee.com/weeekj/weeekj_core/v5/router2/mid"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ParamsPhoneType 需要验证手机号匹配的参数
type ParamsPhoneType struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//是否验证手机号一致性
	// 如果没有启动，则token可以用不同手机号发起验证，验证时不会二次对手机号核对，请慎重关闭
	// 可用于登陆时，关闭后可实现跳过验证手机号一致性，直接完成登陆操作
	Allow bool `json:"allow"`
	//手机区号和手机号
	NationCode string `json:"nationCode"`
	Phone      string `json:"phone"`
}

// CheckSMS 检查短信验证码
func CheckSMS(c any, value string, phoneData ParamsPhoneType) (BaseSMS.FieldsSMS, bool) {
	ctx, _ := getContextData(c)
	//获取token
	tokenInfo := Router2Mid.GetTokenInfo(ctx)
	//检查短信验证码
	if data, b := BaseSMS.CheckSMSAndData(&BaseSMS.ArgsCheckSMSAndData{
		OrgID:    phoneData.OrgID,
		ConfigID: -1,
		Token:    tokenInfo.ID,
		Value:    value,
	}); !b {
		//如果启动debug模式，则无论是false，也会读取data并反馈
		if Router2SystemConfig.Debug {
			return data, true
		}
		Router2Mid.ReportWarnLog(c, "sms phone vcode error", nil, "err_phone_vcode")
		return BaseSMS.FieldsSMS{}, false
	} else {
		if phoneData.Allow {
			if phoneData.NationCode == data.NationCode && phoneData.Phone == data.Phone {
				return data, true
			}
			Router2Mid.ReportWarnLog(c, "sms phone vcode error", nil, "err_phone")
			return BaseSMS.FieldsSMS{}, false
		}
		return data, true
	}
}
