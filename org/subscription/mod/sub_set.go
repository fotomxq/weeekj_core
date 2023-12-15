package OrgSubscriptionMod

import CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"

// ArgsSetSubAdd 向后续约指定时间参数
type ArgsSetSubAdd struct {
	ConfigID int64 `json:"configID"`
	OrgID    int64 `json:"orgID"`
	Unit     int   `json:"unit"`
	OrderID  int64 `json:"orderID"`
}

// SetSubAdd 向后续约指定时间
func SetSubAdd(args *ArgsSetSubAdd) (err error) {
	//通知修改
	CoreNats.PushDataNoErr("/org/sub/set_add", "", 0, "", args)
	//反馈
	return
}
