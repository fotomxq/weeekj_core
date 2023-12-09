package OrgTimeSign

import "time"

//FieldsConfig 考勤设置
// 在考勤时间基础上进行条件性设置
type FieldsConfig struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//掌管该数据的组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
}
