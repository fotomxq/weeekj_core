package OrgActivity

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"time"
)

// FieldsActivity 活动安排表
// 商户将活动挂钩到不同的内容上，实现优惠或其他内容共赢
type FieldsActivity struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//活动ID
	ConfigID int64 `db:"config_id" json:"configID"`
	//参与目标系统
	FromSystem string `db:"from_system" json:"fromSystem"`
	//参与目标ID
	FromID int64 `db:"from_id" json:"fromID"`
	//开始时间
	StartAt time.Time `db:"start_at" json:"startAt"`
	//结束时间
	EndAt time.Time `db:"end_at" json:"endAt"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
