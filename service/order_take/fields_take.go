package OrderTake

import (
	"time"
)

type FieldsTake struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//订单ID
	OrderID int64 `db:"order_id" json:"orderID"`
	//自提代码
	TakeCode string `db:"take_code" json:"takeCode"`
}
