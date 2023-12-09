package CoreFilter

import "strings"

//GetURLNameType 分解URL获取名称和类型
//param sendURL URL地址
//return map[string]string 返回值集合
func GetURLNameType(sendURL string) map[string]string {
	res := map[string]string{
		"full-name": "",
		"only-name": "",
		"type":      "",
	}
	urls := strings.Split(sendURL, "/")
	if len(urls) < 1 {
		return res
	}
	res["full-name"] = urls[len(urls)-1]
	if res["full-name"] == "" {
		res["only-name"] = res["full-name"]
		return res
	}
	names := strings.Split(res["full-name"], ".")
	if len(names) < 2 {
		return res
	}
	res["type"] = names[len(names)-1]
	for i := 0; i <= len(names); i++ {
		if i == len(names)-1 {
			break
		}
		if res["only-name"] == "" {
			res["only-name"] = names[i]
		} else {
			res["only-name"] += "." + names[i]
		}
	}
	return res
}
