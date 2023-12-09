package CoreFilter

// EqID2 对比两个ID是否相同，第一个作为参数可以小于1
func EqID2(argID int64, dataID int64) (b bool) {
	if argID > -1 && argID != dataID {
		return
	}
	return true
}

// EqHaveID2 是否包含ID
func EqHaveID2(argID int64, arr []int64) (b bool) {
	for _, v := range arr {
		if v == argID {
			return true
		}
	}
	return
}
