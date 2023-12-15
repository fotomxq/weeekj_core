package ERPWarehouse

import (
	AnalysisAny2 "github.com/fotomxq/weeekj_core/v5/analysis/any2"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// updateAnalysis 更新基本数据
func updateAnalysis(orgID int64) {
	//统计产品品种个数
	var productSortCount int64
	_ = Router2SystemConfig.MainDB.Get(&productSortCount, "SELECT COUNT(product_id) FROM erp_warehouse_store WHERE org_id = $1 AND count > 0 GROUP BY product_id")
	AnalysisAny2.AppendData("re", "erp_warehouse_store_product_count", CoreFilter.GetNowTime(), orgID, 0, 0, 0, 0, productSortCount)
	//缺货品种个数
	var productSortLackCount int64
	_ = Router2SystemConfig.MainDB.Get(&productSortLackCount, "SELECT COUNT(product_id) FROM erp_warehouse_store WHERE org_id = $1 AND count < 1 GROUP BY product_id")
	AnalysisAny2.AppendData("re", "erp_warehouse_store_product_lack_count", CoreFilter.GetNowTime(), orgID, 0, 0, 0, 0, productSortLackCount)
}
