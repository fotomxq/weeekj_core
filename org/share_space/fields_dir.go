package OrgShareSpace

import "time"

type FieldsDir struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//所属人
	// 如果为0则为机构共享目录
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID"`
	//上级目录
	ParentID int64 `db:"parent_id" json:"parentID"`
	//名称
	Name string `db:"name" json:"name"`
}
