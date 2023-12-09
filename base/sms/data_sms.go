package BaseSMS

import "time"

type DataSMS struct {
	//创建时间
	CreateAt time.Time `json:"createAt"`
	//过期时间
	// 过期后不清理数据，但存在保留的最大时间
	ExpireAt time.Time `json:"expireAt"`
	//发送时间
	SendAt time.Time `json:"sendAt"`
	//组织ID
	OrgID int64 `json:"orgID"`
	//配置ID
	ConfigID int64 `json:"configID"`
	//会话
	Token int64 `json:"token"`
	//用户ID
	UserID int64 `json:"userID"`
	//国家代码
	NationCode string `json:"nationCode"`
	//目标手机号
	// 目标手机号是唯一的标识码
	Phone string `json:"phone"`
	//短信内容
	Des string `json:"des"`
	//失败原因
	// 如果为本地原因则显示错误代码，否则显示API提供方反馈信息
	FailedMsg string `json:"failedMsg"`
	//短信类型
	// check 验证类; des 内容类
	UseType string `json:"useType"`
	//是否已经验证
	HaveCheck time.Time `json:"haveCheck"`
}
