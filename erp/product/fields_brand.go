package ERPProduct

import "time"

type FieldsBrand struct {
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
	//品牌编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"300"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
	//关联分类
	// 可选，如果给予，则可用于检索所属的分类ID
	CategoryID int64 `db:"category_id" json:"categoryID" check:"id" empty:"true"`
}
