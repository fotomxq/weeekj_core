package ServiceOrder

import (
	"fmt"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// 获取订单缓冲标识码
func getOrderCacheMark(id int64) string {
	return fmt.Sprint("service:order:id:", id)
}

func GetOrderURLListCacheMark(id int64) string {
	return fmt.Sprint("service:order:url:list:", id)
}

func GetOrderURLDataCacheMark(id int64) string {
	return fmt.Sprint("service:order:url:data:", id)
}

// 清理缓冲
func deleteOrderCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getOrderCacheMark(id))
	Router2SystemConfig.MainCache.DeleteMark(GetOrderURLListCacheMark(id))
	Router2SystemConfig.MainCache.DeleteMark(GetOrderURLDataCacheMark(id))
}
