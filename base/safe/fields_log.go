package BaseSafe

import "time"

//FieldsLog 安全日志记录
type FieldsLog struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//来源系统
	System string `db:"system" json:"system"`
	//警告级别
	// 0 普通警告；1 中等警告，一些常见但容易混淆的安全问腿；2 高级警告，明显的安全问题警告
	Level int `db:"level" json:"level"`
	//触发IP
	IP string `db:"ip" json:"ip"`
	//触发用户
	UserID int64 `db:"user_id" json:"userID"`
	//触发商户
	OrgID int64 `db:"org_id" json:"orgID"`
	//事件日志信息
	Des string `db:"des" json:"des"`
}
