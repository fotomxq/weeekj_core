package AnalysisIndexValCustom

// ClearByCode 清理指定编码的数据
func ClearByCode(code string) (err error) {
	//清理数据
	err = indexValCustomDB.GetDelete().DeleteByField("code", code)
	if err != nil {
		return
	}
	//反馈
	return
}

// Clear 清理指标数据
func Clear() (err error) {
	//清理数据
	err = indexValCustomDB.GetDelete().Clear()
	if err != nil {
		return
	}
	//反馈
	return
}
