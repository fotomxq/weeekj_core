package EAMWarehouse

import (
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

var (
	//缓冲时间
	cacheWarehouseTime    = 1800
	cacheWarehouseLogTime = 1800
	//数据表
	warehouseDB    CoreSQL2.Client
	warehouseLogDB CoreSQL2.Client
)

// Init 初始化
func Init() {
	//初始化数据表
	warehouseDB.Init(&Router2SystemConfig.MainSQL, "eam_warehouse")
	warehouseLogDB.Init(&Router2SystemConfig.MainSQL, "eam_warehouse_log")
}
