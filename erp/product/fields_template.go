package ERPProduct

import "time"

// FieldsTemplate 产品模板
type FieldsTemplate struct {
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
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
	//插槽主题ID
	// BPM模块插槽主题ID，用于关联插槽主题，产品会自动使用相关的插槽用于表单实现
	BPMThemeID int64 `db:"bpm_theme_id" json:"bpmThemeID" check:"id" empty:"true"`
}
