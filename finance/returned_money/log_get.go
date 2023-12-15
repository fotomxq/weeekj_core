package FinanceReturnedMoney

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

// ArgsGetLogList 获取回款列表参数
type ArgsGetLogList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//回款公司ID
	CompanyID int64 `db:"company_id" json:"companyID" check:"id" empty:"true"`
	//关联订单ID
	OrderID int64 `db:"order_id" json:"orderID" check:"id" empty:"true"`
	//关联其他第三方模块
	BindSystem string `db:"bind_system" json:"bindSystem" check:"mark" empty:"true"`
	BindID     int64  `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	BindMark   string `db:"bind_mark" json:"bindMark"`
	//是否需要是否为回款参数
	NeedIsReturn bool `json:"needIsReturn" check:"bool"`
	IsReturn     bool `db:"is_return" json:"isReturn" check:"bool"`
	//时间范围
	MinAt string `json:"minAt" check:"defaultTime" empty:"true"`
	MaxAt string `json:"maxAt" check:"defaultTime" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetLogList 获取回款列表
func GetLogList(args *ArgsGetLogList) (dataList []FieldsLog, dataCount int64, err error) {
	//组合条件处理
	maps := map[string]interface{}{}
	where := ""
	if args.OrgID > -1 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.CompanyID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "company_id = :company_id"
		maps["company_id"] = args.CompanyID
	}
	if args.OrderID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "order_id = :order_id"
		maps["order_id"] = args.OrderID
	}
	if args.BindSystem != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "bind_system = :bind_system AND bind_id = :bind_id AND bind_mark = :bind_mark"
		maps["bind_system"] = args.BindSystem
		maps["bind_id"] = args.BindID
		maps["bind_mark"] = args.BindMark
	}
	if args.NeedIsReturn {
		if where != "" {
			where = where + " AND "
		}
		where = CoreSQL.GetDeleteSQLField(args.IsReturn, where, "is_return")
	}
	if args.MinAt != "" {
		if where != "" {
			where = where + " AND "
		}
		var minAt time.Time
		minAt, err = CoreFilter.GetTimeByDefault(args.MinAt)
		if err != nil {
			return
		}
		where = where + "create_at >= :min_at"
		maps["min_at"] = minAt
	}
	if args.MaxAt != "" {
		if where != "" {
			where = where + " AND "
		}
		var maxAt time.Time
		maxAt, err = CoreFilter.GetTimeByDefault(args.MaxAt)
		if err != nil {
			return
		}
		where = where + "create_at <= :max_at"
		maps["max_at"] = maxAt
	}
	if args.Search != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "(sn ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	if where == "" {
		where = "true"
	}
	tableName := "finance_returned_money_log"
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
		[]string{"id", "create_at"},
	)
	if err != nil {
		return
	}
	//重组数据
	for _, v := range rawList {
		vData := getLogByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	//反馈
	return
}

// GetLogLastReturnByCompanyID 获取公司最后一条回款日志
func GetLogLastReturnByCompanyID(companyID int64) (logData FieldsLog) {
	var rawData FieldsLog
	_ = Router2SystemConfig.MainDB.Get(&rawData, "SELECT id FROM finance_returned_money_log WHERE company_id = $1 AND is_return = true ORDER BY id DESC LIMIT 1", companyID)
	if rawData.ID < 1 {
		return
	}
	logData = getLogByID(rawData.ID)
	return
}

// 获取当期在途日志总数
func getLogCountByCompanyID(orgID int64, companyID int64) (count int64) {
	var margeData FieldsMarge
	err := Router2SystemConfig.MainDB.Get(&margeData, "SELECT id FROM finance_returned_money_marge WHERE org_id = $1 AND company_id = $2 AND delete_at < to_timestamp(1000000) AND need_at >= NOW()", orgID, companyID)
	if err != nil || margeData.ID < 1 {
		return
	}
	margeData = getMargeByID(margeData.ID)
	err = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM finance_returned_money_log WHERE org_id = $1 AND company_id = $2 AND create_at >= $3 AND create_at <= $4 AND is_return = false AND have_refund = false", margeData.OrgID, margeData.CompanyID, CoreFilter.GetNowTimeCarbon().StartOfDay().Time, margeData.NeedAt)
	if err != nil {
		return
	}
	return
}

// 获取数据
func getLogByID(id int64) (data FieldsLog) {
	cacheMark := getLogCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, sn, org_id, company_id, order_id, pay_id, bind_system, bind_id, bind_mark, is_return, price, des FROM finance_returned_money_log WHERE id = $1", id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 21600)
	return
}
