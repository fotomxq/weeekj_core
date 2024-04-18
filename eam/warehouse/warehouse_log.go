package EAMWarehouse

import (
	"fmt"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetWarehouseLog 获取库存批次列表参数
type ArgsGetWarehouseLog struct {
}

// GetWarehouseLog 获取库存批次列表
func GetWarehouseLog(args *ArgsGetWarehouseLog) (dataList []FieldsWarehouseLog, dataCount int64, err error) {
	return
}

// ArgsInWarehouseLog 批次入库参数
type ArgsInWarehouseLog struct {
}

// InWarehouseLog 批次入库
func InWarehouseLog(args *ArgsInWarehouseLog) (err error) {
	return
}

// ArgsMoveWarehouseLog 批次移库参数
type ArgsMoveWarehouseLog struct {
}

// MoveWarehouseLog 批次移库
func MoveWarehouseLog(args *ArgsMoveWarehouseLog) (err error) {
	return
}

// ArgsOutWarehouseLog 批次出库参数
type ArgsOutWarehouseLog struct {
}

// OutWarehouseLog 批次出库
func OutWarehouseLog(args *ArgsOutWarehouseLog) (err error) {
	return
}

// getWarehouseLogData 获取模板数据
func getWarehouseLogData(id int64) (data FieldsWarehouseLog) {
	cacheMark := getWarehouseLogCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := warehouseDB.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "delete_at", "product_id", "order_id", "count", "total", "price", "warranty_at"}).GetByID(id).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheWarehouseLogTime)
	return
}

// 缓冲
func getWarehouseLogCacheMark(id int64) string {
	return fmt.Sprint("eam:warehouse:log:id.", id)
}

func deleteWarehouseLogCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getWarehouseLogCacheMark(id))
}
