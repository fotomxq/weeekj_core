package WeixinPayV3

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	CoreFile "github.com/fotomxq/weeekj_core/v5/core/file"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	FinancePayMod "github.com/fotomxq/weeekj_core/v5/finance/pay/mod"
	OrgCoreCore "github.com/fotomxq/weeekj_core/v5/org/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

var (
	//全局根URL地址
	baseURL = "https://api.mch.weixin.qq.com"
)

// 组装头部
func getClient(orgID int64) (*core.Client, dataClientConfig, error) {
	//隶属关系
	orgID = FinancePayMod.FixOrgID(orgID)
	//获取配置项
	merchantConfig, err := getClientConfig(orgID)
	if err != nil {
		err = errors.New(fmt.Sprint("get client error, org id: ", orgID, ", err: ", err))
		return &core.Client{}, dataClientConfig{}, err
	}
	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	var mchPrivateKey *rsa.PrivateKey
	/** 禁止从IO本地加载PEM数据
	if merchantConfig.KeyPEM == "" {
		mchPrivateKey, err = utils.LoadPrivateKeyWithPath(fmt.Sprint(merchantConfig.CertDir, CoreFile.Sep, "apiclient_key.pem"))
		if err != nil {
			err = errors.New(fmt.Sprint("load merchant private key error", ", cert src: ", merchantConfig.CertDir, CoreFile.Sep, "apiclient_key.pem"))
			return &core.Client{}, dataClientConfig{}, err
		}
	} else {
		mchPrivateKey, err = utils.LoadPrivateKey(merchantConfig.KeyPEM)
		if err != nil {
			err = errors.New(fmt.Sprint("load merchant private key by config error"))
			return &core.Client{}, dataClientConfig{}, err
		}
	}
	*/
	if merchantConfig.KeyPEM == "" {
		mchPrivateKey, err = utils.LoadPrivateKeyWithPath(fmt.Sprint(merchantConfig.CertDir, CoreFile.Sep, "apiclient_key.pem"))
		if err != nil {
			err = errors.New(fmt.Sprint("load merchant private key error", ", cert src: ", merchantConfig.CertDir, CoreFile.Sep, "apiclient_key.pem"))
			return &core.Client{}, dataClientConfig{}, err
		}
	} else {
		mchPrivateKey, err = utils.LoadPrivateKey(merchantConfig.KeyPEM)
		if err != nil {
			err = errors.New(fmt.Sprint("load merchant private key pem by config error, org id: ", orgID, ", err: ", err))
			return &core.Client{}, dataClientConfig{}, err
		}
	}
	ctx := context.Background()
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(merchantConfig.MerchantID, merchantConfig.CertSN, mchPrivateKey, merchantConfig.KeyV3),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		err = errors.New(fmt.Sprint("new wechat pay client err: ", err))
		CoreLog.Warn("sn: ", merchantConfig.CertSN, ", pem: ", merchantConfig.KeyPEM)
		return &core.Client{}, dataClientConfig{}, err
	}
	return client, merchantConfig, nil
}

// 通用参数结构体
type dataClientConfig struct {
	//商户ID
	MerchantID string
	//子商户ID
	SubMerchantID string
	//证书序列号
	CertSN string
	//商户APIv3 Key
	KeyV3 string
	//证书内容
	KeyPEM string
	//证书路径
	// 将自动识别该目录下的三个证书文件
	// 请勿修改默认名称，分别为：apiclient_cert.p12/apiclient_cert.pem/apiclient_key.pem
	CertDir string
	//访问得URL
	BaseURL string
	//其他扩展项
	Params CoreSQLConfig.FieldsConfigsType
}

func getClientConfig(orgID int64) (data dataClientConfig, err error) {
	//全局通用设置
	data.CertDir = getCertDirSrc()
	data.BaseURL = baseURL
	var financePayOtherInOne bool
	financePayOtherInOne, err = BaseConfig.GetDataBool("FinancePayOtherInOne")
	if err != nil {
		err = nil
		financePayOtherInOne = false
	}
	//如果具有组织ID，优先从组织从抽取数据
	if orgID > 0 && !financePayOtherInOne {
		//获取商户是否开通特定服务
		openFinanceIndependent := OrgCoreCore.CheckOrgPermissionFunc(orgID, "finance_independent")
		if openFinanceIndependent {
			//否则获取指定的商户数据包
			var orgConfig OrgCoreCore.FieldsSystem
			orgConfig, err = OrgCoreCore.GetSystem(&OrgCoreCore.ArgsGetSystem{
				OrgID:      orgID,
				SystemMark: "weixin_pay",
			})
			if err == nil {
				var b bool
				data.MerchantID = orgConfig.Mark
				data.SubMerchantID, _ = orgConfig.Params.GetVal("subMerchantID")
				data.CertSN, b = orgConfig.Params.GetVal("certSN")
				if !b {
					err = errors.New("no pay cert sn")
					return
				}
				data.KeyV3, b = orgConfig.Params.GetVal("keyV3")
				if !b {
					err = errors.New("no pay key v3")
					return
				}
				data.KeyPEM, b = orgConfig.Params.GetVal("keyPEM")
				if !b {
					err = errors.New("no pay key pem")
					return
				}
				data.CertDir = fmt.Sprint(data.CertDir, CoreFile.Sep, "org", CoreFile.Sep, orgID, CoreFile.Sep, "weixin_pay")
				data.Params = orgConfig.Params
			} else {
				err = errors.New(fmt.Sprint("get org system config, org id: ", orgID, ", system mark: weixin_pay, err: ", err))
			}
			return
		}
	}
	//重新加载数据
	data.MerchantID, err = BaseConfig.GetDataString("WeixinXiaochengxuPayMerchantID")
	if err != nil {
		err = errors.New("get config by WeixinXiaochengxuPayMerchantID, " + err.Error())
		return
	}
	data.SubMerchantID, err = BaseConfig.GetDataString("WeixinXiaochengxuPaySubMerchantID")
	if err != nil {
		err = errors.New("get config by WeixinXiaochengxuPaySubMerchantID, " + err.Error())
		return
	}
	data.CertSN, err = BaseConfig.GetDataString("WeixinXiaochengxuPayCertSN")
	if err != nil {
		err = errors.New("get config by WeixinXiaochengxuPayCertSN, " + err.Error())
		return
	}
	data.KeyV3, err = BaseConfig.GetDataString("WeixinXiaochengxuPayMerchantKeyV3")
	if err != nil {
		err = errors.New("get config by WeixinXiaochengxuPayMerchantKeyV3, " + err.Error())
		return
	}
	data.KeyPEM, err = BaseConfig.GetDataString("WeixinXiaochengxuPayMerchantKeyPEM")
	if err != nil {
		err = errors.New("get config by WeixinXiaochengxuPayMerchantKeyPEM, " + err.Error())
		return
	}
	data.CertDir = fmt.Sprint(data.CertDir, CoreFile.Sep, "glob", CoreFile.Sep, "weixin_pay")
	//反馈数据
	return
}

// 设置证书位置
func getCertDirSrc() string {
	return fmt.Sprint(CoreFile.BaseDir(), CoreFile.Sep, "conf", CoreFile.Sep, "cert")
}
