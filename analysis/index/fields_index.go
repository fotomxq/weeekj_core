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
	// 该阈值仅适用于单一提醒，如果是多层级提醒，建议单独开发
	Threshold int64 `db:"threshold" json:"threshold" index:"true"`
	//是否启用
	// 关闭后将不对该指标进行汇总运算
	IsEnable bool `db:"is_enable" json:"isEnable" index:"true" field_list:"true"`
	//指标类型
	// 指标类型方便程序做指标汇总时，进行识别，例如比率型可直接用于汇算；计数型、汇总型、平均型需根据业务特点进行汇总计算（需开发修正数据）
	// ratio 比率型; count 计数型; sum 汇总型; avg 平均型
	// ratio 比率型，必须是0-100%，建议存储为0.00-100.00的浮点数
	// count 计数型，必须是整数类型，建议存储为int64，一般用于次数记录
	// sum 汇总型，代表业务上是合计数据，建议存储为浮点数。一般用于财务指标
	// avg 平均型，代表业务上是平均数据，建议存储为浮点数。一般用于财务指标
	IndexType string `db:"index_type" json:"indexType" check:"des" min:"1" max:"50" index:"true" field_list:"true"`
	//指标方向
	// 用于描述指标好坏，可用于程序识别、业务识别，例如up表示指标越大越好，down表示指标越小越好
	// up 上升型; down 下降型
	IndexDirection string `db:"index_direction" json:"indexDirection" check:"des" min:"1" max:"50" index:"true" field_list:"true"`
}
