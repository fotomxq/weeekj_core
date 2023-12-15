package ERPProduct

import (
	"fmt"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// 缓冲
func getProductCacheMark(id int64) string {
	return fmt.Sprint("erp:product:id:", id)
}

func getProductCodeCacheMark(orgID int64, code string) string {
	return fmt.Sprint("erp:product:code:", orgID, ".", code)
}

func deleteProductCache(id int64) {
	data := getProductByID(id)
	Router2SystemConfig.MainCache.DeleteMark(getProductCacheMark(id))
	if data.ID > 0 {
		Router2SystemConfig.MainCache.DeleteMark(getProductCodeCacheMark(data.OrgID, data.Code))
	}
}
