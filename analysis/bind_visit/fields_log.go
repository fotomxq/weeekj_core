package AnalysisBindVisit

import "time"

type FieldsLog struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//用户
	UserID int64 `db:"user_id" json:"userID"`
	//来源模块
	BindSystem string `db:"bind_system" json:"bindSystem"`
	BindID     int64  `db:"bind_id" json:"bindID"`
}
