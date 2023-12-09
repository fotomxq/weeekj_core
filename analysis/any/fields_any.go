package AnalysisAny

import "time"

//FieldsAny 统计数据支持
type FieldsAny struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	// 可留空
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	// 可留空
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//绑定ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//扩展参数1
	// 可选扩展参数，默认给0无视
	Param1 int64 `db:"params1" json:"params1" check:"id" empty:"true"`
	//扩展参数2
	// 可选扩展参数，默认给0无视
	Param2 int64 `db:"params2" json:"params2" check:"id" empty:"true"`
	//数据配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id"`
	//数据Hash
	Hash string `db:"hash" json:"hash"`
	//数据
	Data    int64  `db:"data" json:"data"`
	DataVal string `db:"data_val" json:"dataVal"`
}
