package FinanceReturnedMoney

import (
	"fmt"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

func getCompanyCacheMark(companyID int64) string {
	return fmt.Sprint("finance:return:money:company:company:id:", companyID)
}

func deleteCompanyCache(companyID int64) {
	Router2SystemConfig.MainCache.DeleteMark(getCompanyCacheMark(companyID))
	Router2SystemConfig.MainCache.DeleteSearchMark(getMargeAnalysisCacheMark(companyID))
	Router2SystemConfig.MainCache.DeleteMark(getMargeAnalysisCompanyCacheMark(companyID))
}
