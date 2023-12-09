package BaseSaving

import "time"

type FieldsSaving struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//过期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	//标识码
	Mark string `db:"mark" json:"mark"`
	//数据集合
	Val string `db:"val" json:"val"`
}