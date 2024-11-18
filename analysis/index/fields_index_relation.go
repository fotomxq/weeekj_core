package AnalysisIndex

import "time"

// FieldsIndexRelation 指标组合关系
type FieldsIndexRelation struct {
	// ID
	ID int64 `db:"id" json:"id" check:"id" unique:"true"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt" default:"now()"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt" default:"now()"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt" default:"0" index:"true"`
	//指标ID
	// 上级指标
	IndexID int64 `db:"index_id" json:"indexID" check:"id" field_list:"true"`
	//关联指标
	// 禁出现嵌套关系，系统将检查并报错
	RelationIndexID int64 `db:"relation_index_id" json:"relationIndexID" check:"id" field_list:"true"`
	//指标权重占比
	// 同一个indexID下，所有指标的权重总和必须为1，否则计算的结果将出现异常
	Weight int64 `db:"weight" json:"weight" check:"int64Than0"`
	//算法自动权重
	// 同一个indexID下，所有指标的权重总和必须为1，否则计算的结果将出现异常
	AutoWeight int64 `db:"auto_weight" json:"autoWeight" check:"int64Than0"`
	//是否启动算法自动权重
	IsAutoWeight bool `db:"is_auto_weight" json:"isAutoWeight" field_list:"true"`
}
