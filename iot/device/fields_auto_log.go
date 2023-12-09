package IOTDevice

import "time"

//FieldsAutoLog 触发累计记录
// 针对外部模块，可以查询本记录表，发现设备最近1小时的触发情况
// 注意超出1小时将自动销毁，避免挤占空间
type FieldsAutoLog struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//触发设备ID
	DeviceID int64 `db:"device_id" json:"deviceID"`
	//触发信息
	InfoID int64 `db:"info_id" json:"infoID"`
	//触发条件
	// 扩展参数mark
	Mark string `db:"mark" json:"mark"`
	// 等式
	// 0 等于; 1 小于; 2 大于; 3 不等于
	Eq int `db:"eq" json:"eq"`
	//条件值
	EqVal string `db:"eq_val" json:"eqVal"`
	//值
	Val string `db:"val" json:"val"`
}