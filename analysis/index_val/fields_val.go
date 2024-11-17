package AnalysisIndexVal

import "time"

// FieldsVal 指标值
type FieldsVal struct {
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
	//年月日
	// 可任意持续，如年，或仅年月
	// 不建议构建小时及以下级别的指标
	// 同一个维度和时间范围，仅会存在一个数据，否则将覆盖
	YearMD string `db:"year_md" json:"yearMD" index:"true"`
	//扩展维度1
	// 可建议特别的维度关系，例如特定供应商的数据、特定地区的数据等
	Extend1 string `db:"extend1" json:"extend1" index:"true"`
	//扩展维度2
	Extend2 string `db:"extend2" json:"extend2" index:"true"`
	//扩展维度3
	Extend3 string `db:"extend3" json:"extend3" index:"true"`
	//扩展维度4
	Extend4 string `db:"extend4" json:"extend4" index:"true"`
	//扩展维度5
	Extend5 string `db:"extend5" json:"extend5" index:"true"`
	//原始值
	ValRaw float64 `db:"val_raw" json:"valRaw" index:"true"`
	//归一化值
	ValNorm float64 `db:"val_norm" json:"valNorm" index:"true"`
	//是否为预测值
	IsForecast bool `db:"is_forecast" json:"isForecast"`
}
