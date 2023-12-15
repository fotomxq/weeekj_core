package BaseWeixinPayClient

import (
	"errors"
	"fmt"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	CoreFile "github.com/fotomxq/weeekj_core/v5/core/file"
	OrgCoreCore "github.com/fotomxq/weeekj_core/v5/org/core"
)

// ClientType 商户支付平台初始化组件
type ClientType struct {
	//配置内容
	ConfigData configType
	BaseURL    string
}

type configType struct {
	//小程序或服务商
	AppID string
	//商户ID
	MerchantID string
	//商户API Key
	Key string
	//商户APIv3 Key
	KeyV3 string
	//PEM证书
	KeyPEM string
	//证书路径
	// 将自动识别该目录下的三个证书文件
	// 请勿修改默认名称，分别为：apiclient_cert.p12/apiclient_cert.pem/apiclient_key.pem
	CertDir string
}

func getAppID() (string, string, error) {
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

// 获取商户接口参数
func getPayMerchantBase() (string, string, error) {
	mchID, err := BaseConfig.GetDataString("WeixinXiaochengxuPayMerchantID")
	if err != nil {
		return "", "", errors.New("get config by WeixinXiaochengxuPayMerchantID, " + err.Error())
	}
	mchKey, err := BaseConfig.GetDataString("WeixinXiaochengxuPayMerchantKey")
	if err != nil {
		return "", "", errors.New("get config by WeixinXiaochengxuPayMerchantKey, " + err.Error())
	}
	return mchID, mchKey, nil
}

// 获取商户v3 key
func getPayMerchantV3Key() (string, error) {
	key, err := BaseConfig.GetDataString("WeixinXiaochengxuPayMerchantKeyV3")
	if err != nil {
		return "", errors.New("get config by WeixinXiaochengxuPayMerchantID, " + err.Error())
	}
	return key, nil
}

// 设置证书位置
func getCertDirSrc() string {
	return fmt.Sprint(CoreFile.BaseDir(), CoreFile.Sep, "conf", CoreFile.Sep, "cert")
}

// GetMerchantClient 获取指定ID的商户client数据
// fromSystem 来源于不同的渠道，例如小程序，将影响支付过程中的appID组成部分
// 支持: 0 微信小程序
func GetMerchantClient(fromSystem int, orgID int64) (client ClientType, err error) {
	//全局通用设置
	client.ConfigData.CertDir = getCertDirSrc()
	client.BaseURL = "https://api.mch.weixin.qq.com"
	//如果具有组织ID，优先从组织从抽取数据
	if orgID > 0 {
		//否则获取指定的商户数据包
		var orgConfig OrgCoreCore.FieldsSystem
		orgConfig, err = OrgCoreCore.GetSystem(&OrgCoreCore.ArgsGetSystem{
			OrgID:      orgID,
			SystemMark: "weixin_pay",
		})
		if err != nil {
			err = errors.New(fmt.Sprint("get org: ", orgID, ", system config by weixin_pay, ", err))
			return
		}
		var orgWxxConfig OrgCoreCore.FieldsSystem
		orgWxxConfig, err = OrgCoreCore.GetSystem(&OrgCoreCore.ArgsGetSystem{
			OrgID:      orgID,
			SystemMark: "weixin_wxx",
		})
		if err != nil {
			err = errors.New(fmt.Sprint("get org: ", orgID, ", system config by weixin_wxx, ", err))
			return
		}
		var b bool
		client.ConfigData.AppID = orgWxxConfig.Mark
		client.ConfigData.MerchantID = orgConfig.Mark
		client.ConfigData.Key, b = orgConfig.Params.GetVal("key")
		if !b {
			err = errors.New("no pay key")
			return
		}
		client.ConfigData.KeyV3, b = orgConfig.Params.GetVal("keyV3")
		if !b {
			err = errors.New("no pay key v3")
			return
		}
		client.ConfigData.KeyPEM, b = orgConfig.Params.GetVal("keyPEM")
		if !b {
			err = errors.New("no pay key pem")
			return
		}
		client.ConfigData.CertDir = fmt.Sprint(client.ConfigData.CertDir, CoreFile.Sep, "org", CoreFile.Sep, orgID, CoreFile.Sep, "weixin_pay")
		return
	} else {
		//重新加载数据
		client.ConfigData.AppID, _, err = getAppID()
		if err != nil {
			return
		}
		client.ConfigData.MerchantID, client.ConfigData.Key, err = getPayMerchantBase()
		if err != nil {
			return
		}
		client.ConfigData.KeyV3, err = getPayMerchantV3Key()
		if err != nil {
			return
		}
		client.ConfigData.KeyPEM, err = BaseConfig.GetDataString("WeixinXiaochengxuPayMerchantKeyPEM")
		if err != nil {
			err = errors.New("get config by WeixinXiaochengxuPayMerchantKeyPEM, " + err.Error())
			return
		}
		client.ConfigData.CertDir = fmt.Sprint(client.ConfigData.CertDir, CoreFile.Sep, "glob", CoreFile.Sep, "weixin_pay")
		//反馈数据
		return
	}
}
