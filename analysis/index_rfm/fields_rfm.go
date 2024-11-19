package AnalysisIndexRFM

import "time"

type FieldsRFM struct {
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
	//年月
	YearM string `db:"year_m" json:"yearM" index:"true" field_list:"true"`
	///////////////////////////////////////////////////////////////////////////////////////////////////
	// 维度和IndexVals模块一致
	///////////////////////////////////////////////////////////////////////////////////////////////////
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
	///////////////////////////////////////////////////////////////////////////////////////////////////
	// 计算过程值
	///////////////////////////////////////////////////////////////////////////////////////////////////
	//R
	RVal float64 `db:"r_val" json:"rVal"`
	//RMin
	RMin float64 `db:"r_min" json:"rMin"`
	//RMax
	RMax float64 `db:"r_max" json:"rMax"`
	//R 权重
	RWeight float64 `db:"r_weight" json:"rWeight"`
	//F
	FVal float64 `db:"f_val" json:"fVal"`
	//FMin
	FMin float64 `db:"f_min" json:"fMin"`
	//FMax
	FMax float64 `db:"f_max" json:"fMax"`
	//F 权重
	FWeight float64 `db:"f_weight" json:"fWeight"`
	//M
	MVal float64 `db:"m_val" json:"mVal"`
	//MMin
	MMin float64 `db:"m_min" json:"mMin"`
	//MMax
	MMax float64 `db:"m_max" json:"mMax"`
	//M 权重
	MWeight float64 `db:"m_weight" json:"mWeight"`
	///////////////////////////////////////////////////////////////////////////////////////////////////
	// 计算结果
	///////////////////////////////////////////////////////////////////////////////////////////////////
	RFMResult float64 `db:"rfm_result" json:"rfmResult"`
}
