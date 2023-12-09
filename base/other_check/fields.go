package BaseOtherCheck

import "time"

type FieldsCheck struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//过期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	//路由地址
	URL string `db:"url" json:"url"`
	//数据
	Data string `db:"data" json:"data"`
}