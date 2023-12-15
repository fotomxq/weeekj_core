package TMSSelfOtherKuai100

import (
	"errors"
	"fmt"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	OrgCore "github.com/fotomxq/weeekj_core/v5/org/core"
	"strings"
)

// getCustomer 获取授权码
func getCustomer(orgID int64) (customer string) {
	//获取商户配置
	customer = OrgCore.Config.GetConfigValNoErr(orgID, "TMSSelfOtherKuai100Customer")
	if customer != "" {
		return
	}
	//获取平台配置
	customer = BaseConfig.GetDataStringNoErr("TMSSelfOtherKuai100Customer")
	//反馈
	return
}

// getKey 获取快递100key
func getKey(orgID int64) (key string) {
	//获取商户配置
	key = OrgCore.Config.GetConfigValNoErr(orgID, "TMSSelfOtherKuai100Key")
	if key != "" {
		return
	}
	//获取平台配置
	key = BaseConfig.GetDataStringNoErr("TMSSelfOtherKuai100Key")
	//反馈
	return
}

// getSign 计算密钥
func getSign(orgID int64, param string) (sign string, errCode string, err error) {
	//获取授权码和密钥
	customer := getCustomer(orgID)
	if customer == "" {
		errCode = "err_config"
		err = errors.New("customer is empty")
		return
	}
	key := getKey(orgID)
	if key == "" {
		errCode = "err_config"
		err = errors.New("key is empty")
		return
	}
	//计算
	md5Str := fmt.Sprint(param, key, customer)
	sign = getSignMD5(md5Str)
	//反馈
	return
}

// 计算密钥的算法底层支持
func getSignMD5(str string) (sign string) {
	//计算
	sign = CoreFilter.GetMd5StrByStr(str)
	//转为大写
	sign = strings.ToUpper(sign)
	//反馈
	return
}
