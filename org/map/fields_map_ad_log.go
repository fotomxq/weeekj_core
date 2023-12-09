package OrgMap

import "time"

// FieldsMapAdLog 地图广告点击收益记录
type FieldsMapAdLog struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//查看时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//完成时间
	FinishAt time.Time `db:"finish_at" json:"finishAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//点击用户
	ClickUserID int64 `db:"click_user_id" json:"clickUserID"`
	//地图ID
	MapID int64 `db:"map_id" json:"mapID"`
	//扣除的点击次数
	Count int64 `db:"integral_count" json:"count"`
	// 奖励金额
	Bonus int64 `db:"bonus" json:"bonus"`
}
