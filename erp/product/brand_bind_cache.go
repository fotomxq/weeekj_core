package ERPProduct

import (
	"fmt"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// 缓冲
func getBrandBindCacheMark(orgID, brandID, companyID, productID int64) string {
	return fmt.Sprint("erp:product:brand:bind:org.", orgID, ".brand.", brandID, ".company.", companyID, ".product.", productID)
}

func deleteBrandBindCache(orgID, brandID, companyID, productID int64) {
	Router2SystemConfig.MainCache.DeleteMark(getBrandBindCacheMark(orgID, brandID, companyID, productID))
}
