package BaseToken2

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// Check 检查token
func Check(action, urlAction, timestamp, nonce, secretID, signatureKey, signatureMethod string) (tokenID int64, err error) {
	//获取token ID
	tokenID, _ = CoreFilter.GetInt64ByString(secretID)
	if tokenID < 1 {
		err = errors.New("secretID not int64")
		return
	}
	//获取token
	data := getByID(tokenID)
	if data.ID < 1 {
		err = errors.New(fmt.Sprint("get token not data by id: ", tokenID))
		return
	}
	if data.ExpireAt.Unix() < CoreFilter.GetNowTime().Unix() {
		err = errors.New("token is expire")
		return
	}
	//检查密钥
	err = CheckByKey(action, urlAction, timestamp, nonce, secretID, signatureKey, signatureMethod, data.Key)
	if err != nil {
		err = errors.New(fmt.Sprint("token check key, ", err))
		return
	}
	//更新token过期时间
	updateExpire(&data)
	//反馈
	return
}

// CheckByKey 计算密钥数据包
// 原先位于RouterMid中间件声明
// action/timestamp/nonce/secretID/signatureKey/signatureMethod 为header头部提交的数据
// urlAction 实际接口地址
// key 密钥数据
func CheckByKey(action, urlAction, timestamp, nonce, secretID, signatureKey, signatureMethod, key string) (err error) {
	if action == "" || action != urlAction || timestamp == "" || nonce == "" || secretID == "" || signatureKey == "" || signatureMethod == "" {
		err = errors.New("params is empty")
		return
	}
	timestampInt64, _ := CoreFilter.GetInt64ByString(timestamp)
	if timestampInt64 < 1 {
		err = errors.New("timestamp not int64")
		return
	}
	nowUnix := CoreFilter.GetNowTime().Unix()
	if Router2SystemConfig.GlobConfig.Safe.SafeRouterTimeBlocker {
		if timestampInt64 < nowUnix-300 || timestampInt64 > nowUnix+300 {
			err = errors.New("timestamp is error area")
			return
		}
	}
	str := action + timestamp + nonce + secretID + key
	var strKey string
	switch signatureMethod {
	case "sha256":
		var strSha256 []byte
		strSha256, err = CoreFilter.GetSha256([]byte(str))
		if err != nil {
			err = errors.New(fmt.Sprint("get sha256 error, ", err))
			return
		}
		strKey = string(strSha256)
	case "sha1":
		var strSHA1 []byte
		strSHA1, err = CoreFilter.GetSha1([]byte(str))
		if err != nil {
			err = errors.New(fmt.Sprint("get sha1 error, ", err))
			return
		}
		strKey = string(strSHA1)
	default:
		return
	}
	if strKey != signatureKey {
		//err = errors.New(fmt.Sprint("signatureKey is error, strKEY: ", strKey, ", signatureKey: ", signatureKey, ", raw str: ", str))
		err = errors.New("signatureKey is error")
		return
	}
	return
}
