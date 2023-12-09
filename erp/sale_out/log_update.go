package ERPSaleOut

import (
	"errors"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsUpdateLogPurchaseFixPrice 更新成本底价参数
type ArgsUpdateLogPurchaseFixPrice struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//调整后的成本价（不含税）
	PurchaseFixPrice int64 `db:"purchase_fix_price" json:"purchaseFixPrice" check:"price" empty:"true"`
	//调整后的成本价（含税）
	PurchaseFixTaxPrice int64 `db:"purchase_fix_tax_price" json:"purchaseFixTaxPrice" check:"price" empty:"true"`
}

// UpdateLogPurchaseFixPrice 更新成本底价
// 将自动核算新的毛利
func UpdateLogPurchaseFixPrice(args *ArgsUpdateLogPurchaseFixPrice) (err error) {
	data := GetLogByID(args.ID, args.OrgID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	data.ProfitPrice = data.SellFinalPrice - args.PurchaseFixPrice
	data.ProfitTaxPrice = data.SellFinalTaxPrice - args.PurchaseFixTaxPrice
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE erp_sale_out SET purchase_fix_price = :purchase_fix_price, purchase_fix_tax_price = :purchase_fix_tax_price, profit_price = :profit_price, profit_tax_price = :profit_tax_price WHERE id = :id", map[string]interface{}{
		"id":                     data.ID,
		"purchase_fix_price":     args.PurchaseFixPrice,
		"purchase_fix_tax_price": args.PurchaseFixTaxPrice,
		"profit_price":           data.ProfitPrice,
		"profit_tax_price":       data.ProfitTaxPrice,
	})
	if err != nil {
		return
	}
	deleteLogCache(data.ID)
	return
}
