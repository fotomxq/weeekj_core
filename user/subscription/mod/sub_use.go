package UserSubscriptionMod

import (
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
)

// ArgsUseSub 使用目标订阅参数
type ArgsUseSub struct {
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//使用来源
	UseFrom     string `db:"use_from" json:"useFrom"`
	UseFromName string `db:"use_from_name" json:"useFromName"`
}

// UseSub 使用目标订阅
func UseSub(args ArgsUseSub) (err error) {
	//通知nats
	CoreNats.PushDataNoErr("/user/sub/use", "", 0, "", args)
	//反馈
	return
}
