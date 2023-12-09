package MallCoreMod

import "fmt"

func getProductCacheMark(id int64) string {
	return fmt.Sprint("mall:core:product:id:", id)
}
