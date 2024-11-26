package AnalysisIndexEvent

// DeleteEventByCode 批量删除指标的所有事件
func DeleteEventByCode(code string) (err error) {
	err = eventDB.GetDelete().DeleteByField("code", code)
	if err != nil {
		return
	}
	return
}

// DeleteEventByID 删除指定事件
func DeleteEventByID(id int64) (err error) {
	err = eventDB.GetDelete().DeleteByID(id)
	if err != nil {
		return
	}
	return
}
