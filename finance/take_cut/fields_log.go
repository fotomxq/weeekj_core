package FinanceTakeCut

import "time"

//FieldsLog 抽抽日志
type FieldsLog struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//订单金额
	OrderPrice int64 `db:"order_price" json:"orderPrice" check:"price"`
	//抽成金额
	CutPrice int64 `db:"cut_price" json:"cutPrice" check:"price"`
	//抽成时的配置比例
	// 单位10万分之%，例如300000对应3%；3对应0.00003%
	CutPriceProportion int64 `db:"cut_price_proportion" json:"cutPriceProportion"`
	//关联的订单
	OrderID int64 `db:"order_id" json:"orderID" check:"id"`
}
