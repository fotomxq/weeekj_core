package BaseConfig

import (
	"time"
)

type FieldsConfigType struct {
	//标识码
	Mark string `db:"mark" json:"mark"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//是否可以公开
	AllowPublic bool `db:"allow_public" json:"allowPublic"`
	//验证Hash
	UpdateHash string `db:"update_hash" json:"updateHash"`
	//名称
	Name string `db:"name" json:"name"`
	//分组
	GroupMark string `db:"group_mark" json:"groupMark"`
	//描述
	Des string `db:"des" json:"des"`
	//结构
	// string / bool / int / int64 / float64
	// 结构也可用于前端判定某个特殊的样式，如时间样式、过期时间样式等，程序内不做任何限定，只是标记
	ValueType string `db:"value_type" json:"valueType"`
	//值
	Value string `db:"value" json:"value"`
}
