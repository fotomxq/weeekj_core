package ERPProduct

import "time"

// FieldsTemplateBind 模板和分类、品牌绑定关系
type FieldsTemplateBind struct {
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
	//分类ID
	// 关联ERPProduct.Sort模块
	CategoryID int64 `db:"category_id" json:"categoryID" check:"id" empty:"true"`
	//品牌ID
	BrandID int64 `db:"brand_id" json:"brandID" check:"id" empty:"true"`
}
