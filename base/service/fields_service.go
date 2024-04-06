package BaseService

import "time"

// FieldsService 服务信息
type FieldsService struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//过期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
	//描述
	Description string `db:"description" json:"description" check:"des" min:"1" max:"500" empty:"true"`
	//事件订阅方式
	// server 服务器订阅; client 客户端订阅; all 服务器和客户端都订阅
	EventSubType string `db:"event_sub_type" json:"eventSubType" check:"intThan0"`
	//事件编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"300"`
	//事件类型
	// nats - NATS事件
	EventType string `db:"event_type" json:"eventType" check:"intThan0"`
	//事件地址
	// nats - 触发的地址
	EventURL string `db:"event_url" json:"eventURL" check:"des" min:"1" max:"600"`
	//事件固定参数
	// nats - 事件附带的固定参数，如果为空则根据流程阶段事件触发填入
	EventParams string `db:"event_params" json:"eventParams" check:"des" min:"1" max:"1000" empty:"true"`
}
