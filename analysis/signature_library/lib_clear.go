package AnalysisSignatureLibrary

import "fmt"

// ClearAllIndexData 清理所有指标的相关数据
func ClearAllIndexData(indexCode string) {
	_, _ = libDB.GetClient().DB.GetPostgresql().Exec(fmt.Sprintf("delete from %s where code1 = $1 or code2 = $1", libDB.GetClient().TableName), indexCode)
	return
}
