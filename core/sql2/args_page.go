package CoreSQL2

import "fmt"

// ArgsPages 分页结构体
type ArgsPages struct {
	Page int64  `json:"page" check:"page"`
	Max  int64  `json:"max" check:"max"`
	Sort string `json:"sort" check:"sort"`
	Desc bool   `json:"desc" check:"desc"`
}

// GetCacheMark 获取缓冲名称
func (t *ArgsPages) GetCacheMark() string {
	return fmt.Sprint(t.Page, ".", t.Max, ".", t.Sort, ".", t.Desc)
}
