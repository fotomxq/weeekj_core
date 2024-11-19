package AnalysisIndexDimensions

import "time"

// FieldsDimensions 维度定义
type FieldsDimensions struct {
	// ID
	ID int64 `db:"id" json:"id" check:"id" unique:"true"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt" default:"now()"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt" default:"now()"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt" default:"0" index:"true"`
	//编码
	// 维度编码，用于程序内部识别
	Code string `db:"code" json:"code" check:"des" min:"1" max:"50" index:"true" field_list:"true"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300" empty:"true" field_search:"true" field_list:"true"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"0" max:"500" empty:"true" field_search:"true"`
	//约定Extend字段名称
	// 约定指标、指标值
	// 例如：extend1
	ExtendIndex string `db:"extend_index" json:"extendIndex" index:"true"`
}
