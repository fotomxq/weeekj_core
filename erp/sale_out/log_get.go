package ERPSaleOut

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	CoreSQLTime2 "github.com/fotomxq/weeekj_core/v5/core/sql/time2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

type ArgsGetLogList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
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
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

func GetLogList(args *ArgsGetLogList) (dataList []FieldsLog, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.CreateOrgBindID > -1 {
		where = where + " AND create_org_bind_id = :create_org_bind_id"
		maps["create_org_bind_id"] = args.CreateOrgBindID
	}
	if args.FromOrgBindID > -1 {
		where = where + " AND from_org_bind_id = :from_org_bind_id"
		maps["from_org_bind_id"] = args.FromOrgBindID
	}
	if args.SellOrgBindID > -1 {
		where = where + " AND sell_org_bind_id = :sell_org_bind_id"
		maps["sell_org_bind_id"] = args.SellOrgBindID
	}
	if args.ConfirmOrgBindID > -1 {
		where = where + " AND confirm_org_bind_id = :confirm_org_bind_id"
		maps["confirm_org_bind_id"] = args.ConfirmOrgBindID
	}
	if args.SaleType > -1 {
		where = where + " AND sale_type = :sale_type"
		maps["sale_type"] = args.SaleType
	}
	if args.SyncSystem != "" {
		where = where + " AND sync_system = :sync_system"
		maps["sync_system"] = args.SyncSystem
		if args.SyncHash != "" {
			where = where + " AND sync_hash = :sync_hash"
			maps["sync_hash"] = args.SyncHash
		}
	}
	if args.OrderID > -1 {
		where = where + " AND order_id = :order_id"
		maps["order_id"] = args.OrderID
	}
	if args.MallProductID > -1 {
		where = where + " AND mall_product_id = :mall_product_id"
		maps["mall_product_id"] = args.MallProductID
	}
	if args.ERPProductID > -1 {
		where = where + " AND erp_product_id = :erp_product_id"
		maps["erp_product_id"] = args.ERPProductID
	}
	if args.FromCompanyID > -1 {
		where = where + " AND from_company_id = :from_company_id"
		maps["from_company_id"] = args.FromCompanyID
	}
	if args.BuyUserID > -1 {
		where = where + " AND buy_user_id = :buy_user_id"
		maps["buy_user_id"] = args.BuyUserID
	}
	if args.BuyCompanyID > -1 {
		where = where + " AND buy_company_id = :buy_company_id"
		maps["buy_company_id"] = args.BuyCompanyID
	}
	if args.BetweenAt.MinTime != "" {
		where = where + " AND create_at >= :min_at"
		maps["min_at"] = args.BetweenAt.MinTime
	}
	if args.BetweenAt.MaxTime != "" {
		where = where + " AND create_at <= :max_at"
		maps["max_at"] = args.BetweenAt.MaxTime
	}
	if args.Search != "" {
		where = where + " AND (to_address ->> 'address' ILIKE '%' || :search || '%' OR to_address ->> 'name' ILIKE '%' || :search || '%' OR to_address ->> 'phone' ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "erp_sale_out"
	var rawList []FieldsLog
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		"SELECT id FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at", "sell_count", "sell_raw_price", "sell_raw_tax_price", "sell_final_price", "sell_final_tax_price", "sell_ex_price", "purchase_price", "purchase_tax_price", "purchase_fix_price", "purchase_fix_tax_price", "profit_price", "profit_tax_price"},
	)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getLogByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

func GetLogByID(id int64, orgID int64) (data FieldsLog) {
	data = getLogByID(id)
	if data.ID < 1 || !CoreFilter.EqID2(orgID, data.OrgID) {
		data = FieldsLog{}
		return
	}
	return
}

func getLogByID(id int64) (data FieldsLog) {
	cacheMark := getLogCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, create_org_bind_id, from_org_bind_id, sell_org_bind_id, confirm_org_bind_id, sale_type, sync_system, sync_hash, order_id, mall_product_id, erp_product_id, from_company_id, buy_user_id, buy_company_id, to_address, sell_count, sell_raw_price, sell_raw_tax_price, sell_final_price, sell_final_tax_price, sell_ex_price, purchase_price, purchase_tax_price, purchase_fix_price, purchase_fix_tax_price, profit_price, profit_tax_price FROM erp_sale_out WHERE id = $1", id)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 3600)
	return
}
