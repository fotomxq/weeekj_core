package BaseApprover

import "time"

// FieldsConfigItem 审批配置
type FieldsConfigItem struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id"`
	//审批顺序
	FlowOrder int `db:"flow_order" json:"flowOrder" check:"intThan0" empty:"true"`
	//组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//审批人用户ID
	// 用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
}
