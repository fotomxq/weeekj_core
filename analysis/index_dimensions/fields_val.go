package AnalysisIndexDimensions

import "time"

// FieldsVal 维度值定义
// 该结构体主要用于自定义的维度值枚举值定义，当维度没有定义来源表和字段，那么可使用本模块直接定义维度值
type FieldsVal struct {
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
	Code string `db:"code" json:"code" check:"des" min:"1" max:"600" index:"true" field_search:"true" field_list:"true"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"600" empty:"true" field_search:"true" field_list:"true"`
}
