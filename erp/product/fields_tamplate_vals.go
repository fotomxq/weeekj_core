package ERPProduct

import "time"

// FieldsTemplateVals 模板字段值
type FieldsTemplateVals struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//模板ID
	TemplateID int64 `db:"template_id" json:"templateID" check:"id"`
	//字段类型
	FieldType string `db:"field_type" json:"fieldType" check:"des" min:"1" max:"300"`
	//字段默认值
	DefaultValue string `db:"default_value" json:"defaultValue" check:"des" min:"1" max:"1000"`
	//扩展参数
	ExtraParams string `db:"extra_params" json:"extraParams" check:"des" min:"1" max:"1000"`
}
