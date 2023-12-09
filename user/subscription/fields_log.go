package UserSubscription

import "time"

//FieldsLog 使用日志
type FieldsLog struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//使用来源
	UseFrom string `db:"use_from" json:"useFrom"`
	//使用日志
	Des string `db:"des" json:"des"`
}