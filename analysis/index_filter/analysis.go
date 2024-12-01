package AnalysisIndexFilter

// GetCount 获取指标的数据量
func GetCount(code string) (count int64) {
	count = filterDB.GetAnalysis().GetCountByField("code", code)
	return
}
