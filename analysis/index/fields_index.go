package AnalysisIndex

import "time"

// FieldsIndex 指标定义
type FieldsIndex struct {
	// ID
	ID int64 `db:"id" json:"id" check:"id" unique:"true"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt" default:"now()"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt" default:"now()"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt" default:"0" index:"true"`
	//指标编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"50" index:"true"`
	//指标名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300" index:"true" field_search:"true" field_list:"true"`
	//是否内置
	// 前端应拦截内置指标的删除操作，以免影响系统正常运行，启动重启后将自动恢复，所以删除操作是无法生效的
	IsSystem bool `db:"is_system" json:"isSystem" index:"true" field_list:"true"`
	//指标描述
	Description string `db:"description" json:"description" check:"des" min:"1" max:"-1" field_search:"true" field_list:"true" empty:"true"`
	//指标决策建议
	Decision string `db:"decision" json:"decision" check:"des" min:"1" max:"-1" empty:"true" field_search:"true"`
	//指标预警阈值
	// 0-100，归一化后的数据，超出此范围将可触发预警事件记录
	Threshold int64 `db:"threshold" json:"threshold" index:"true"`
	//是否启用
	// 关闭后将不对该指标进行汇总运算
	IsEnable bool `db:"is_enable" json:"isEnable" index:"true" field_list:"true"`
}
