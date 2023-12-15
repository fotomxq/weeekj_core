package ERPSaleOut

import (
	"errors"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLAddress "github.com/fotomxq/weeekj_core/v5/core/sql/address"
	ERPProduct "github.com/fotomxq/weeekj_core/v5/erp/product"
	MallCore "github.com/fotomxq/weeekj_core/v5/mall/core"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

type ArgsCreateLog struct {
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//创建人
	CreateOrgBindID int64 `db:"create_org_bind_id" json:"createOrgBindID" check:"id" empty:"true"`
	//进货对接人
	FromOrgBindID int64 `db:"from_org_bind_id" json:"fromOrgBindID" check:"id" empty:"true"`
	//销售对接人
	SellOrgBindID int64 `db:"sell_org_bind_id" json:"sellOrgBindID" check:"id" empty:"true"`
	//表单核对人
	ConfirmOrgBindID int64 `db:"confirm_org_bind_id" json:"confirmOrgBindID" check:"id" empty:"true"`
	//出货单类型
	// 0 普通出货单; 1 退货; 2 补差价
	SaleType int `db:"sale_type" json:"saleType" check:"intThan0" empty:"true"`
	//同步来源
	SyncSystem string `db:"sync_system" json:"syncSystem" check:"mark" empty:"true"`
	SyncHash   string `db:"sync_hash" json:"syncHash" check:"mark" empty:"true"`
	//订单来源
	OrderID int64 `db:"order_id" json:"orderID" check:"id" empty:"true"`
	//商品来源
	MallProductID int64 `db:"mall_product_id" json:"mallProductID" check:"id" empty:"true"`
	//产品来源
	ERPProductID int64 `db:"erp_product_id" json:"erpProductID" check:"id" empty:"true"`
	//供货商来源
	// 如果产品没有选择供货商，则查询绑定的供货商执行，如果还是没有将标记0，同时其他周边成本信息将标记为0
	FromCompanyID int64 `db:"from_company_id" json:"fromCompanyID" check:"id" empty:"true"`
	//购买人来源
	BuyUserID    int64 `db:"buy_user_id" json:"buyUserID" check:"id" empty:"true"`
	BuyCompanyID int64 `db:"buy_company_id" json:"buyCompanyID" check:"id" empty:"true"`
	//发货地址
	ToAddress CoreSQLAddress.FieldsAddress `db:"to_address" json:"toAddress"`
	//订单产品数量
	SellCount int64 `db:"sell_count" json:"sellCount" check:"int64Than0" empty:"true"`
	//折扣前原售价（不含税）
	SellRawPrice int64 `db:"sell_raw_price" json:"sellRawPrice" check:"price" empty:"true"`
	//折扣前原售价（含税）
	SellRawTaxPrice int64 `db:"sell_raw_tax_price" json:"sellRawTaxPrice" check:"price" empty:"true"`
	//最终售价（不含税）
	SellFinalPrice int64 `db:"sell_final_price" json:"sellFinalPrice" check:"price" empty:"true"`
	//最终售价（含税）
	SellFinalTaxPrice int64 `db:"sell_final_tax_price" json:"sellFinalTaxPrice" check:"price" empty:"true"`
	//享受折扣金额
	SellExPrice int64 `db:"sell_ex_price" json:"sellExPrice" check:"price" empty:"true"`
	//进货时供货商成本价（不含税）
	PurchasePrice int64 `db:"purchase_price" json:"purchasePrice" check:"price" empty:"true"`
	//进货时供货商成本价（含税）
	PurchaseTaxPrice int64 `db:"purchase_tax_price" json:"purchaseTaxPrice" check:"price" empty:"true"`
	//调整后的成本价（不含税）
	PurchaseFixPrice int64 `db:"purchase_fix_price" json:"purchaseFixPrice" check:"price" empty:"true"`
	//调整后的成本价（含税）
	PurchaseFixTaxPrice int64 `db:"purchase_fix_tax_price" json:"purchaseFixTaxPrice" check:"price" empty:"true"`
	//最终毛利（不含税）
	ProfitPrice int64 `db:"profit_price" json:"profitPrice" check:"price" empty:"true"`
	//最终毛利（含税）
	ProfitTaxPrice int64 `db:"profit_tax_price" json:"profitTaxPrice" check:"price" empty:"true"`
}

func CreateLog(args *ArgsCreateLog) (data FieldsLog, errCode string, err error) {
	//检查类型
	if !checkSaleType(args.SaleType) {
		errCode = "err_erp_sale_type"
		err = errors.New("sale type error")
		return
	}
	//修正时间
	if !CoreSQL.CheckTimeHaveData(args.CreateAt) {
		args.CreateAt = CoreFilter.GetNowTime()
	}
	//检查同步系统
	if !checkLogSyncSystem(args.OrgID, args.SyncSystem, args.SyncHash) {
		return
	}
	//创建数据
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "erp_sale_out", "INSERT INTO erp_sale_out (create_at, org_id, create_org_bind_id, from_org_bind_id, sell_org_bind_id, confirm_org_bind_id, sale_type, sync_system, sync_hash, order_id, mall_product_id, erp_product_id, from_company_id, buy_user_id, buy_company_id, to_address, sell_count, sell_raw_price, sell_raw_tax_price, sell_final_price, sell_final_tax_price, sell_ex_price, purchase_price, purchase_tax_price, purchase_fix_price, purchase_fix_tax_price, profit_price, profit_tax_price) VALUES (:create_at,:org_id,:create_org_bind_id,:create_org_bind_id,:sell_org_bind_id,:confirm_org_bind_id,:sale_type,:sync_system,:sync_hash,:order_id,:mall_product_id,:erp_product_id,:from_company_id,:buy_user_id,:buy_company_id,:to_address,:sell_count,:sell_raw_price,:sell_raw_tax_price,:sell_final_price,:sell_final_tax_price,:sell_ex_price,:purchase_price,:purchase_tax_price,:purchase_fix_price,:purchase_fix_tax_price,:profit_price,:profit_tax_price)", args, &data)
	if err != nil {
		errCode = "err_insert"
		return
	}
	//反馈
	return
}

// ArgsCreateLogMarge 联动创建参数
type ArgsCreateLogMarge struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//创建人
	CreateOrgBindID int64 `db:"create_org_bind_id" json:"createOrgBindID" check:"id" empty:"true"`
	//进货对接人
	FromOrgBindID int64 `db:"from_org_bind_id" json:"fromOrgBindID" check:"id" empty:"true"`
	//销售对接人
	SellOrgBindID int64 `db:"sell_org_bind_id" json:"sellOrgBindID" check:"id" empty:"true"`
	//表单核对人
	ConfirmOrgBindID int64 `db:"confirm_org_bind_id" json:"confirmOrgBindID" check:"id" empty:"true"`
	//订单来源
	OrderID int64 `db:"order_id" json:"orderID" check:"id"`
	//商品来源
	MallProductID int64 `db:"mall_product_id" json:"mallProductID" check:"id" empty:"true"`
	//产品来源
	ERPProductID int64 `db:"erp_product_id" json:"erpProductID" check:"id"`
	//供货商来源
	// 如果产品没有选择供货商，则查询绑定的供货商执行，如果还是没有将标记0，同时其他周边成本信息将标记为0
	FromCompanyID int64 `db:"from_company_id" json:"fromCompanyID" check:"id" empty:"true"`
	//购买人来源
	BuyUserID    int64 `db:"buy_user_id" json:"buyUserID" check:"id" empty:"true"`
	BuyCompanyID int64 `db:"buy_company_id" json:"buyCompanyID" check:"id" empty:"true"`
	//发货地址
	ToAddress CoreSQLAddress.FieldsAddress `db:"to_address" json:"toAddress"`
	//订单产品数量
	SellCount int64 `db:"sell_count" json:"sellCount" check:"int64Than0" empty:"true"`
}

// CreateLogMarge 联动创建
func CreateLogMarge(args *ArgsCreateLogMarge) (data FieldsLog, errCode string, err error) {
	//查询商品
	mallProductData, _ := MallCore.GetProduct(&MallCore.ArgsGetProduct{
		ID:    args.MallProductID,
		OrgID: args.OrgID,
	})
	if mallProductData.ID < 1 {
		errCode = "err_no_data"
		err = errors.New("mall product not exist")
		return
	}
	mallProductPrice := mallProductData.Price
	if CoreSQL.CheckTimeThanNow(mallProductData.PriceExpireAt) {
		mallProductPrice = mallProductData.PriceReal
	}
	//查询产品
	erpProductData, _ := ERPProduct.GetProductByID(&ERPProduct.ArgsGetProductByID{
		ID:    args.ERPProductID,
		OrgID: args.OrgID,
	})
	if erpProductData.ID < 1 {
		errCode = "err_no_data"
		err = errors.New("erp product not exist")
		return
	}
	//查询供货商
	var erpFromCompanyData ERPProduct.FieldsProductCompany
	if args.FromCompanyID > 0 {
		erpFromCompanyData = ERPProduct.GetProductCompany(args.OrgID, erpProductData.ID, args.FromCompanyID)
		if erpFromCompanyData.ID < 1 {
			errCode = "err_erp_product_company"
			err = errors.New("erp product not have company")
			return
		}
	} else {
		//查询公司绑定的供货商
		if erpProductData.CompanyID > 0 {
			erpFromCompanyData = ERPProduct.GetProductCompany(args.OrgID, erpProductData.ID, erpProductData.CompanyID)
			if erpFromCompanyData.ID < 1 {
				errCode = "err_erp_product_company"
				err = errors.New("erp product not have company")
				return
			}
		} else {
			//根据最低成本价随机抽取一个供货商
			erpFromCompanyData = ERPProduct.GetProductCompanyRand(args.ERPProductID, 1)
		}
	}
	if erpFromCompanyData.ID < 1 {
		errCode = "err_erp_product_company"
		err = errors.New("erp product not have company")
		return
	}
	//创建数据
	data, errCode, err = CreateLog(&ArgsCreateLog{
		CreateAt:            time.Time{},
		OrgID:               args.OrgID,
		CreateOrgBindID:     args.CreateOrgBindID,
		FromOrgBindID:       args.FromOrgBindID,
		SellOrgBindID:       args.SellOrgBindID,
		ConfirmOrgBindID:    args.ConfirmOrgBindID,
		SaleType:            0,
		SyncSystem:          "",
		SyncHash:            "",
		OrderID:             args.OrderID,
		MallProductID:       args.MallProductID,
		ERPProductID:        args.ERPProductID,
		FromCompanyID:       erpFromCompanyData.CompanyID,
		BuyUserID:           args.BuyUserID,
		BuyCompanyID:        args.BuyCompanyID,
		ToAddress:           args.ToAddress,
		SellCount:           args.SellCount,
		SellRawPrice:        args.SellCount * mallProductData.Price,
		SellRawTaxPrice:     args.SellCount * mallProductData.Price,
		SellFinalPrice:      args.SellCount * mallProductPrice,
		SellFinalTaxPrice:   args.SellCount * mallProductPrice,
		SellExPrice:         args.SellCount * (mallProductPrice - mallProductData.Price),
		PurchasePrice:       args.SellCount * erpFromCompanyData.CostPrice,
		PurchaseTaxPrice:    args.SellCount * erpFromCompanyData.TaxCostPrice,
		PurchaseFixPrice:    args.SellCount * erpFromCompanyData.CostPrice,
		PurchaseFixTaxPrice: args.SellCount * erpFromCompanyData.TaxCostPrice,
		ProfitPrice:         args.SellCount*mallProductPrice - args.SellCount*erpFromCompanyData.CostPrice,
		ProfitTaxPrice:      args.SellCount*mallProductPrice - args.SellCount*erpFromCompanyData.TaxCostPrice,
	})
	//反馈
	return
}

// 检查同步系统
func checkLogSyncSystem(orgID int64, system string, hash string) bool {
	var id int64
	err := Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM erp_sale_out WHERE org_id = $1 AND sync_system = $2 AND sync_hash = $3", orgID, system, hash)
	if err != nil || id < 1 {
		return true
	}
	return false
}
