package BaseIPAddr

import (
	"time"
)

//FieldsIPAddr 基本结构体
type FieldsIPAddr struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//过期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	//要匹配的IP地址
	IP string `db:"ip" json:"ip"`
	//是否为正则表达式
	// 正则表达式将匹配IP字段指定的范围区间
	// 注意发生错误将自动跳过
	IsMatch bool `db:"is_match" json:"isMatch"`
	//是否列入黑名单
	IsBan bool `db:"is_ban" json:"isBan"`
	//是否白名单
	IsWhite bool `db:"is_white" json:"isWhite"`
}
