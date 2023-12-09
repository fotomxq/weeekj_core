package BlogCoreMod

import "fmt"

// 获取文章缓冲标识码
func getContentCacheMark(id int64) string {
	return fmt.Sprint("blog:core:id:", id)
}
