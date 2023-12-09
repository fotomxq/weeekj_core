package TMSTransport

import "time"

//FieldsAnalysis 统计
type FieldsAnalysis struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//配送人员
	BindID int64 `db:"bind_id" json:"bindID"`
	//客户档案ID
	InfoID int64 `db:"info_id" json:"infoID"`
	//客户用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//配送单ID
	TransportID int64 `db:"transport_id" json:"transportID"`
	//公里数
	KM int64 `db:"km" json:"km"`
	//总耗时
	OverTime int64 `db:"over_time" json:"overTime"`
	//评级
	// 1-5 级别
	Level int `db:"level" json:"level"`
}