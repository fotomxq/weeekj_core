package ServiceCompany

import (
	"fmt"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// 获取缓冲名称
func getCompanyCacheMark(id int64) string {
	return fmt.Sprint("service:company:id:", id)
}

func getCompanyCacheHashMark(orgID int64, hash string) string {
	return fmt.Sprint("service:company:hash:", orgID, ".", hash)
}

// 删除公司缓冲
func deleteCompanyCache(id int64) {
	data := getCompany(id)
	Router2SystemConfig.MainCache.DeleteMark(getCompanyCacheMark(id))
	if data.ID > 0 {
		Router2SystemConfig.MainCache.DeleteMark(getCompanyCacheHashMark(data.OrgID, data.Hash))
	}
}
