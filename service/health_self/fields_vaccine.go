package ServiceHealthSelf

import "time"

// FieldsVaccine 疫苗接种记录
type FieldsVaccine struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//疫苗名称
	Name string `db:"name" json:"name"`
	//接种地点
	Address string `db:"address" json:"address"`
}
