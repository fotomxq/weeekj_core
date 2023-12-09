package IOTError

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"time"
)

type FieldsError struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//是否已经处理
	Done bool `db:"done" json:"done"`
	//是否推送了预警信息
	SendEW bool `db:"send_ew" json:"sendEW"`
	//组织ID
	// 设备所属的组织，也可能为0
	OrgID int64 `db:"org_id" json:"orgID"`
	//设备分组
	GroupID int64 `db:"group_id" json:"groupID"`
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID"`
	//错误标识码
	Code string `db:"code" json:"code"`
	//日志内容
	Content string `db:"content" json:"content"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
