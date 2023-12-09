package CoreFilter

import "fmt"

// DerefString 将引用字符串改为非引用关系
func DerefString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

// DerefInt64 将引用int64改为非引用关系
func DerefInt64(s *int64) int64 {
	if s != nil {
		return *s
	}
	return 0
}

// CutStringAndEncrypt 将文本进行加密处理
// 自动保留前N位和后N位
// startN 开头保留前几位
// endStartN 结尾保留后几位
func CutStringAndEncrypt(str string, startN int, endN int) string {
	strLen := len(str)
	if strLen == 0 {
		return ""
	}
	if startN >= strLen {
		startN = strLen - 1
	}
	if strLen < startN+endN {
		return str
	}
	endStartN := strLen - endN
	if endStartN >= strLen {
		endStartN = strLen - 1
	}
	newStr := fmt.Sprint(str[0:startN], "***", str[endStartN:strLen])
	return newStr
}
