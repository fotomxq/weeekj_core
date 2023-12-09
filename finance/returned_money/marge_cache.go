package FinanceReturnedMoney

import (
	"fmt"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// 获取缓冲
func getMargeCacheMark(id int64) string {
	return fmt.Sprint("finance:returned:money:marge:id:", id)
}

func getMargeAnalysisCacheMark(companyID int64) string {
	return fmt.Sprint("finance:returned:money:marge:analysis:", companyID)
}

func getMargeAnalysisCompanyCacheMark(companyID int64) string {
	return fmt.Sprint("finance:returned:money:marge:analysis:company:", companyID)
}

func getMargeAnalysisDayCacheMark(companyID int64, dateStr string) string {
	return fmt.Sprint("finance:returned:money:marge:analysis:", companyID, ".", dateStr)
}

func deleteMargeCache(id int64) {
	data := getMargeByID(id)
	Router2SystemConfig.MainCache.DeleteMark(getMargeCacheMark(id))
	if data.ID > 0 {
		Router2SystemConfig.MainCache.DeleteSearchMark(getMargeAnalysisCacheMark(data.CompanyID))
		Router2SystemConfig.MainCache.DeleteMark(getMargeAnalysisCompanyCacheMark(data.CompanyID))
	}
}
