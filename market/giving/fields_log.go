package MarketGiving

import "time"

type FieldsLog struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//对表成员的用户ID
	// 和成员对等，可用于一次性推荐的记录处理
	UserID int64 `db:"user_id" json:"userID"`
	//奖励来源
	BindFrom string `db:"bind_from" json:"bindFrom"`
	//奖励ID
	BindID int64 `db:"bind_id" json:"bindID"`
	//奖励数量/金额
	Count int64 `db:"count" json:"count"`
	//奖励描述
	Des string `db:"des" json:"des"`
}
