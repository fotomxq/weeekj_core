package Router2APISecret

import (
	"errors"
	"fmt"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	"github.com/gin-gonic/gin"
)

//API验证模块

// GetSecretID 单独获取secret_id值
func GetSecretID(c *gin.Context) string {
	secretID := c.GetHeader("secret_id")
	return secretID
}

// CheckAPI 从header解析并验证密钥
// 必须给定secretID
// param c *gin.Context
// param urlAction string 验证的URL动作
// param key string 私钥
// return error
func CheckAPI(c *gin.Context, urlAction string, key string) error {
	//获取参数
	action, timestamp, nonce, secretID, signatureKey, signatureMethod := getHeaderParams(c)
	//所有值必须存在，如果secretID为空
	if action == "" || action != urlAction || timestamp == "" || nonce == "" || secretID == "" || signatureKey == "" || signatureMethod == "" {
		//输出错误
		return errors.New("api params error, cannot post empty")
	}
	//确保key值
	return checkSignatureKey(action, timestamp, nonce, secretID, signatureKey, signatureMethod, key)
}

// checkSignatureKey 验证一组数据是否合法？
// param action string URL接口名称
// param timestamp string 时间戳
// param nonce string 随机数
// param secretID string 配对id
// param signatureKey string 配对密钥
// param signatureMethod string 密钥算法
// param key string 私钥
// return error
func checkSignatureKey(action string, timestamp string, nonce string, secretID string, signatureKey string, signatureMethod string, key string) error {
	/**
	TODO：未来配送系统和养老系统需迭代至少3-5个版本，再加回来
	timestampInt64, _ := CoreFilter.GetInt64ByString(timestamp)
	if timestampInt64 < 1 || timestampInt64 > CoreFilter.GetNowTime().Unix()+300 || timestampInt64 < CoreFilter.GetNowTime().Unix()-300 {
		return errors.New("key timestamp error")
	}
	*/
	str, err := makeSignatureKey(action, timestamp, nonce, secretID, key, signatureMethod)
	if err != nil {
		return errors.New("system error, cannot get key, " + err.Error())
	}
	if str == signatureKey {
		return nil
	}
	return errors.New(fmt.Sprint("signature key error, str: "+str+", signatureKey: "+signatureKey, ", action: ", action, ", times: ", timestamp, ", nonce: ", nonce, ", secretID: ", secretID, ", signatureMethod: ", signatureMethod))
	//return errors.New("signature key error")
}

// makeSignatureKey 生成signature_key
// param action string URL接口名称
// param timestamp string 时间戳
// param nonce string 随机数
// param secretID string 配对id
// param key string 私钥
// param signatureMethod string 密钥算法
// return string signature_key
// return error
func makeSignatureKey(action string, timestamp string, nonce string, secretID string, key string, signatureMethod string) (string, error) {
	str := action + timestamp + nonce + secretID + key
	if signatureMethod == "sha256" {
		strSha256, err := CoreFilter.GetSha256([]byte(str))
		if err != nil {
			return "", err
		}
		return string(strSha256), nil
	}
	return "", errors.New("cannot find signature_method")
}

// getHeaderParams 从header获取数据
func getHeaderParams(c *gin.Context) (string, string, string, string, string, string) {
	//从form获取数据
	action := c.GetHeader("action")
	timestamp := c.GetHeader("timestamp")
	nonce := c.GetHeader("nonce")
	secretID := c.GetHeader("secret_id")
	signatureKey := c.GetHeader("signature_key")
	signatureMethod := c.GetHeader("signature_method")
	//反馈数据
	return action, timestamp, nonce, secretID, signatureKey, signatureMethod
}
