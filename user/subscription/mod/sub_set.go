package UserSubscriptionMod

import (
	CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"
	"time"
)

// ArgsSetSub 设置订阅信息参数
type ArgsSetSub struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//新的到期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt" check:"isoTime"`
	//是否为继续订阅
	// 否则将覆盖过期时间
	HaveExpire bool `db:"have_expire" json:"haveExpire" check:"bool"`
	//使用来源
	UseFrom     string `db:"use_from" json:"useFrom"`
	UseFromName string `db:"use_from_name" json:"useFromName"`
}

// SetSub 设置订阅信息
func SetSub(args ArgsSetSub) (err error) {
	//通知修改
	CoreNats.PushDataNoErr("/user/sub/set", "", 0, "", args)
	//反馈
	return
}

// ArgsSetSubAdd 向后续约指定时间参数
type ArgsSetSubAdd struct {
	ConfigID int64 `json:"configID"`
	UserID   int64 `json:"userID"`
	Unit     int   `json:"unit"`
	OrderID  int64 `json:"orderID"`
}

// SetSubAdd 向后续约指定时间
func SetSubAdd(args *ArgsSetSubAdd) (err error) {
	//通知修改
	CoreNats.PushDataNoErr("/user/sub/set_add", "", 0, "", args)
	//反馈
	return
}
