package BaseToken2

import "time"

// FieldsTokenS 短验证模块表
// 用于浏览器URL注入字符串形式，验证访问的有效性
type FieldsTokenS struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//绑定会话ID
	TokenID int64 `db:"token_id" json:"tokenID"`
	//匹配值
	Val string `db:"val" json:"val"`
}
