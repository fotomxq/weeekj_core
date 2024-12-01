package AnalysisIndexFilter

import "time"

// FieldsFilter 筛选结果
type FieldsFilter struct {
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
	//年月日
	// 可选，如果不指定则不是时间维度的数据
	// 可任意持续，如年，或仅年月
	// 不建议构建小时及以下级别的指标
	// 同一个维度和时间范围，仅会存在一个数据，否则将覆盖
	YearMD string `db:"year_md" json:"yearMD" index:"true" field_list:"true"`
	//扩展维度1
	// 可建议特别的维度关系，例如特定供应商的数据、特定地区的数据等
	Extend1 string `db:"extend1" json:"extend1" index:"true" field_list:"true"`
	//扩展维度2
	Extend2 string `db:"extend2" json:"extend2" index:"true" field_list:"true"`
	//扩展维度3
	Extend3 string `db:"extend3" json:"extend3" index:"true" field_list:"true"`
	//扩展维度4
	Extend4 string `db:"extend4" json:"extend4" index:"true" field_list:"true"`
	//扩展维度5
	Extend5 string `db:"extend5" json:"extend5" index:"true" field_list:"true"`
	//追踪来源表名称
	FromTable string `db:"from_table" json:"fromTable" check:"des" min:"1" max:"100" index:"true" field_list:"true"`
	//追踪来源表ID
	FromID int64 `db:"from_id" json:"fromID" check:"id" index:"true" field_list:"true"`
	//数据值1
	// 用于标记筛选出结果后的一些数值内容
	Val1 float64 `db:"val1" json:"val1"`
	//数据值2
	Val2 float64 `db:"val2" json:"val2"`
	//数据值3
	Val3 float64 `db:"val3" json:"val3"`
	//数据值4
	Val4 float64 `db:"val4" json:"val4"`
	//数据值5
	Val5 float64 `db:"val5" json:"val5"`
	//数据值6
	Val6 float64 `db:"val6" json:"val6"`
	//数据值7
	Val7 float64 `db:"val7" json:"val7"`
	//数据值8
	Val8 float64 `db:"val8" json:"val8"`
	//数据值9
	Val9 float64 `db:"val9" json:"val9"`
	//数据值10
	Val10 float64 `db:"val10" json:"val10"`
	//备注值1
	// 用于标记一些字符串信息，也可以当作备注使用
	Des1 string `db:"des1" json:"des1" check:"des" min:"1" max:"3000" empty:"true" index:"true" field_list:"true" field_search:"true"`
	//备注值2
	Des2 string `db:"des2" json:"des2" check:"des" min:"1" max:"3000" empty:"true" index:"true" field_list:"true" field_search:"true"`
	//备注值3
	Des3 string `db:"des3" json:"des3" check:"des" min:"1" max:"3000" empty:"true" index:"true" field_list:"true" field_search:"true"`
	//备注值4
	Des4 string `db:"des4" json:"des4" check:"des" min:"1" max:"3000" empty:"true" index:"true" field_list:"true" field_search:"true"`
	//备注值5
	Des5 string `db:"des5" json:"des5" check:"des" min:"1" max:"3000" empty:"true" index:"true" field_list:"true" field_search:"true"`
}
