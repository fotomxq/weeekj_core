package ServiceDistribution

import "time"

//FieldsDistribution 分销商设置
type FieldsDistribution struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//分销商姓名
	Name string `db:"name" json:"name"`
	//绑定用户
	UserID int64 `db:"user_id" json:"userID"`
}
