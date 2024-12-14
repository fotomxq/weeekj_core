package AnalysisIndexSQL

import "time"

type FieldsSQL struct {
	// ID
	ID int64 `db:"id" json:"id" check:"id" unique:"true"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt" default:"now()"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt" default:"now()"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt" default:"0" index:"true"`
	//指标编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"50" index:"true" field_list:"true"`
	///////////////////////////////////////////////////////////////////////////////////////////////////
	// 维度和IndexVals模块一致
	// 由于本模块特殊性，每个SQL只能对应一个维度，否则需进行区分SQL
	///////////////////////////////////////////////////////////////////////////////////////////////////
	//扩展维度名称
	// 可建议特别的维度关系，例如特定供应商的数据、特定地区的数据等
	ExtendName string `db:"extend_name" json:"extendName" check:"des" min:"1" max:"50" index:"true" field_list:"true"`
	//维度值
	ExtendValue string `db:"extend_value" json:"extendValue" check:"des" min:"1" max:"-1" index:"true" field_list:"true"`
	///////////////////////////////////////////////////////////////////////////////////////////////////
	// 数据
	///////////////////////////////////////////////////////////////////////////////////////////////////
	//SQL内容
	SQLData string `db:"sql_data" json:"sqlData"`
}
