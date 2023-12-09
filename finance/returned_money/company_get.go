package FinanceReturnedMoney

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/golang-module/carbon"
	"github.com/lib/pq"
)

// ArgsGetCompanyList 获取公司设置列表参数
type ArgsGetCompanyList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//回款公司ID
	CompanyID int64 `db:"company_id" json:"companyID" check:"id" empty:"true"`
	//销售人员
	SellOrgBindID int64 `db:"sell_org_bind_id" json:"sellOrgBindID" check:"id" empty:"true"`
	//催收负责人
	ReturnOrgBindID int64 `db:"return_org_bind_id" json:"returnOrgBindID" check:"id" empty:"true"`
	//是否坏账
	NeedIsBan bool `json:"needIsBan" check:"bool"`
	IsBan     bool `db:"is_ban" json:"isBan" check:"bool"`
	//当前超期状态
	// -1 跳过; 0 没有应收; 1 存在应收尚未逾期; 2 存在应收已经逾期; 3 已经完成回款；4 存在逾期; 5 严重逾期30天; 6 违约60天; 7 违约90天; 8 违约365天
	ReturnStatus pq.Int64Array `db:"return_status" json:"returnStatus"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
}

// GetCompanyList 获取公司设置列表
func GetCompanyList(args *ArgsGetCompanyList) (dataList []FieldsCompany, dataCount int64, err error) {
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
	if args.SellOrgBindID > -1 {
		where = where + " AND sell_org_bind_id = :sell_org_bind_id"
		maps["sell_org_bind_id"] = args.SellOrgBindID
	}
	if args.ReturnOrgBindID > -1 {
		where = where + " AND return_org_bind_id = :return_org_bind_id"
		maps["return_org_bind_id"] = args.ReturnOrgBindID
	}
	if args.NeedIsBan {
		where = where + " AND is_ban = :is_ban"
		maps["is_ban"] = args.IsBan
	}
	if len(args.ReturnStatus) > 0 {
		where = where + " AND return_status = ANY(:return_status)"
		maps["return_status"] = args.ReturnStatus
	}
	tableName := "finance_returned_money_company"
	var rawList []FieldsCompany
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		"SELECT company_id "+"FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	if err != nil {
		return
	}
	//重组数据
	for _, v := range rawList {
		vData := getCompanyID(v.CompanyID)
		if vData.CompanyID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	//反馈
	return
}

// GetCompanyByCompanyID 通过公司获取数据
func GetCompanyByCompanyID(companyID int64, orgID int64) (data FieldsCompany) {
	var rawData FieldsCompany
	_ = Router2SystemConfig.MainDB.Get(&rawData, "SELECT company_id FROM finance_returned_money_company WHERE company_id = $1 AND ($2 < 1 OR org_id = $2) AND delete_at < to_timestamp(1000000)", companyID, orgID)
	if rawData.CompanyID < 1 {
		return
	}
	data = getCompanyID(rawData.CompanyID)
	return
}

// 获取数据
func getCompanyID(companyID int64) (data FieldsCompany) {
	cacheMark := getCompanyCacheMark(companyID)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, company_id, need_take_add_month, need_take_start_day, need_take_day, sell_org_bind_id, return_org_bind_id, is_ban, return_location, return_status FROM finance_returned_money_company WHERE company_id = $1", companyID)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 21600)
	return
}

// getCompanySubTime 获取公司上期的截止日期
func getCompanySubTime(companyID int64) (subEndAt carbon.Carbon) {
	data := GetCompanyByCompanyID(companyID, -1)
	if data.ID < 1 {
		return
	}
	nowAt := CoreFilter.GetNowTimeCarbon()
	switch data.NeedTakeDay {
	case -1:
		//每月最后一天
		subEndAt = nowAt.EndOfDay()
	case 0:
		//每月第一天
		subEndAt = nowAt.StartOfDay()
	default:
		//其他日期
		subEndAt = nowAt.SetDay(data.NeedTakeDay)
	}
	if subEndAt.Time.Unix() >= nowAt.Time.Unix() {
		subEndAt = subEndAt.SubMonth()
	}
	//反馈时间
	return
}

// getCompanyNextTime 获取公司下一轮时间节点
func getCompanyNextTime(orgID int64, companyID int64) (nextStartAt carbon.Carbon) {
	//获取公司回款信息
	data := GetCompanyByCompanyID(companyID, orgID)
	if data.ID < 1 {
		nextStartAt = CoreFilter.GetNowTimeCarbon().AddMonth()
		return
	}
	//当前时间
	nextStartAt = CoreFilter.GetNowTimeCarbon()
	//获取公司最后一次回款汇总表记录
	var margeData FieldsMarge
	err := Router2SystemConfig.MainDB.Get(&margeData, "SELECT id FROM finance_returned_money_marge WHERE org_id = $1 AND company_id = $2 ORDER BY create_at DESC LIMIT 1", data.OrgID, data.CompanyID)
	if err == nil && margeData.ID > 0 {
		margeData = getMargeByID(margeData.ID)
		//修正当前时间，作为下一轮时间
		if margeData.NeedAt.Unix() < CoreFilter.GetNowTime().Unix() {
			nextStartAt = CoreFilter.GetCarbonByTime(margeData.NeedAt)
		}
	}
	//获取下一轮时间节点
	// 如果时间小于当前时间，则继续获取
	for {
		if data.NeedTakeAddMonth > 0 {
			nextStartAt = nextStartAt.AddMonths(data.NeedTakeAddMonth)
		} else {
			nextStartAt = nextStartAt.AddMonth()
		}
		switch data.NeedTakeDay {
		case -1:
			//每月最后一天
			nextStartAt = nextStartAt.EndOfDay()
		case 0:
			//每月第一天
			nextStartAt = nextStartAt.StartOfDay()
		default:
			//其他日期
			nextStartAt = nextStartAt.SetDay(data.NeedTakeDay)
		}
		//如果超出当前时间，则跳出
		if nextStartAt.Time.Unix() >= CoreFilter.GetNowTime().Unix() {
			break
		}
	}
	//反馈时间
	return
}
