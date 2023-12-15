package ERPCore

import CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"

// FieldsComponentVal 节点组
type FieldsComponentVal struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//所属
	BindID int64 `db:"bind_id" json:"bindID"`
	//组件key
	// 单个节点内必须唯一
	Key string `db:"key" json:"key"`
	//展示顺序
	Sort int `db:"sort" json:"sort"`
	//组件类型
	ComponentType string `db:"component_type" json:"componentType"`
	//组件名称
	Name string `db:"name" json:"name"`
	//帮助描述
	HelpDes string `db:"help_des" json:"helpDes"`
	//组件默认值
	Val string `db:"val" json:"val"`
	//整数（内部记录用）
	ValInt64 int64 `db:"val_int64" json:"valInt64"`
	//浮点数（内部记录用）
	ValFloat64 float64 `db:"val_float64" json:"valFloat64"`
	//验证用的正则表达式
	CheckVal string `db:"check_val" json:"checkVal"`
	//是否必填
	IsRequire bool `db:"is_require" json:"isRequire"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
