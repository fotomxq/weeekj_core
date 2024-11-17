package AnalysisIndex

import "time"

// FieldsIndexParam 参数定义
// 仅适用于内置指标
type FieldsIndexParam struct {
	// ID
	ID int64 `db:"id" json:"id" check:"id" unique:"true"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt" default:"now()"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt" default:"now()"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt" default:"0" index:"true"`
	//指标ID
	IndexID int64 `db:"index_id" json:"indexID" check:"id" index:"true"`
	//参数编码
	// 用于程序内识别内置指标的参数
	Code string `db:"code" json:"code" check:"des" min:"1" max:"50" field_search:"true" field_list:"true"`
	//参数值
	ParamVal string `db:"param_val" json:"paramVal"`
}
