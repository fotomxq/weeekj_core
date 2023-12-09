package RouterVcode

import (
	BaseSMS "gitee.com/weeekj/weeekj_core/v5/base/sms"
	BaseVcode "gitee.com/weeekj/weeekj_core/v5/base/vcode"
	RouterReport "gitee.com/weeekj/weeekj_core/v5/router/report"
	Router2Mid "gitee.com/weeekj/weeekj_core/v5/router2/mid"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/gin-gonic/gin"
)

//验证码通用规则
// 用于根据验证码和token，自动检查数据并验证参数正确性
// 可用于一般验证码和短信验证码

// CheckImage 检查图形验证码
func CheckImage(c *gin.Context, value string) bool {
	//如果启动debug，则自动忽略
	if Router2SystemConfig.Debug {
		return true
	}
	//获取token
	tokenInfo := Router2Mid.GetTokenInfo(c)
	//验证验证码是否正确
	if !BaseVcode.Check(&BaseVcode.ArgsCheck{
		Token: tokenInfo.ID, Value: value,
	}) {
		RouterReport.BaseError(c, "vcode_error", "验证码错误")
		return false
	}
	return true
}

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
func CheckSMS(c *gin.Context, value string, phoneData ParamsPhoneType) (BaseSMS.FieldsSMS, bool) {
	//获取token
	tokenInfo := Router2Mid.GetTokenInfo(c)
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
		RouterReport.BaseError(c, "err_sms_vcode", "验证码错误")
		return BaseSMS.FieldsSMS{}, false
	} else {
		if phoneData.Allow {
			if phoneData.NationCode == data.NationCode && phoneData.Phone == data.Phone {
				return data, true
			}
			RouterReport.BaseError(c, "err_sms_phone", "电话错误")
			return BaseSMS.FieldsSMS{}, false
		}
		return data, true
	}
}

// CheckSMS2 检查短信验证码
func CheckSMS2(c any, value string, phoneData ParamsPhoneType) (BaseSMS.FieldsSMS, bool) {
	//获取上下文
	ctx := Router2Mid.GetContext(c)
	//获取token
	tokenID := Router2Mid.GetTokenID(ctx)
	//检查短信验证码
	if data, b := BaseSMS.CheckSMSAndData(&BaseSMS.ArgsCheckSMSAndData{
		OrgID:    phoneData.OrgID,
		ConfigID: -1,
		Token:    tokenID,
		Value:    value,
	}); !b {
		//如果启动debug模式，则无论是false，也会读取data并反馈
		if Router2SystemConfig.Debug {
			return data, true
		}
		Router2Mid.ReportBaseError(c, "err_sms_vcode")
		return BaseSMS.FieldsSMS{}, false
	} else {
		if phoneData.Allow {
			if phoneData.NationCode == data.NationCode && phoneData.Phone == data.Phone {
				return data, true
			}
			Router2Mid.ReportBaseError(c, "err_sms_phone")
			return BaseSMS.FieldsSMS{}, false
		}
		return data, true
	}
}
