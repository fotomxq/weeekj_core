package UserIntegral

import "time"

//FieldsIntegral 积分主表
type FieldsIntegral struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//所属组织
	// 如果为0，则代表平台方
	OrgID int64 `db:"org_id" json:"orgID"`
	//所属用户
	UserID int64 `db:"user_id" json:"userID"`
	//积分
	Count int64 `db:"count" json:"count"`
}

//FieldsLog 积分变动表
type FieldsLog struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//所属组织
	// 如果为0，则代表平台方
	OrgID int64 `db:"org_id" json:"orgID"`
	//所属用户
	UserID int64 `db:"user_id" json:"userID"`
	//变动分数
	AddCount int64 `db:"add_count" json:"addCount"`
	//备注
	Des string `db:"des" json:"des"`
}