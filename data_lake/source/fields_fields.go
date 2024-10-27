package DataLakeSource

import "time"

// FieldsFields 表结构
type FieldsFields struct {
	//ID
	ID int64 `db:"id" json:"id" unique:"true"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt" default:"now()"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt" default:"now()"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt" default:"0"`
	//表ID
	TableID int64 `db:"table_id" json:"tableId" index:"true"`
	//字段名
	FieldName string `db:"field_name" json:"fieldName" index:"true" field_search:"true"`
	//提示名称
	FieldLabel string `db:"field_label" json:"fieldLabel" field_search:"true"`
	//字段表单类型
	InputType string `db:"input_type" json:"inputType" field_search:"true"`
	//字段表单长度
	InputLength int `db:"input_length" json:"inputLength"`
	//字段表单默认值
	InputDefault string `db:"input_default" json:"inputDefault"`
	//字段表单是否必填
	InputRequired bool `db:"input_required" json:"inputRequired"`
	//字段表单正则表达式
	InputPattern string `db:"input_pattern" json:"inputPattern"`
	//字段数据类型
	DataType string `db:"data_type" json:"dataType" field_search:"true"`
	//字段描述
	FieldDesc string `db:"field_desc" json:"fieldDesc" field_search:"true"`
}
