package BaseService

import "time"

// FieldsAnalysis 服务统计
type FieldsAnalysis struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//创建时间
	// 每小时创建一次
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//服务ID
	ServiceID int64 `db:"service_id" json:"serviceID" check:"id"`
	//服务端发送消息次数
	SendCount int `db:"send_count" json:"sendCount" check:"intThan0"`
	//服务端接收次数
	ReceiveCount int `db:"receive_count" json:"receiveCount" check:"intThan0"`
}
