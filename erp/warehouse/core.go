package ERPWarehouse

import (
	AnalysisAny2 "gitee.com/weeekj/weeekj_core/v5/analysis/any2"
	CoreSQL2 "gitee.com/weeekj/weeekj_core/v5/core/sql2"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"sync"
)

//仓储模块
/**
1. 划分不同的仓库
2. 为仓库划分不同的区域
3. 存储不同的产品
*/
/**
关于批次设计
批次设计可以替代直接修改库存的设计。
1. 使用批次对应方法，进行出入库操作
2. 批次会影响产品的总库存记录，原先的store作为库存台账数据
*/
// TODO: 完成批次设计
// TODO：完成库存货位设计，该设计不影响其他任何内容，后续支持即可
// TODO: 核对库存日志，增设新的动作支持
// TODO: 全部完成后核对后续内容：
/**
1. 注意批次后续可以开关，根据库存设置启用，建议配置放入未来计划新增的BaseConfig2模块实现
2. 批次不应该存在货位设计，因为本身无法做到精细化管控
*/

var (
	////库存移动锁定
	//moveProductLock map[int64]sync.Mutex
	////库存锁定表变动锁定
	//moveProductListLock sync.Mutex
	//moveProductLock 仓储锁定
	moveProductLock sync.Mutex
	//OpenAnalysis 是否启动通用框架体系的统计支持
	OpenAnalysis = false
	//OpenSub 订阅
	OpenSub = false
	// SQL
	batchSQL    CoreSQL2.Client
	batchOutSQL CoreSQL2.Client
	storeSQL    CoreSQL2.Client
	//批次处理锁定
	batchWriteLock sync.Mutex
)

func Init() {
	//初始化表
	batchSQL.Init(&Router2SystemConfig.MainSQL, "erp_warehouse_batch")
	storeSQL.Init(&Router2SystemConfig.MainSQL, "erp_warehouse_store")
	batchOutSQL.Init(&Router2SystemConfig.MainSQL, "erp_warehouse_batch_out")
	//统计
	if OpenAnalysis {
		AnalysisAny2.SetConfigBeforeNoErr("erp_warehouse_store_product_count", 1, 365)
		AnalysisAny2.SetConfigBeforeNoErr("erp_warehouse_store_product_price", 1, 365)
		AnalysisAny2.SetConfigBeforeNoErr("erp_warehouse_store_count", 1, 365)
		AnalysisAny2.SetConfigBeforeNoErr("erp_warehouse_store_product_lack_count", 1, 365)
	}
	//nats
	if OpenSub {
		subNats()
	}
}
