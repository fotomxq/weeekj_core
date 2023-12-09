package BlogStuRead

import "time"

type FieldsLog struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//结束时间
	EndAt time.Time `db:"end_at" json:"endAt"`
	//学习时间长度
	RunTime int `db:"run_time" json:"runTime"`
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//学习内容
	ContentID int64 `db:"content_id" json:"contentID"`
}
