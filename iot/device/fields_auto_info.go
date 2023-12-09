package IOTDevice

import "time"

//FieldsAutoInfo 设备扩展信息传递表
// 当设备修改某个状态值，符合条件的将自动通知对应设备
type FieldsAutoInfo struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//触发设备ID
	DeviceID int64 `db:"device_id" json:"deviceID"`
	//采用模版
	// 如果存在模版ID，自定义触发条件将无效
	TemplateID int64 `db:"template_id" json:"templateID"`
	//触发条件
	// 扩展参数mark
	Mark string `db:"mark" json:"mark"`
	// 等式
	// 0 等于; 1 小于; 2 大于; 3 不等于
	Eq int `db:"eq" json:"eq"`
	//值
	Val string `db:"val" json:"val"`
	//冷却时间
	WaitTime int64 `db:"wait_time" json:"waitTime"`
	//反馈设备ID
	// 如果没有指定反馈设备ID
	ReportDeviceID int64 `db:"report_device_id" json:"reportDeviceID"`
	//发送任务指令
	// 留空则发送触发条件的数据包
	SendAction string `db:"send_action" json:"sendAction"`
	//发送参数
	ParamsData []byte `db:"params_data" json:"paramsData"`
}
