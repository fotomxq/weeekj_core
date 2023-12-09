package ServiceInfoExchangeMod

import "fmt"

// 获取信息交互缓冲名称
func getInfoCacheMark(id int64) string {
	return fmt.Sprint("service:info:exchange:id:", id)
}
