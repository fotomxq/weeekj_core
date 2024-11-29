package AnalysisSignatureLibrary

import "time"

// FieldsLib 相关性识别模块
type FieldsLib struct {
	// ID
	ID int64 `db:"id" json:"id" check:"id" unique:"true"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt" default:"now()"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt" default:"now()"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt" default:"0" index:"true"`
	//算法模型类型
	// 1.皮尔森相关系数 CoreMathArraySimilarityPPMCC
	// 2.斯皮尔曼相关系数 CoreMathArraySimilaritySpearman
	LibType string `db:"lib_type" json:"libType" check:"des" min:"1" max:"50" index:"true"`
	//指标1编码
	Code1 string `db:"code1" json:"code1" check:"des" min:"1" max:"50" index:"true"`
	//指标2编码
	Code2 string `db:"code2" json:"code2" check:"des" min:"1" max:"50" index:"true"`
	//指标时间范围
	MinYearMD string `db:"min_year_md" json:"minYearMD" index:"true"`
	MaxYearMD string `db:"max_year_md" json:"maxYearMD" index:"true"`
	//相似度得分
	Score float64 `db:"score" json:"score"`
}
