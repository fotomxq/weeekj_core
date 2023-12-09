package BlogExam

import "time"

// FieldsTopicGroupLog 学习日志记录
type FieldsTopicGroupLog struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//日期
	CreateDay time.Time `db:"create_day" json:"createDay"`
	//耗时
	// 秒
	RunTime int `db:"run_time" json:"runTime"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
}
