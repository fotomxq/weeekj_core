package ToolsAppUpdate

import "time"

//FieldsCount 总数统计
type FieldsCount struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//统计时间
	// 20200101 没有小时等
	DayTime time.Time `db:"day_time" json:"dayTime"`
	//组织ID
	// 设备所属的组织，也可能为0
	// 组织也可以发布自己的APP
	OrgID int64 `db:"org_id" json:"orgID"`
	//APP ID
	AppID int64 `db:"app_id" json:"appID"`
	//版本ID
	UpdateID int64 `db:"update_id" json:"updateID"`
	//次数
	Count int64 `db:"count" json:"count"`
}