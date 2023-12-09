package ERPSaleOut

import (
	CoreSQLAddress "gitee.com/weeekj/weeekj_core/v5/core/sql/address"
	"time"
)

type FieldsLog struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	// 作废时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//创建人
	CreateOrgBindID int64 `db:"create_org_bind_id" json:"createOrgBindID"`
	//进货对接人
	FromOrgBindID int64 `db:"from_org_bind_id" json:"fromOrgBindID"`
	//销售对接人
	SellOrgBindID int64 `db:"sell_org_bind_id" json:"sellOrgBindID"`
	//表单核对人
	ConfirmOrgBindID int64 `db:"confirm_org_bind_id" json:"confirmOrgBindID"`
	//出货单类型
	// 0 普通出货单; 1 退货; 2 补差价
	SaleType int `db:"sale_type" json:"saleType"`
	//同步来源
	SyncSystem string `db:"sync_system" json:"syncSystem"`
	SyncHash   string `db:"sync_hash" json:"syncHash"`
	//订单来源
	OrderID int64 `db:"order_id" json:"orderID"`
	//商品来源
	MallProductID int64 `db:"mall_product_id" json:"mallProductID"`
	//产品来源
	ERPProductID int64 `db:"erp_product_id" json:"erpProductID"`
	//供货商来源
	// 如果产品没有选择供货商，则查询绑定的供货商执行，如果还是没有将标记0，同时其他周边成本信息将标记为0
	FromCompanyID int64 `db:"from_company_id" json:"fromCompanyID"`
	//购买人来源
	BuyUserID    int64 `db:"buy_user_id" json:"buyUserID"`
	BuyCompanyID int64 `db:"buy_company_id" json:"buyCompanyID"`
	//发货地址
	ToAddress CoreSQLAddress.FieldsAddress `db:"to_address" json:"toAddress"`
	//订单产品数量
	SellCount int64 `db:"sell_count" json:"sellCount"`
	//折扣前原售价（不含税）
	SellRawPrice int64 `db:"sell_raw_price" json:"sellRawPrice"`
	//折扣前原售价（含税）
	SellRawTaxPrice int64 `db:"sell_raw_tax_price" json:"sellRawTaxPrice"`
	//最终售价（不含税）
	SellFinalPrice int64 `db:"sell_final_price" json:"sellFinalPrice"`
	//最终售价（含税）
	SellFinalTaxPrice int64 `db:"sell_final_tax_price" json:"sellFinalTaxPrice"`
	//享受折扣金额
	SellExPrice int64 `db:"sell_ex_price" json:"sellExPrice"`
	//进货时供货商成本价（不含税）
	PurchasePrice int64 `db:"purchase_price" json:"purchasePrice"`
	//进货时供货商成本价（含税）
	PurchaseTaxPrice int64 `db:"purchase_tax_price" json:"purchaseTaxPrice"`
	//调整后的成本价（不含税）
	PurchaseFixPrice int64 `db:"purchase_fix_price" json:"purchaseFixPrice"`
	//调整后的成本价（含税）
	PurchaseFixTaxPrice int64 `db:"purchase_fix_tax_price" json:"purchaseFixTaxPrice"`
	//最终毛利（不含税）
	ProfitPrice int64 `db:"profit_price" json:"profitPrice"`
	//最终毛利（含税）
	ProfitTaxPrice int64 `db:"profit_tax_price" json:"profitTaxPrice"`
}
