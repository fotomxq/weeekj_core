package AnalysisUserVisit

import "time"

//FieldsWaitTime 客户停留时间分析
type FieldsWaitTime struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	// 每隔1小时统计一次
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	// 如果存在数据，则表明该数据隶属于指定组织
	// 组织依可查看该数据
	OrgID int64 `db:"org_id" json:"orgID"`
	//系统类型
	System string `db:"system" json:"system"`
	FromMark string `db:"from_mark" json:"fromMark"`
	FromID int64 `db:"from_id" json:"fromID"`
	//访问数量
	Count int64 `db:"count" json:"count"`
	//时间总长度
	WaitTime int64 `db:"wait_time" json:"waitTime"`
}
