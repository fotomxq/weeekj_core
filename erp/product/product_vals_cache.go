package ERPProduct

import (
	"fmt"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

func getProductValsCacheMark(orgID int64, productID int64) (val string) {
	return fmt.Sprint("erp:product:brand:bind:org.", orgID, ".product.", productID)
}

func deleteProductValsCache(orgID, productID int64) {
	Router2SystemConfig.MainCache.DeleteMark(getProductValsCacheMark(orgID, productID))
}
