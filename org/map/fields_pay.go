package OrgMap

import "time"

type FieldsMapPay struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//支付成功时间
	FinishAt time.Time `db:"finish_at" json:"finishAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//地图ID
	MapID int64 `db:"map_id" json:"mapID"`
	//支付ID
	PayID int64 `db:"pay_id" json:"payID"`
	//购买的点击次数
	Count int64 `db:"count" json:"count"`
}
