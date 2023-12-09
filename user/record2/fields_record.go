package UserRecord2

import "time"

type FieldsRecord struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID"`
	//成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//系统来源
	System string `db:"system" json:"system"`
	//影响ID
	ModID int64 `db:"mod_id" json:"modID"`
	//操作内容标识码
	Mark string `db:"mark" json:"mark"`
	//操作内容概述
	Des string `db:"des" json:"des"`
}
