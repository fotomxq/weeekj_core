package ServiceOrderAnalysis

import (
	"time"
)

//FieldsOrg 订单组织标准统计
type FieldsOrg struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//统计周期
	// 完全相同的一个来源体系，1小时仅构建一条数据
	DayTime time.Time `db:"day_time" json:"dayTime"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//创建来源
	// 0 商户平台下单、1 用户APP下单、2 微信小程序下单、3 支付宝小程序下单、4 电话系统下单、5 线下柜台下单、6 辅助设备下单、7 其他渠道
	CreateFrom int `db:"create_from" json:"createFrom"`
	//订单个数
	Count int64 `db:"count" json:"count"`
	//货物货币类型
	Currency int `db:"currency" json:"currency"`
	//金额合计
	Price int64 `db:"price" json:"price"`
}

//FieldsUser 订单用户标准统计
type FieldsUser struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//统计周期
	// 完全相同的一个来源体系，1小时仅构建一条数据
	DayTime time.Time `db:"day_time" json:"dayTime"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//订单个数
	Count int64 `db:"count" json:"count"`
	//货物货币类型
	Currency int `db:"currency" json:"currency"`
	//金额合计
	Price int64 `db:"price" json:"price"`
}

//FieldsFromCount 商品数量统计
type FieldsFromCount struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//统计周期
	// 完全相同的一个来源体系，1小时仅构建一条数据
	DayTime time.Time `db:"day_time" json:"dayTime"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//商品来源
	FromSystem string `db:"from_system" json:"fromSystem"`
	//商品ID
	FromID int64 `db:"from_id" json:"fromID"`
	//商品数量
	BuyCount int64 `db:"buy_count" json:"buyCount"`
}

//FieldsOrgRefund 订单退单统计
//TODO: 暂未支持，由订单核心模块提供相关数据支持
type FieldsOrgRefund struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//统计周期
	// 完全相同的一个来源体系，1小时仅构建一条数据
	DayTime time.Time `db:"day_time" json:"dayTime"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//订单个数
	Count int64 `db:"count" json:"count"`
	//货物货币类型
	Currency int `db:"currency" json:"currency"`
	//金额合计
	Price int64 `db:"price" json:"price"`
}

//FieldsOrgExemption 优惠统计
type FieldsOrgExemption struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//统计周期
	// 完全相同的一个来源体系，1小时仅构建一条数据
	DayTime time.Time `db:"day_time" json:"dayTime"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//商品来源
	FromSystem string `db:"from_system" json:"fromSystem"`
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID"`
	//使用张数
	Count int64 `db:"count" json:"count"`
	//抵扣费用
	Price int64 `db:"price" json:"price"`
}
