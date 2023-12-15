package BaseWeixinWXXClient

import (
	"errors"
	"fmt"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	OrgCoreCore "github.com/fotomxq/weeekj_core/v5/org/core"
	"time"
)

// ClientType 小程序初始化组件
type ClientType struct {
	//组织ID
	OrgID int64
	// BaseURL 微信请求基础URL
	BaseURL string
	//配置内容
	ConfigData ConfigType
	//token
	AccessToken string
	//token 过期时间
	AccessTokenExpireTime time.Time
	//获取过期时间间隔
	AccessTokenExpireDuration time.Duration
}

type ConfigType struct {
	//小程序或服务商
	AppID string
	//API Key
	Key string
}

// GetMerchantClient 获取指定ID的商户client数据
func GetMerchantClient(orgID int64) (client ClientType, err error) {
	//全局设置
	client.BaseURL = "https://api.weixin.qq.com"
	//获取支付体系控制开关
	var financePayOtherInOne bool
	financePayOtherInOne, err = BaseConfig.GetDataBool("FinancePayOtherInOne")
	if err != nil {
		err = nil
		financePayOtherInOne = false
	}
	//如果具有商户，则优先查询
	// 如果 不存在、删除，则按照没有商户处理
	if orgID > 0 && !financePayOtherInOne {
		//获取商户是否开通特定服务
		openFinanceIndependent := OrgCoreCore.CheckOrgPermissionFunc(orgID, "finance_independent")
		if openFinanceIndependent {
			//否则获取指定的商户数据包
			var orgConfig OrgCoreCore.FieldsSystem
			orgConfig, err = OrgCoreCore.GetSystem(&OrgCoreCore.ArgsGetSystem{
				OrgID:      orgID,
				SystemMark: "weixin_wxx",
			})
			if err == nil {
				var b bool
				client.ConfigData.AppID = orgConfig.Mark
				client.ConfigData.Key, b = orgConfig.Params.GetVal("key")
				if !b {
					err = errors.New("no key")
					return
				}
				client.OrgID = orgID
			} else {
				err = errors.New(fmt.Sprint("get org system config by old wxx pay, org id: ", orgID, ", system mark: weixin_wxx, err: ", err))
			}
			return
		}
	}
	client.ConfigData.AppID, client.ConfigData.Key, err = getGlobAppID()
	if err != nil {
		err = errors.New("get wxx app id or key failed, " + err.Error())
		return
	}
	client.OrgID = 0
	//反馈数据
	return
}

func getGlobAppID() (string, string, error) {
	appID, err := BaseConfig.GetDataString("WeixinXiaochengxuAppID")
	if err != nil {
		return "", "", errors.New("weixin xiaochengxu config is error, " + err.Error())
	}
	secret, err := BaseConfig.GetDataString("WeixinXiaochengxuSecret")
	if err != nil {
		return "", "", errors.New("weixin xiaochengxu config is error, " + err.Error())
	}
	return appID, secret, nil
}
