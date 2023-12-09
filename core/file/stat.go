package CoreFile

import (
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"os"
)

//获取文件信息
//param src string 文件路径
//return os.FileInfo,error 文件信息，错误
func GetFileInfo(src string) (os.FileInfo, error) {
	c, err := os.Stat(src)
	return c, err
}

//GetFileSize 获取文件大小
//param src string 文件路径
//return int64,bool 文件大小，错误
func GetFileSize(src string) (int64, error) {
	info, err := os.Stat(src)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

//GetFileListCount 查询文件夹下文件个数
//param src string 文件夹路径
//return int,error 文件个数,错误
func GetFileListCount(src string) (int, error) {
	dir, err := ioutil.ReadDir(src)
	if err != nil {
		return 0, err
	}
	var res int
	for range dir {
		res += 1
	}
	return res, nil
}

//获取文件SHA1值
//param src string 文件路径
//return string,error SHA1值,错误
func GetFileSha1(src string) (string, error) {
	content, err := LoadFile(src)
	if err != nil {
		return "", err
	}
	if content != nil {
		sha := sha1.New()
		_, err = sha.Write(content)
		if err != nil {
			return "", err
		}
		res := sha.Sum(nil)
		return hex.EncodeToString(res), nil
	}
	return "", nil
}

//判断是否为文件夹
//param src string 文件夹路径
//return bool 是否为文件夹
func IsFolder(src string) bool {
	info, err := os.Stat(src)
	return err == nil && info.IsDir()
}

//判断是否为文件
//param src string 文件路径
//return bool 是否为文件
func IsFile(src string) bool {
	info, err := os.Stat(src)
	return err == nil && !info.IsDir()
}

//判断文件或文件夹是否存在
//param src string 文件路径
//return bool 是否存在
func IsExist(src string) bool {
	_, err := os.Stat(src)
	return err == nil && !os.IsNotExist(err)
}
