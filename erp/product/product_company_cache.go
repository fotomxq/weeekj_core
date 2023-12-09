package ERPProduct

import (
	"fmt"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

func getProductCompanyCacheMark(id int64) string {
	return fmt.Sprint("erp:product:company:id:", id)
}

func deleteProductCompanyCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getProductCompanyCacheMark(id))
}
