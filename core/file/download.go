package CoreFile

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// DownloadFile 下载文件处理
// 给定一个文件序列组，该序列组是经过严格判定符合标准的，且不允许出现../的字符串结构
// param c *gin.Context
// param src string 文件路径
// param name string 文件名称
func DownloadFile(c *gin.Context, src string, name string) error {
	//检查src内是否包含..
	if strings.Count(src, "..") > 0 {
		return errors.New("src have \"...\", security mechanism blocking.")
	}
	//打开文件
	fd, err := LoadFile(src)
	if err != nil {
		return err
	}
	//添加头信息
	c.Header("Content-Type", "application/octet-stream")
	c.Header("content-disposition", "attachment; filename=\""+name+"\"")
	//写入文件
	_, err = c.Writer.Write(fd)
	if err != nil {
		return err
	}
	//返回成功
	return nil
}

// DownloadFileByByte 下载文件
// 注入byte方式
func DownloadFileByByte(c *gin.Context, data []byte, name string) error {
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", strconv.Itoa(len(data)))
	c.Header("Content-Disposition", "attachment; filename=\""+name+"\"")
	_, err := c.Writer.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// CreateDownloadLS 创建临时文件到指定目录
// 文件名称以文件内容的SHA1为主
// param dirSrc string 存放目录 eg : abc/
// param content []byte 文件内容
// param FileType string 文件类型
// return string,string,string,error 存储路径，相对路径，文件名称，错误
func CreateDownloadLS(dirSrc string, content []byte, FileType string) (string, string, string, error) {
	//计算SHA1
	hasher := sha1.New()
	_, err := hasher.Write(content)
	if err != nil {
		return "", "", "", err
	}
	sha := hasher.Sum(nil)
	shaStr := hex.EncodeToString(sha)
	//创建文件路径
	fileName := shaStr + "." + FileType
	nowTime := CoreFilter.GetNowTime().Format("2006010215")
	dsrc := "ls" + Sep + nowTime
	src := dirSrc + Sep + dsrc
	if err = CreateFolder(src); err != nil {
		return "", "", "", err
	}
	src += Sep + fileName
	//删除60分钟之前的数据
	sinceTime := CoreFilter.GetNowTime().Add(-time.Minute * 30).Format("2006010215")
	sinceTimeDSrc := dirSrc + Sep + sinceTime
	fileList, err := GetFileList(dirSrc, []string{}, false)
	if err == nil {
		for _, v := range fileList {
			if v == nowTime || v == sinceTime {
				continue
			}
			err = DeleteF(sinceTimeDSrc + Sep + "v")
			if err != nil {
				return src, dsrc, fileName, WriteFile(src, content)
			}
		}
	}
	return src, dsrc, fileName, WriteFile(src, content)
}

// DownloadByURLToTemp 直接下载文件存储临时文件中
func DownloadByURLToTemp(url string, path string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	// Create output file
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = out.Close()
	}()
	// copy stream
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	//反馈
	return nil
}

// DownloadByURLToByte 直接下载文件到二进制数据
func DownloadByURLToByte(url string) ([]byte, error) {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	data, err := io.ReadAll(resp.Body)
	return data, err
}
