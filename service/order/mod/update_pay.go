package ServiceOrderMod

import (
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
)

// ArgsUpdatePayID 更新payID参数
type ArgsUpdatePayID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 可选，作为验证
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//日志
	//操作组织人员ID
	// 如果留空则说明为系统自动调整或创建人产生
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//支付ID
	PayID int64 `db:"pay_id" json:"payID" check:"id"`
}

// UpdatePayID 更新payID
func UpdatePayID(args ArgsUpdatePayID) {
	CoreNats.PushDataNoErr("service_order_pay_id", "/service/order/pay_id", "", args.ID, "", args)
}

// ArgsUpdatePrice 修改订单金额参数
type ArgsUpdatePrice struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 可选，作为验证
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//日志
	//操作组织人员ID
	// 如果留空则说明为系统自动调整或创建人产生
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//费用组成
	PriceList FieldsPrices `db:"price_list" json:"priceList"`
}

// UpdatePrice 修改订单金额参数
// 必须在付款之前修改
func UpdatePrice(args ArgsUpdatePrice) {
	CoreNats.PushDataNoErr("service_order_pay_price", "/service/order/pay_price", "", args.ID, "", args)
}
