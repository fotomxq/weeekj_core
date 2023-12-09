package CoreFilter

import (
	"fmt"
	"regexp"
	"strings"
)

//过滤器模块
// 提供标准，多方位的过滤和转化方案

// FilterStrForce 强过滤字符串
// param str string 要过滤的字符串
// return string 过滤后的字符串
func FilterStrForce(str string) (newStr string) {
	newStr = strings.Replace(str, "~", "～", -1)
	newStr = strings.Replace(newStr, "*", "＊", -1)
	newStr = strings.Replace(newStr, "<", "〈", -1)
	newStr = strings.Replace(newStr, ">", "〉", -1)
	newStr = strings.Replace(newStr, "$", "￥", -1)
	newStr = strings.Replace(newStr, "!", "！", -1)
	newStr = strings.Replace(newStr, "[", "【", -1)
	newStr = strings.Replace(newStr, "]", "】", -1)
	newStr = strings.Replace(newStr, "{", "｛", -1)
	newStr = strings.Replace(newStr, "}", "｝", -1)
	newStr = strings.Replace(newStr, "/", "／", -1)
	newStr = strings.Replace(newStr, "\\", "﹨", -1)
	return
}

// FilterPage 处理page
// param postPage string 用户提交的page
// return int 过滤后的页数
func FilterPage(postPage string) int64 {
	res, err := GetInt64ByString(postPage)
	if err != nil {
		res = 1
	}
	return FilterPageInt(res)
}
func FilterPageInt(postPage int64) int64 {
	if postPage < 1 {
		postPage = 1
	}
	return postPage
}

// FilterMax 处理max
// 限制最小值为1，最大值为999
// param postMax string 用户提交的max
// return int 过滤后的页数
func FilterMax(postMax string) int64 {
	res, err := GetInt64ByString(postMax)
	if err != nil {
		res = 1
	}
	return FilterMaxInt(res)
}
func FilterMaxInt(postMax int64) int64 {
	if postMax < 1 {
		postMax = 1
	}
	if postMax > 999 {
		postMax = 999
	}
	return postMax
}

// MatchStr
// param mStr string 验证
// param str string 要验证的字符串
// return bool 是否成功
func MatchStr(mStr string, str string) bool {
	res, err := regexp.MatchString(mStr, str)
	if err != nil {
		return false
	}
	return res
}

// SubStr 截取字符串
// param str string 要截取的字符串
// param star int 开始位置
// param length int 长度
// return string 新字符串
func SubStr(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0
	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length
	if start > end {
		start, end = end, start
	}
	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}

// SubStrQuick 快速截取字符串并输出...
func SubStrQuick(str string, limit int) string {
	if str == "" {
		return str
	}
	if len(str) > limit {
		newStr := SubStr(str, 0, limit)
		return fmt.Sprint(newStr, "...")
	}
	return str
}
