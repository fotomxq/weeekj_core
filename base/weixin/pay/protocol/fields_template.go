package BaseWeixinPayProtocol

import "time"

//FieldsTemplate 签约模版
type FieldsTemplate struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//组织ID
	// 设备所属的组织，也可能为0
	OrgID int64 `db:"org_id" json:"orgID"`
	//模版在微信的编号
	// 该编号在组织下唯一
	Code string `db:"code" json:"code"`
	//名称
	Name string `db:"name" json:"name"`
}