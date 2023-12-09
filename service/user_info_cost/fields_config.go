package ServiceUserInfoCost

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"time"
)

// FieldsConfig 能耗费用配置
type FieldsConfig struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//名称
	Name string `db:"name" json:"name"`
	//房间场景值
	// 设备和房间绑定关系的mark值
	RoomBindMark string `db:"room_bind_mark" json:"roomBindMark"`
	//数据类型标识码
	// 遥感数据及传感器数据值
	SensorMark string `db:"sensor_mark" json:"sensorMark"`
	//计算方式
	// 0 合并计算，将时间阶段内的所有遥感数据合并进行统计计算
	// 1 平均值计算，将时间段内的数据平均化计算
	CountType int `db:"count_type" json:"countType"`
	//每小时能耗值
	EachUnit float64 `db:"each_unit" json:"eachUnit"`
	//每小时费用
	// 每累计产生EachUnit，将增加该金额1次，不足将不增加继续等待累计
	Currency  int   `db:"currency" json:"currency"`
	EachPrice int64 `db:"each_price" json:"eachPrice"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
