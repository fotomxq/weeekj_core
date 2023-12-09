package ServiceOrderMod

import "fmt"

// 获取订单缓冲标识码
func getOrderCacheMark(id int64) string {
	return fmt.Sprint("service:order:id:", id)
}
