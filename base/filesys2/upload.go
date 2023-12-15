package BaseFileSys2

import (
	"errors"
	"fmt"
	CoreFile "github.com/fotomxq/weeekj_core/v5/core/file"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

// DataUploadFileType 上传文件结构体
type DataUploadFileType struct {
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

// ArgsUploadFile 上传文件参数
type ArgsUploadFile struct {
	//文件路径
	FileSrc string `json:"fileSrc"`
	//目标路径
	TargetSrc string `json:"targetSrc"`
	//文件最大尺寸
	MaxSize int64 `json:"maxSize"`
	//限制格式
	FilterType []string `json:"filterType"`
	//创建用户
	UserID int64 `json:"userID" check:"id" empty:"true"`
	//创建组织
	OrgID int64 `json:"orgID" check:"id" empty:"true"`
	//IP地址
	IP string `json:"ip"`
	//过期时间
	ExpireAt time.Time `json:"expireAt" check:"defaultTime" empty:"true"`
	//是否为公开的文件
	IsPublic bool `json:"isPublic" check:"bool" empty:"true"`
	//描述
	Des string `json:"des" check:"des" min:"1" max:"600" empty:"true"`
}

// UploadFile 上传文件
// 上传到本地，然后根据存储规格区别处理
func UploadFile(c *gin.Context, args *ArgsUploadFile) (newClaimID int64, errCode string, err error) {
	//将文件存储到临时文件下
	//生成子目录结构
	//nowTime := CoreFilter.GetNowTime()
	//newDir := localDefaultDir + CoreFile.Sep + nowTime.Format("200601"+CoreFile.Sep+"02"+CoreFile.Sep+"15") + CoreFile.Sep
	//保存文件
	//var newFileData DataUploadFileType
	//newFileData, err = uploadFileToLocal(c, &argsUploadFileToLocal{
	//	FormName:   "file",
	//	TargetSrc:  newDir,
	//	MaxSize:    args.MaxSize,
	//	FilterType: args.FilterType,
	//})
	//if err != nil {
	//	errCode = "文件存储失败"
	//	err = errors.New("cannot save upload file, " + err.Error())
	//	return
	//}
	//保存到数据库结构体
	//var newClaimData FieldsFileClaim
	//var newCoreData FieldsFile
	//newClaimData, newCoreData, errCode, err = createClaim(&argsCreateClaim{
	//	CreateIP:   c.ClientIP(),
	//	CreateInfo: args.CreateInfo,
	//	UserID:     args.UserInfo.Info.ID,
	//	OrgID:      args.UserInfo.Info.OrgID,
	//	FileSize:   newFileData.Size,
	//	FileType:   newFileData.Type,
	//	FileHash:   newFileData.SHA256,
	//	FileSrc:    newFileData.Src,
	//	ExpireAt:   args.ExpireAt,
	//	FromInfo:   args.FromInfo,
	//	Infos:      args.Infos,
	//	ClaimInfos: args.ClaimInfos,
	//	Des:        args.Des,
	//})
	//// 读取文件
	//src, err := file.Open()
	//if err != nil {
	//	return "", err
	//}
	//defer src.Close()
	//
	//// 创建文件
	//dst, err := os.Create(path)
	//if err != nil {
	//	return "", err
	//}
	//defer dst.Close()
	//
	//// 复制文件
	//if _, err := io.Copy(dst, src); err != nil {
	//	return "", err
	//}
	//return path, nil
	return
}

// argsUploadFileToLocal 上传文件参数
type argsUploadFileToLocal struct {
	//表单名称
	FormName string
	//目标路径，末尾必须添加Sep
	TargetSrc string
	//文件最大大小，如果为0则不限制
	MaxSize int64
	//文件类别限制
	FilterType []string
}

// uploadFileToLocal 上传文件参数
// 可利用该方法，实现任意文件、目标得上传
// 注意，目标路径不要交给用户设计，否则将出现重大系统级漏洞，影响系统安全
// 存储后，默认根据创建"Unix时间戳_"结构设计文件名称
// param c *gin.Context
// return FileUploadType 文件类型
// return error 错误信息
func uploadFileToLocal(c *gin.Context, args *argsUploadFileToLocal) (DataUploadFileType, error) {
	//检查前置并生成数据结构
	res, dataByte, err := checkUpload(c, args.FormName, args.MaxSize, args.FilterType)
	if err != nil {
		return res, err
	}
	//构建新的文件名称
	res.NewName = strconv.FormatInt(res.CreateTime, 10) + "_" + res.SHA256 + "." + res.Type
	//确保上传文件夹存在
	if !CoreFile.IsExist(args.TargetSrc) {
		err = CoreFile.CreateFolder(args.TargetSrc)
		if err != nil {
			return res, err
		}
	}
	//构建文件路径
	res.Src = args.TargetSrc + CoreFile.Sep + res.NewName
	//文件不能已经存在
	if CoreFile.IsExist(res.Src) {
		return res, errors.New("upload file is exist")
	}
	//创建并保存文件
	err = CoreFile.WriteFile(res.Src, dataByte)
	if err != nil {
		return res, err
	}
	//返回
	return res, nil
}

// checkUpload 上传文件前置判断
func checkUpload(c *gin.Context, formName string, maxSize int64, filterType []string) (DataUploadFileType, []byte, error) {
	//初始化
	var res DataUploadFileType
	//获取文件
	formFile, header, err := c.Request.FormFile(formName)
	if err != nil {
		err = errors.New(fmt.Sprint("get file by from file, ", err))
		return res, []byte{}, err
	}
	defer func() {
		if err := formFile.Close(); err != nil {
			CoreLog.Error(err)
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
