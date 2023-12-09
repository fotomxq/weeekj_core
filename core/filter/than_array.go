package CoreFilter

// CheckArray 比较两组[]int64结构
// 检查每一项，左侧最大。反馈总数量最大的
// 用于版本号检查，最多支持4位
// [1.0.0] -> [0.0.9]
// 左侧最大则反馈true，否则false
// 建议左侧为app实际安装，右侧为系统存储的版本，对比可得出最后是否需升级？
func CheckArray(a, b []int64) bool {
	offset := 0
	for {
		if offset < 4 {
			if len(a) > offset && len(b) > offset {
				if a[offset] > b[offset] {
					return true
				}
				if a[offset] < b[offset] {
					return false
				}
			}
		}
		if offset >= 4 {
			break
		}
		offset += 1
	}
	return false
}

// CheckArrayEq 检查两个数组是否相等
func CheckArrayEq(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	for _, v := range a {
		for _, v2 := range b {
			if v != v2 {
				return false
			}
		}
	}
	return true
}

// CheckArrayStringHave 检查两个数据是否存在交集
func CheckArrayStringHave(a, b []string) bool {
	for _, v := range a {
		for _, v2 := range b {
			if v == v2 {
				return true
			}
		}
	}
	for _, v := range b {
		for _, v2 := range a {
			if v == v2 {
				return true
			}
		}
	}
	return false
}

// CheckArrayStringLeftMustHaveRight 检查右侧数据是否完整被左侧包含
func CheckArrayStringLeftMustHaveRight(a, b []string) bool {
	isFind := false
	for _, v := range b {
		for _, v2 := range a {
			if v == v2 {
				isFind = true
				break
			}
		}
		if !isFind {
			return false
		}
	}
	return isFind
}

// MargeArrayString 合并两个数组并排重
func MargeArrayString(a, b []string) []string {
	for _, v := range b {
		f := false
		for _, v2 := range a {
			if v == v2 {
				f = true
				break
			}
		}
		if !f {
			a = append(a, v)
		}
	}
	return a
}

// MargeNoReplaceArrayInt64 去重后写入数组
func MargeNoReplaceArrayInt64(a []int64, b int64) []int64 {
	for _, v := range a {
		if v == b {
			return a
		}
	}
	a = append(a, b)
	return a
}
