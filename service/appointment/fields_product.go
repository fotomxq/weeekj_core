package ServiceAppointment

import "time"

// FieldsProduct 预约项目
type FieldsProduct struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//发布时间
	PublishAt time.Time `db:"publish_at" json:"publishAt"`
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID"`
	//允许预约时间范围，每天时间点周期
	AppendMinAt string `db:"append_min_at" json:"appendMinAt"`
	AppendMaxAt string `db:"append_max_at" json:"appendMaxAt"`
}
