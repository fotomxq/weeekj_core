package BaseSQLTools

import (
	"fmt"
)

// QuickAnalysis 获取统计
type QuickAnalysis struct {
	//Quick
	quickClient *Quick
}

// GetCountByField 获取指定条件的数据量
func (c *QuickAnalysis) GetCountByField(field string, val any) (count int64) {
	//检查ID
	if field == "" {
		return
	}
	//获取数据
	count = c.quickClient.client.Analysis().Count(fmt.Sprintf("%s = $1", field), val)
	//反馈
	return
}
