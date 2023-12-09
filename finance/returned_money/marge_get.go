package FinanceReturnedMoney

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"time"
)

// ArgsGetMargeList 获取汇总表列表参数
type ArgsGetMargeList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//回款公司ID
	CompanyID int64 `db:"company_id" json:"companyID" check:"id" empty:"true"`
	//是否需要已经回款参数
	NeedHaveAt bool `json:"needHaveAt" check:"bool"`
	HaveAt     bool `json:"haveAt" check:"bool"`
	//以下数据根据公司设置填入，该设计主要为保留历史数据记录
	//销售人员
	SellOrgBindID int64 `db:"sell_org_bind_id" json:"sellOrgBindID" check:"id" empty:"true"`
	//催收负责人
	ReturnOrgBindID int64 `db:"return_org_bind_id" json:"returnOrgBindID" check:"id" empty:"true"`
	//逾期天数
	NoReturnDay int `json:"noReturnDay" check:"intThan0" empty:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove"`
}

// GetMargeList 获取汇总表列表
func GetMargeList(args *ArgsGetMargeList) (dataList []FieldsMarge, dataCount int64, err error) {
	//组合条件处理
	maps := map[string]interface{}{}
	where := CoreSQL.GetDeleteSQL(args.IsRemove, "")
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.CompanyID > -1 {
		where = where + " AND company_id = :company_id"
		maps["company_id"] = args.CompanyID
	}
	if args.NeedHaveAt {
		where = CoreSQL.GetDeleteSQLField(args.HaveAt, where, "have_at")
	}
	if args.SellOrgBindID > -1 {
		where = where + " AND sell_org_bind_id = :sell_org_bind_id"
		maps["sell_org_bind_id"] = args.SellOrgBindID
	}
	if args.ReturnOrgBindID > -1 {
		where = where + " AND return_org_bind_id = :return_org_bind_id"
		maps["return_org_bind_id"] = args.ReturnOrgBindID
	}
	if args.NoReturnDay > -1 {
		needAt := CoreFilter.GetNowTimeCarbon().SubDays(args.NoReturnDay)
		where = where + " AND need_at <= :need_at"
		maps["need_at"] = needAt.Time
	}
	tableName := "finance_returned_money_marge"
	var rawList []FieldsMarge
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		"SELECT id FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	if err != nil {
		return
	}
	//重组数据
	for _, v := range rawList {
		vData := getMargeByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	//反馈
	return
}

// ArgsGetMargeByID 获取指定ID参数
type ArgsGetMargeByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//催收负责人
	ReturnOrgBindID int64 `db:"return_org_bind_id" json:"returnOrgBindID" check:"id" empty:"true"`
}

// GetMargeByID 获取指定ID
func GetMargeByID(args *ArgsGetMargeByID) (data FieldsMarge) {
	data = getMargeByID(args.ID)
	if data.ID < 1 {
		return
	}
	if CoreSQL.CheckTimeHaveData(data.DeleteAt) || !CoreFilter.EqID2(args.OrgID, data.OrgID) || !CoreFilter.EqID2(args.ReturnOrgBindID, data.ReturnOrgBindID) {
		data = FieldsMarge{}
		return
	}
	return
}

// GetMargeNoTakePrice 获取汇总表历史未回款金额总额
// needPrice 历史总等待回款金额
// returnPrice 历史总已经回款金额
// lastPrice 历史需继续回款的总金额
func GetMargeNoTakePrice(companyID int64, beforeAt time.Time) (needPrice int64, returnPrice int64, lastPrice int64) {
	//缓冲
	type cacheDataType struct {
		NeedPrice   int64
		ReturnPrice int64
		LastPrice   int64
	}
	var cacheData cacheDataType
	cacheMark := getMargeAnalysisDayCacheMark(companyID, CoreFilter.GetTimeToDefaultDate(beforeAt))
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &cacheData); err == nil && cacheData.NeedPrice > 0 {
		needPrice = cacheData.NeedPrice
		returnPrice = cacheData.ReturnPrice
		lastPrice = cacheData.LastPrice
		return
	}
	//计算数据
	err := Router2SystemConfig.MainDB.Get(&needPrice, "SELECT SUM(need_price) FROM finance_returned_money_marge WHERE delete_at < to_timestamp(1000000) AND company_id = $1 AND need_at <= $2 LIMIT 1", companyID, beforeAt)
	if err != nil {
		//无错误
	}
	err = Router2SystemConfig.MainDB.Get(&returnPrice, "SELECT SUM(have_price) FROM finance_returned_money_marge WHERE delete_at < to_timestamp(1000000) AND company_id = $1 AND need_at <= $2", companyID, beforeAt)
	if err != nil {
		//无错误
	}
	//计算等待回款金额
	lastPrice = needPrice - returnPrice
	//缓冲
	cacheData = cacheDataType{
		NeedPrice:   needPrice,
		ReturnPrice: returnPrice,
		LastPrice:   lastPrice,
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, cacheData, 18400)
	//反馈
	return
}

// DataGetMargeAnalysisCompany 获取公司的统筹信息组参数
type DataGetMargeAnalysisCompany struct {
	//公司ID
	CompanyID int64 `json:"companyID"`
	//历史总应付款
	AllNeedPrice int64 `json:"allNeedPrice"`
	//历史已付款
	AllReturnPrice int64 `json:"allReturnPrice"`
	//历史尚未付款
	AllLastPrice int64 `json:"allLastPrice"`
	//最新一期笔数
	LastLogCount int64 `json:"lastLogCount"`
	//发生逾期
	Last1Price int64 `json:"last1Price"`
	//历史超出30天的尚未付款
	Last30Price int64 `json:"last30Price"`
	//历史超出60天的尚未付款
	Last60Price int64 `json:"last60Price"`
	//历史超出90天的尚未付款
	Last90Price int64 `json:"last90Price"`
	//历史超出365天的尚未付款
	Last365Price int64 `json:"last365Price"`
	//当前超期状态
	// 0 没有应收; 1 存在应收尚未逾期; 2 存在应收已经逾期; 3 已经完成回款；4 存在逾期; 5 严重逾期30天; 6 违约60天; 7 违约90天; 8 违约365天
	ReturnStatus int `db:"return_status" json:"returnStatus"`
}

// GetMargeAnalysisCompany 获取公司的统筹信息组
func GetMargeAnalysisCompany(orgID int64, companyID int64) (data DataGetMargeAnalysisCompany) {
	//缓冲
	cacheMark := getMargeAnalysisCompanyCacheMark(companyID)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.CompanyID > 0 {
		return
	}
	// 关键节点时间
	nowAt := CoreFilter.GetNowTimeCarbon()
	//公司设置
	companySetData := GetCompanyByCompanyID(companyID, -1)
	//设置状态
	data.ReturnStatus = companySetData.ReturnStatus
	//获取相关统计数据
	data.CompanyID = companyID
	data.AllNeedPrice, data.AllReturnPrice, data.AllLastPrice = GetMargeNoTakePrice(companyID, nowAt.AddYear().Time)
	data.LastLogCount = getLogCountByCompanyID(orgID, companyID)
	_, _, data.Last1Price = GetMargeNoTakePrice(companyID, nowAt.Time)
	_, _, data.Last30Price = GetMargeNoTakePrice(companyID, nowAt.SubMonth().Time)
	_, _, data.Last60Price = GetMargeNoTakePrice(companyID, nowAt.SubMonths(2).Time)
	_, _, data.Last90Price = GetMargeNoTakePrice(companyID, nowAt.SubMonths(3).Time)
	_, _, data.Last365Price = GetMargeNoTakePrice(companyID, nowAt.SubYear().Time)
	//设置缓冲
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 21600)
	//反馈
	return
}

// 获取汇总表ID
func getMargeByID(id int64) (data FieldsMarge) {
	cacheMark := getMargeCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, company_id, need_price, need_at, have_price, have_at, sell_org_bind_id, return_org_bind_id, return_confirm_at, start_at FROM finance_returned_money_marge WHERE id = $1", id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 21600)
	return
}
