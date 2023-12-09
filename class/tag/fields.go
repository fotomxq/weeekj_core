package ClassTag

import (
	"time"
)

// FieldsTag 标签配置
type FieldsTag struct {
	//基础
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//来源ID
	BindID int64 `db:"bind_id" json:"bindID"`
	//名称
	Name string `db:"name" json:"name"`
}
