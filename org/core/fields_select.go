package OrgCoreCore

import "time"

//FieldsSelectOrg 选择组织列
// 记录用户选择对应的组织情况
type FieldsSelectOrg struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//最后一次选择时间
	LastAt time.Time `db:"last_at" json:"lastAt"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//组织中绑定的ID
	BindID int64 `db:"bind_id" json:"bindID"`
}
