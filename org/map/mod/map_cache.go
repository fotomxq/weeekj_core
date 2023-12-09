package OrgMapMod

import "fmt"

// 获取缓冲名称
func getMapCacheMark(id int64) string {
	return fmt.Sprint("org:map:id:", id)
}
