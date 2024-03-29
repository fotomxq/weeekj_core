package BaseBPM

import "time"

// FieldsTheme 主题
type FieldsTheme struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//创建时间
	CreatedAt time.Time `db:"create_at" json:"createdAt"`
	//更新时间
	UpdatedAt time.Time `db:"update_at" json:"updatedAt"`
	//删除时间
	DeletedAt time.Time `db:"delete_at" json:"deletedAt"`
	//所属主题分类
	CategoryID int64 `db:"category_id" json:"categoryId" check:"id"`
	//主题名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
	//主题描述
	Description string `db:"description" json:"description" check:"des" min:"1" max:"500" empty:"true"`
}
