package FinancePhysicalPay

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	"time"
)

// FieldsLog 请求记录
type FieldsLog struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//配送人员
	BindID int64 `db:"bind_id" json:"bindID"`
	//用户
	UserID int64 `db:"user_id" json:"userID"`
	//支付分发渠道
	// order 订单 / tms 配送 / housekeeping 家政服务
	System string `db:"system" json:"system"`
	//标的物
	PhysicalID int64 `db:"physical_id" json:"physicalID"`
	//给予标的物数量
	PhysicalCount int64 `db:"physical_count" json:"physicalCount"`
	//货物来源标的物
	BindFrom CoreSQLFrom.FieldsFrom `db:"bind_from" json:"bindFrom"`
	//置换商品的数量
	BindCount int64 `db:"bind_count" json:"bindCount"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
