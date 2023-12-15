package ERPSaleOut

import (
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

func getLogCacheMark(id int64) string {
	return fmt.Sprint("erp:sale:out:log:id:", id)
}

func getAnalysisCacheMark(args string) string {
	return fmt.Sprint("erp:sale:out:analysis:marge:", args)
}

func getAnalysisBuyCompanySortCacheMark(orgID int64, args string) string {
	return fmt.Sprint("erp:sale:out:analysis:buy:company:sort:", orgID, ".", args)
}

func getAnalysisBuyCompanyDayCacheMark(orgID int64, buyCompanyID int64, dayTime string) string {
	return fmt.Sprint("erp:sale:out:analysis:buy:company:sort:", orgID, ".", buyCompanyID, ".", dayTime)
}

func deleteLogCache(id int64) {
	logData := getLogByID(id)
	Router2SystemConfig.MainCache.DeleteMark(getLogCacheMark(id))
	Router2SystemConfig.MainCache.DeleteSearchMark(getAnalysisCacheMark(""))
	if logData.OrgID > 0 {
		Router2SystemConfig.MainCache.DeleteSearchMark(getAnalysisBuyCompanySortCacheMark(logData.OrgID, ""))
	}
	if logData.BuyCompanyID > 0 {
		Router2SystemConfig.MainCache.DeleteMark(getAnalysisBuyCompanyDayCacheMark(logData.OrgID, logData.BuyCompanyID, CoreFilter.GetCarbonByTime(logData.CreateAt).Format("20060102")))
	}
}
