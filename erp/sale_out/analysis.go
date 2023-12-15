package ERPSaleOut

import (
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	CoreSQLTime2 "github.com/fotomxq/weeekj_core/v5/core/sql/time2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

type ArgsGetAnalysisMarge struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
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
	SaleType int `db:"sale_type" json:"saleType"`
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
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//时间范围
	BetweenAt CoreSQLTime2.DataCoreTime `json:"betweenAt"`
}

type DataGetAnalysisMarge struct {
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

func GetAnalysisMarge(args *ArgsGetAnalysisMarge) (data DataGetAnalysisMarge) {
	//获取缓冲
	cacheMark := getAnalysisCacheMark(CoreFilter.GetSha1Str(fmt.Sprint(args)))
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.SellCount > 0 {
		return
	}
	//请求数据
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT SUM(sell_count) as sell_count, SUM(sell_raw_price) as sell_raw_price, SUM(sell_raw_tax_price) as sell_raw_tax_price, SUM(sell_final_price) as sell_final_price, SUM(sell_final_tax_price) as sell_final_tax_price, SUM(sell_ex_price) as sell_ex_price, SUM(purchase_price) as purchase_price, SUM(purchase_tax_price) as purchase_tax_price, SUM(purchase_fix_price) as purchase_fix_price, SUM(purchase_fix_tax_price) as purchase_fix_tax_price, SUM(profit_price) as profit_price, SUM(profit_tax_price) as profit_tax_price FROM erp_sale_out WHERE (($1 = true AND delete_at >= to_timestamp(1000000)) OR ($1 = false AND delete_at < to_timestamp(1000000))) AND org_id = $2 AND ($3 < 0 OR create_org_bind_id = $3) AND ($4 < 0 OR from_org_bind_id = $4) AND ($5 < 0 OR sell_org_bind_id = $5) AND ($6 < 0 OR confirm_org_bind_id = $6) AND ($7 < 0 OR sale_type = $7) AND ($8 = '' OR sync_system = $8) AND ($9 = '' OR sync_hash = $9) AND ($10 < 0 OR order_id = $10) AND ($11 < 0 OR mall_product_id = $11) AND ($12 < 0 OR erp_product_id = $12) AND ($13 < 0 OR from_company_id = $13) AND ($14 < 0 OR buy_user_id = $14) AND ($15 < 0 OR buy_company_id = $15) AND create_at >= $16 AND create_at <= $17", args.IsRemove, args.OrgID, args.CreateOrgBindID, args.FromOrgBindID, args.SellOrgBindID, args.ConfirmOrgBindID, args.SaleType, args.SyncSystem, args.SyncHash, args.OrderID, args.MallProductID, args.ERPProductID, args.FromCompanyID, args.BuyUserID, args.BuyCompanyID, args.BetweenAt.MinTime, args.BetweenAt.MaxTime)
	if err != nil {
		return
	}
	//保存缓冲
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 1800)
	//反馈
	return
}

// ArgsGetAnalysisBuyCompanySort 获取公司消费排名参数
type ArgsGetAnalysisBuyCompanySort struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DataGetAnalysisBuyCompanySort 获取公司消费排名数据
type DataGetAnalysisBuyCompanySort struct {
	//购买公司ID
	BuyCompanyID int64 `db:"buy_company_id" json:"buyCompanyID"`
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

// GetAnalysisBuyCompanySort 获取公司消费排名
func GetAnalysisBuyCompanySort(args *ArgsGetAnalysisBuyCompanySort) (dataList []DataGetAnalysisBuyCompanySort) {
	//获取缓冲
	cacheMark := getAnalysisBuyCompanySortCacheMark(args.OrgID, CoreFilter.GetSha1Str(fmt.Sprint(args)))
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &dataList); err == nil && len(dataList) > 0 {
		return
	}
	//请求数据
	err := Router2SystemConfig.MainDB.Select(&dataList, fmt.Sprint("SELECT buy_company_id, SUM(sell_count) as sell_count, SUM(sell_raw_price) as sell_raw_price, SUM(sell_raw_tax_price) as sell_raw_tax_price, SUM(sell_final_price) as sell_final_price, SUM(sell_final_tax_price) as sell_final_tax_price, SUM(sell_ex_price) as sell_ex_price, SUM(purchase_price) as purchase_price, SUM(purchase_tax_price) as purchase_tax_price, SUM(purchase_fix_price) as purchase_fix_price, SUM(purchase_fix_tax_price) as purchase_fix_tax_price, SUM(profit_price) as profit_price, SUM(profit_tax_price) as profit_tax_price FROM erp_sale_out WHERE org_id = $1 AND delete_at < to_timestamp(1000000) GROUP BY buy_company_id ORDER BY SUM(", CoreSQLPages.FilterSQLSort(args.Pages.Sort, []string{"sell_count", "sell_raw_price", "sell_raw_tax_price", "sell_final_price", "sell_final_tax_price", "sell_ex_price", "purchase_price", "purchase_tax_price", "purchase_fix_price", "purchase_fix_tax_price", "profit_price", "profit_tax_price"}), ") ", CoreSQLPages.GetSQLDesc(args.Pages.Desc), " LIMIT $2 OFFSET $3"), args.OrgID, args.Pages.Max, (args.Pages.Page-1)*args.Pages.Max)
	if err != nil {
		return
	}
	//保存缓冲
	Router2SystemConfig.MainCache.SetStruct(cacheMark, dataList, 1800)
	//反馈
	return
}

// DataGetAnalysisBuyCompanyMarge 获取购买公司的最近统计数据包
type DataGetAnalysisBuyCompanyMarge struct {
	//最近30天数据
	Last30Day []DataGetAnalysisBuyCompanyDay `json:"last30Day"`
}

// GetAnalysisBuyCompanyMarge 获取购买公司的最近统计数据包
func GetAnalysisBuyCompanyMarge(orgID int64, buyCompanyID int64) (data DataGetAnalysisBuyCompanyMarge) {
	startAt := CoreFilter.GetNowTimeCarbon().SubMonth()
	endAt := CoreFilter.GetNowTimeCarbon()
	for {
		if startAt.Time.Unix() > endAt.Time.Unix() {
			break
		}
		data.Last30Day = append(data.Last30Day, getAnalysisBuyCompanyDay(orgID, buyCompanyID, CoreFilter.GetTimeToDefaultDate(startAt.Time)))
		startAt = startAt.AddDay()
	}
	return
}

// DataGetAnalysisBuyCompanyDay 获取指定日期的公司购买聚合数据数据
type DataGetAnalysisBuyCompanyDay struct {
	//时间
	DayAt string `json:"dayAt"`
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

// getAnalysisBuyCompanyDay 获取指定日期的公司购买聚合数据
func getAnalysisBuyCompanyDay(orgID int64, buyCompanyID int64, nowDay string) (data DataGetAnalysisBuyCompanyDay) {
	//记录时间
	data.DayAt = nowDay
	//分析时间
	nowDayAt, err := CoreFilter.GetTimeByDefault(nowDay)
	if err != nil {
		return
	}
	minAt := CoreFilter.GetCarbonByTime(nowDayAt).StartOfDay()
	maxAt := CoreFilter.GetCarbonByTime(nowDayAt).EndOfDay()
	//获取缓冲
	cacheMark := getAnalysisBuyCompanyDayCacheMark(orgID, buyCompanyID, CoreFilter.GetCarbonByTime(nowDayAt).Format("20060102"))
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.SellCount > 0 {
		return
	}
	//请求数据
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT SUM(sell_count) as sell_count, SUM(sell_raw_price) as sell_raw_price, SUM(sell_raw_tax_price) as sell_raw_tax_price, SUM(sell_final_price) as sell_final_price, SUM(sell_final_tax_price) as sell_final_tax_price, SUM(sell_ex_price) as sell_ex_price, SUM(purchase_price) as purchase_price, SUM(purchase_tax_price) as purchase_tax_price, SUM(purchase_fix_price) as purchase_fix_price, SUM(purchase_fix_tax_price) as purchase_fix_tax_price, SUM(profit_price) as profit_price, SUM(profit_tax_price) as profit_tax_price FROM erp_sale_out WHERE org_id = $1 AND delete_at < to_timestamp(1000000) AND buy_company_id = $2 AND create_at >= $3 AND create_at <= $4", orgID, buyCompanyID, minAt.Time, maxAt.Time)
	if err != nil {
		return
	}
	//保存缓冲
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 2592000)
	//反馈
	return
}
