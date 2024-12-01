package AnalysisSignatureLibrary

import "fmt"

// ClearAllIndexData 清理所有指标的相关数据
func ClearAllIndexData(libType string, indexCode string) {
	_, _ = libDB.GetClient().DB.GetPostgresql().Exec(fmt.Sprintf("delete from %s where lib_type = $1 AND (code1 = $2 or code2 = $2)", libDB.GetClient().TableName), libType, indexCode)
	return
}
