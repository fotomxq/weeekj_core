package AnalysisIndexFilter

// Clear 清理掉数据
func Clear(code string) (err error) {
	err = filterDB.GetDelete().DeleteByFieldReal("code", code)
	if err != nil {
		return
	}
	return
}
