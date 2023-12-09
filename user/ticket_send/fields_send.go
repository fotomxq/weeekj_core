package UserTicketSend

import "time"

//FieldsSend 群发优惠券约定
type FieldsSend struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//完成时间
	FinishAt time.Time `db:"finish_at" json:"finishAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//完成发放用户个数
	SendCount int64 `db:"send_count" json:"sendCount"`
	//是否必须是会员配置ID
	NeedUserSubConfigID int64 `db:"need_user_sub_config_id" json:"needUserSubConfigID"`
	//是否自动发放，如果不是，则需绑定广告
	NeedAuto bool `db:"need_auto" json:"needAuto"`
	//发放的票据配置ID
	ConfigID int64 `db:"config_id" json:"configID"`
	//每个用户发放几张
	PerCount int64 `db:"per_count" json:"perCount"`
}
