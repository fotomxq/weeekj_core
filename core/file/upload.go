package CoreFile

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

// SaveUploadFileToTemp 将上传文件存储到临时文件
func SaveUploadFileToTemp(c *gin.Context, targetDir string, formName string, maxSize int64, filterType []string) (DataGetUploadFileData, error) {
	//检查前置并生成数据结构
	res, dataByte, err := GetUploadFileData(c, formName, maxSize, filterType)
	if err != nil {
		return res, err
	}
	//构建新的文件名称
	res.NewName = strconv.FormatInt(res.CreateTime, 10) + "_" + res.SHA256 + "." + res.Type
	//确保上传文件夹存在
	if !IsExist(targetDir) {
		err = CreateFolder(targetDir)
		if err != nil {
			return res, err
		}
	}
	//构建文件路径
	res.Src = targetDir + Sep + res.NewName
	//文件不能已经存在
	if IsExist(res.Src) {
		return res, errors.New("upload file is exist")
	}
	//创建并保存文件
	err = WriteFile(res.Src, dataByte)
	if err != nil {
		return res, err
	}
	//返回
	return res, nil
}

// DataGetUploadFileData 上传文件结构体
type DataGetUploadFileData struct {
	//文件尺寸
	Size int64
	//文件名称，含类别
	Name string
	//文件名称，不含类别
	OnlyName string
	//新的文件名称
	NewName string
	//文件类别
	Type string
	//创建时间
	CreateTime int64
	//存储路径
	Src string
	//SHA256摘要
	SHA256 string
}

// GetUploadFileData 加载上传文件
// formName 表单名称
// maxSize 最大尺寸，字节类型
// filterType 过滤类型
func GetUploadFileData(c *gin.Context, formName string, maxSize int64, filterType []string) (DataGetUploadFileData, []byte, error) {
	//初始化
	var res DataGetUploadFileData
	//获取文件
	formFile, header, err := c.Request.FormFile(formName)
	if err != nil {
		err = errors.New(fmt.Sprint("get file by from file, ", err))
		return res, []byte{}, err
	}
	defer func() {
		if err := formFile.Close(); err != nil {
			return
		}
	}()
	//判断文件尺寸
	if header.Size > maxSize && maxSize > 0 {
		return res, []byte{}, errors.New("upload file size too lager")
	}
	res.Size = header.Size
	//获取文件名
	res.Name = header.Filename
	//过滤Names，不允许非英文等特殊字符
	res.Name = CoreFilter.CheckFilterStr(res.Name, 1, 250)
	//继续拆解文件名、类型
	names := strings.Split(res.Name, ".")
	res.Type = names[len(names)-1]
	res.OnlyName = names[0]
	for k, v := range names {
		if k == 0 {
			continue
		}
		if k >= len(names)-1 {
			break
		}
		res.OnlyName += "." + v
	}
	//甄别filterType
	if len(filterType) > 0 {
		isOK := false
		for _, v := range filterType {
			if v == res.Type {
				isOK = true
			}
		}
		if !isOK {
			return res, []byte{}, errors.New("upload file type is ban")
		}
	}
	//创建时间
	res.CreateTime = CoreFilter.GetNowTime().Unix()
	//获取字节数据
	buf := make([]byte, res.Size)
	n, err := formFile.Read(buf)
	if err != nil {
		return res, []byte{}, err
	}
	dataByte := buf[:n]
	//计算SHA256
	res.SHA256, err = CoreFilter.GetSha256Str(string(dataByte))
	if err != nil {
		return res, dataByte, errors.New("sha256 is error, " + err.Error())
	}
	if res.SHA256 == "" {
		return res, dataByte, errors.New("sha256 is error")
	}
	//完成后反馈数据
	return res, dataByte, nil
}
