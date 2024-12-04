package AnalysisIndexRFM

// ClearByCode 清理指定指标的数据
func ClearByCode(code string) (err error) {
	err = rfmDB.GetDelete().DeleteByField("code", code)
	return
}
