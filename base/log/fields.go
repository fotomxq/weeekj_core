package BaseLog

import "time"

type FieldsLog struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//上报主机名称
	Mark string `db:"mark" json:"mark"`
	//上报主机IP
	IP string `db:"ip" json:"ip"`
	//日志类型
	LogType string `db:"log_type" json:"logType"`
	//时间类型
	// YYYY-MM-DD_HH | YYYY-MM-DD
	TimeType string `db:"time_type" json:"timeType"`
	//日志尺寸
	Size int `db:"size" json:"size"`
	//日志内容
	Contents string `db:"contents" json:"contents"`
}
