package CoreFilter

//GetMapKey 从一组map中找到对应的键位，如果不存在则反馈空字符串
func GetMapKey(key int64, mapData map[int64]string) string {
	for k, v := range mapData {
		if k == key {
			return v
		}
	}
	return ""
}

func GetMapKeys(keys []int64, mapData map[int64]string) []string {
	var result []string
	for k, v := range mapData {
		for _, v2 := range keys {
			if k == v2 {
				result = append(result, v)
			}
		}
	}
	return result
}
