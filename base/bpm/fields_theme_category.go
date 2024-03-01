package BaseBPM

import "time"

// FieldsThemeCategory 主题分类
type FieldsThemeCategory struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//创建时间
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	//更新时间
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
	//删除时间
	DeletedAt time.Time `db:"deleted_at" json:"deletedAt"`
	//主题名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
	//主题描述
	Description string `db:"description" json:"description" check:"des" min:"1" max:"500" empty:"true"`
}
