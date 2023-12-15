package UserUserMQTT

import (
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
)

// CheckType 检查数据聚合
type CheckType struct {
	//用户ID
	UserID int64 `json:"userID"`
	//token
	Token int64 `json:"token"`
	//随机字符
	Nonce string `json:"nonce"`
	//计算密钥结果
	// sha1(MQTT主题地址+用户ID+token+nonce)
	Key string `json:"key"`
}

// CheckKey 检查用户合法性
func CheckKey(topic string, keys CheckType) (b bool) {
	sha1Str := fmt.Sprint(topic, keys.UserID, keys.Token, keys.Nonce)
	sha1 := CoreFilter.GetSha1Str(sha1Str)
	if sha1 == "" {
		return
	}
	b = sha1 == keys.Key
	if !b {
		CoreLog.Warn("user mqtt, check permission failed, user id: ", keys.UserID, ", topic: ", topic)
	}
	return
}
