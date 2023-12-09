package ServiceHousekeeping

import "time"

// FieldsConfig 服务配置
type FieldsConfig struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//排序
	Sort int `db:"sort" json:"sort"`
	//描述标题
	Title string `db:"title" json:"title"`
	//服务时间范围
	StartAt time.Time `db:"start_at" json:"startAt"`
	EndAt   time.Time `db:"end_at" json:"endAt"`
	//服务限制
	LimitCount int `db:"limit_count" json:"limitCount"`
	//累计预约个数
	// 不会清零，一直会累计下去，除非编辑过数据
	NowCount int `db:"now_count" json:"nowCount"`
}
