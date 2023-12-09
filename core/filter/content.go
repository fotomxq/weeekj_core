package CoreFilter

import (
	"regexp"
	"strings"
)

// GetFileIDListByContent 检索content文本内的所有file数据
// 识别方案：({\$file:)[a-zA-Z0-9]*(})
func GetFileIDListByContent(str string) ([]string, error) {
	//搜索字符串
	re, err := regexp.Compile("({file:)[a-zA-Z0-9]*(})")
	if err != nil {
		return []string{}, err
	}
	res := re.FindAllString(str, -1)
	return res, nil
}

func GetIDsInString(str string, split string) []int64 {
	l := strings.Split(str, split)
	var r []int64
	for k := 0; k < len(l); k++ {
		v := l[k]
		r = append(r, GetInt64ByStringNoErr(v))
	}
	return r
}
