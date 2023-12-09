package ClassConfig

import "time"

//FieldsConfigDefault 默认配置项
type FieldsConfigDefault struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//标识码
	Mark string `db:"mark" json:"mark"`
	//名称
	Name string `db:"name" json:"name"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//是否可以公开
	AllowPublic bool `db:"allow_public" json:"allowPublic"`
	//是否允许自身查看
	AllowSelfView bool `db:"allow_self_view" json:"allowSelfView"`
	//是否允许自身修改
	AllowSelfSet bool `db:"allow_self_set" json:"allowSelfSet"`
	//结构
	// 0 string / 1 bool / 2 int / 3 int64 / 4 float64
	// 5 time 时间 / 6 daytime 带有日期的时间 / 7 unix 时间戳
	// 8 fileID 文件ID / 9 fileIDList 文件ID列
	// 10 userID 用户ID / 11 userIDList 用户ID列
	// 结构也可用于前端判定某个特殊的样式，如时间样式、过期时间样式等，程序内不做任何限定，只是标记
	ValueType int `db:"value_type" json:"valueType"`
	//正则表达式
	ValueCheck string `db:"value_check" json:"valueCheck"`
	//默认值
	ValueDefault string `db:"value_default" json:"valueDefault"`
}

//FieldsConfig 组织附加配置表
type FieldsConfig struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//绑定ID
	BindID int64 `db:"bind_id" json:"bindID"`
	//配置标识码
	Mark string `db:"mark" json:"mark"`
	//值
	Val string `db:"val" json:"val"`
}