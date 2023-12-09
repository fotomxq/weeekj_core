package CoreFile

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

//文件操作底层封装

var (
	//BaseSrc 根路径
	BaseSrc = ""
	//Sep 分隔符
	Sep = string(os.PathSeparator)
)

//BaseDir 获取文件根路径
func BaseDir() string{
	fileSrc := os.Args[0]
	BaseSrc = filepath.Dir(fileSrc)
	return BaseSrc
}

//BaseWDDir 获取命令行路径
func BaseWDDir() (string, error){
	return os.Getwd()
}

//MoveF 移动文件或文件夹
//param src string 文件路径
//param dest string 新路径
//return error
func MoveF(src string, dest string) error {
	return os.Rename(src, dest)
}

//DeleteF 删除文件或文件夹
//param src string 文件路径
//return error
func DeleteF(src string) error {
	return os.RemoveAll(src)
}

//LoadFile 读取文件
//param src string 文件路径
//return []byte,error 文件数据,错误
func LoadFile(src string) ([]byte, error) {
	return ioutil.ReadFile(src)
}

//WriteFile 写入文件
//param src string 文件路径
//param content []byte 写入内容
//return error
func WriteFile(src string, content []byte) error {
	return ioutil.WriteFile(src, content, 0666)
}

//WriteFileAppend 追加写入文件
//param src string 文件路径
//param content []byte 写入内容
//return error
func WriteFileAppend(src string, content []byte) error {
	f, err := os.OpenFile(src, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(content)
	if err != nil {
		return err
	}
	return nil
}

//CopyFile 复制文件
//param src string 文件路径
//param dest string 新路径
//return error
func CopyFile(src string, dest string) error {
	srcF, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcF.Close()
	destF, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destF.Close()
	_, err = io.Copy(destF, srcF)
	if err != nil {
		return err
	}
	return nil
}

//CreateFolder 创建多级文件夹
//param src string 新文件夹路径
//return error
func CreateFolder(src string) error {
	return os.MkdirAll(src, os.ModePerm)
}

//CopyFolder 复制文件夹
// 自递归复制文件夹
//param src string 源路径
//param dest string 目标路径
//return bool 是否成功
func CopyFolder(src string, dest string) bool {
	dir, err := ioutil.ReadDir(src)
	if err != nil {
		return false
	}
	if CreateFolder(dest) != nil {
		return false
	}
	for _, v := range dir {
		vSrc := src + Sep + v.Name()
		vDest := dest + Sep + v.Name()
		if v.IsDir()  {
			if CreateFolder(vDest) != nil {
				return false
			}
			if !CopyFolder(vSrc, vDest) {
				return false
			}
		} else {
			if CopyFile(vSrc, vDest) != nil {
				return false
			}
		}
	}
	return true
}

//GetFileList 获取文件列表
// 按照文件名，倒叙排列返回
//param src string 查询的文件夹路径,eg: /var/data
//param filters []string 仅保留的文件，文件夹除外
//param isSrc bool 返回是否为文件路径
//return []string,error 文件列表,错误
func GetFileList(src string, filters []string, isSrc bool) ([]string, error) {
	//初始化
	var fs []string
	//读取目录
	dir, err := ioutil.ReadDir(src)
	if err != nil {
		return nil, err
	}
	//遍历目录文件
	for _, v := range dir {
		var appendSrc string
		if isSrc  {
			appendSrc = src + Sep + v.Name()
		} else {
			appendSrc = v.Name()
		}
		if v.IsDir()  || len(filters) < 1 {
			fs = append(fs, appendSrc)
			continue
		}
		names := strings.Split(v.Name(), ".")
		if len(names) == 1 {
			fs = append(fs, appendSrc)
			continue
		}
		t := names[len(names)-1]
		for _, filterValue := range filters {
			if t != filterValue {
				continue
			}
			fs = append(fs, appendSrc)
		}
	}
	//对数组进行倒叙排序
	sort.Sort(sort.Reverse(sort.StringSlice(fs)))
	//返回
	return fs, nil
}
