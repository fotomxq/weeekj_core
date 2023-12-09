package AnalysisAny2

import "time"

// FieldsConfig 统计配置
type FieldsConfig struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//数据标识码
	Mark string `db:"mark" json:"mark" check:"mark"`
	//颗粒度
	// 0 小时（默认）/ 1 1天 / 2 1周 / 3 1月 / 4 1年
	Particle int `db:"particle" json:"particle"`
	//归档天数
	// 自动生成的数据为30天
	// 必须指定，小于1则强制按照3天计算
	FileDay int `db:"file_day" json:"fileDay" check:"intThan0"`
}
