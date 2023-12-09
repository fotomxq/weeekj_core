package ServiceHealthSelf

import "time"

// FieldsLog 记录信息
type FieldsLog struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//健康码状态
	// 0 正常（绿）; 1 警告（黄）; 2 危险（红）
	HealthStatus int `db:"health_status" json:"healthStatus"`
	//健康码附加文件
	HealthFileID int64 `db:"health_file_id" json:"healthFileID"`
	//行程卡状态
	// 0 正常（绿）; 1 警告（黄）; 2 危险（红）
	TravelStatus int `db:"travel_status" json:"travelStatus"`
	//行程卡附加文件
	TravelFileID int64 `db:"travel_file_id" json:"travelFileID"`
	//体温
	// 小数点保留2位数x100
	BodyTemperature int `db:"body_temperature" json:"bodyTemperature"`
	//核酸报告截图
	NAReportFileID int64 `db:"na_report_file_id" json:"naReportFileID"`
	//核酸结果
	// 0 正常（阴性）; 1 异常（阳性）
	NAReportStatus int `db:"na_report_status" json:"naReportStatus"`
	//总的检查结果
	// 0 正常（绿）; 1 警告（黄）; 2 危险（红）
	Result int `db:"result" json:"result"`
}
