package FinancePhysicalPay

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	"time"
)

// FieldsPhysical 实物标的物
type FieldsPhysical struct {
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
	//名称
	Name string `db:"name" json:"name"`
	//可用的来源标的物
	BindFrom CoreSQLFrom.FieldsFrom `db:"bind_from" json:"bindFrom"`
	//置换一件商品需对应几个标的物
	NeedCount int64 `db:"need_count" json:"needCount"`
	//标的物市场总投放量限制
	LimitCount int64 `db:"limit_count" json:"limitCount"`
	//已经投放的数量
	TakeCount int64 `db:"take_count" json:"takeCount"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
