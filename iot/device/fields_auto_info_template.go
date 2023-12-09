package IOTDevice

import "time"

//FieldsAutoInfoTemplate info规则套用设计
type FieldsAutoInfoTemplate struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//名称
	Name string `db:"name" json:"name"`
	//冷却时间
	WaitTime int64 `db:"wait_time" json:"waitTime"`
	//触发条件
	// 扩展参数mark
	Mark string `db:"mark" json:"mark"`
	// 等式
	// 0 等于; 1 小于; 2 大于; 3 不等于
	Eq int `db:"eq" json:"eq"`
	//值
	Val string `db:"val" json:"val"`
	//发送任务指令
	// 留空则发送触发条件的数据包
	SendAction string `db:"send_action" json:"sendAction"`
	//发送参数
	ParamsData []byte `db:"params_data" json:"paramsData"`
}