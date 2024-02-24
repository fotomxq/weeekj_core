package BaseUnit

import "time"

type FieldsUnit struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"100"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
}
