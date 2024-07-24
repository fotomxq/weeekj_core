package CoreFilter

// CheckInt64InArray 检查int64是否在列内
func CheckInt64InArray(arr []int64, c int64) bool {
	for _, v := range arr {
		if v == c {
			return true
		}
	}
	return false
}

// CheckStringInArray 检查string是否在列内
func CheckStringInArray(arr []string, c string) bool {
	for _, v := range arr {
		if v == c {
			return true
		}
	}
	return false
}
