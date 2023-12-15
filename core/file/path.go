package CoreFile

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	"os"
	"runtime"
	"strings"
)

//路径处理

// GetTimeDirSrc 获取并创建时间序列创建的多级文件夹
// eg : Return and create the path ,"[src]/201611/"
// eg : Return and create the path ,"[src]/201611/2016110102-03[appendFileType]"
// param src string 文件路径
// param appendFileType string 是否末尾追加文件类型，如果指定值，则返回
// return string,error 新时间周期目录，错误
func GetTimeDirSrc(src string, appendFileType string) (string, error) {
	newSrc := src + Sep + CoreFilter.GetNowTime().Format("200601")
	err := CreateFolder(newSrc)
	if err != nil {
		return "", err
	}
	newSrc = newSrc + Sep
	if appendFileType != "" {
		newSrc = newSrc + CoreFilter.GetNowTime().Format("20060102-03") + appendFileType
	}
	return newSrc, nil
}

// GetDir 获取目录路径
// param path string 地址路径
// return string 返回值
func GetDir(path string) string {
	return SubString(path, 0, strings.LastIndex(path, "/"))
}

// SubString 截取字符串
// param str string 字符串
// param start int 开始位置
// param end int 结束位置
// return string 结果字符串
func SubString(str string, start, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < start || end > length {
		panic("end is wrong")
	}

	return string(rs[start:end])
}

// GetFileNames 获取文件名称分割序列
// param src string 文件路径
// return map[string]string,error 文件名称序列，错误 eg : {"name","abc.jpg","type":"jpg","only-name":"abc"}
func GetFileNames(src string) (map[string]string, error) {
	info, err := os.Stat(src)
	if err != nil {
		return nil, err
	}
	res := map[string]string{
		"name":      info.Name(),
		"type":      "",
		"only-name": info.Name(),
	}
	names := strings.Split(res["name"], ".")
	if len(names) < 2 {
		return res, nil
	}
	res["type"] = names[len(names)-1]
	res["only-name"] = names[0]
	for i := range names {
		if i != 0 && i < len(names)-1 {
			res["only-name"] = res["only-name"] + "." + names[i]
		}
	}
	return res, nil
}

// GetNowFileName 获取当前运行程序的名称
func GetNowFileName() (string, error) {
	_, fileSrc, _, _ := runtime.Caller(0)
	fileNames, err := GetFileNames(fileSrc)
	if err != nil {
		return "", err
	}
	return fileNames["only-name"], nil
}
