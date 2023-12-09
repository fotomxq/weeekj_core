package FinanceReturnedMoney

import "time"

type FieldsLog struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//回款单号
	SN string `db:"sn" json:"sn"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//回款公司ID
	CompanyID int64 `db:"company_id" json:"companyID"`
	//关联订单ID
	OrderID int64 `db:"order_id" json:"orderID"`
	//支付ID
	PayID int64 `db:"pay_id" json:"payID"`
	//关联其他第三方模块
	BindSystem string `db:"bind_system" json:"bindSystem"`
	BindID     int64  `db:"bind_id" json:"bindID"`
	BindMark   string `db:"bind_mark" json:"bindMark"`
	//是否回款, 否则为入账
	IsReturn bool `db:"is_return" json:"isReturn"`
	//是否发生退款
	HaveRefund bool `db:"have_refund" json:"haveRefund"`
	//回款金额
	Price int64 `db:"price" json:"price"`
	//备注历史
	Des string `db:"des" json:"des"`
}
