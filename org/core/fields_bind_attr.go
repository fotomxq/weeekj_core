package OrgCoreCore

import "time"

type FieldsBindAttr struct {
	//ID
	ID int64 `db:"id" json:"id" unique:"true"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt" default:"now()"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt" default:"now()"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt" default:"0"`
	//所属企业ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" index:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" index:"true"`
	//组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" index:"true"`
	//标识码
	AttrCode string `db:"attr_code" json:"attrCode" index:"true"`
	//值
	AttrValue string `db:"attr_value" json:"attrValue"`
	//整数
	AttrInt int64 `db:"attr_int" json:"attrInt"`
	//浮点数
	AttrFloat float64 `db:"attr_float" json:"attrFloat"`
}
