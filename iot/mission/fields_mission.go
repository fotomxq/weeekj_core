package IOTMission

import (
	"time"
)

type FieldsMission struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//过期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	//组织ID
	// 设备所属的组织，也可能为0
	OrgID int64 `db:"org_id" json:"orgID"`
	//设备分组
	GroupID int64 `db:"group_id" json:"groupID"`
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID"`
	//任务状态
	// 0 wait 等待发起 / 1 send 已经发送 / 2 success 已经完成 / 3 failed 已经失败 / 4 cancel 取消
	Status int `db:"status" json:"status"`
	//发送请求数据集合
	ParamsData []byte `db:"params_data" json:"paramsData"`
	//回收数据
	// 回收数据如果过大，将不会被存储到本地
	ReportData []byte `db:"report_data" json:"reportData"`
	//任务动作
	Action string `db:"action" json:"action"`
	//任务动作连接方案
	// mqtt_client MQTT单一推送 ; mqtt_group MQTT分组推送 ; none 驱动主动处理
	ConnectType string `db:"connect_type" json:"connectType"`
}