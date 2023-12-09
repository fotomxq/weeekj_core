package BaseFileUpload

import (
	"encoding/base64"
	"errors"
	"fmt"
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	BaseFileSys "gitee.com/weeekj/weeekj_core/v5/base/filesys"
	BaseQiniu "gitee.com/weeekj/weeekj_core/v5/base/qiniu"
	CoreFile "gitee.com/weeekj/weeekj_core/v5/core/file"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	UserCore "gitee.com/weeekj/weeekj_core/v5/user/core"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

//上传文件模块

// 上传文件类别封装
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

var (
	//文件系统的本地存储路径
	// 默认程序运行位置
	// 会自动在下面创建子目录分级处理
	localDefaultDir = "./files"
	//文件系统的文件最大上限
	// 默认60MB
	localFileMaxSize int64 = 60 * 1024 * 1024
	//文件系统默认上传限制
	localFileFilterType []string
)

// SetLocalConfig 设置本地存储默认设定
func SetLocalConfig(tLocalDefaultDir string, tLocalFileMaxSize int64, tLocalFileFilterType []string) {
	localDefaultDir = tLocalDefaultDir
	localFileMaxSize = tLocalFileMaxSize
	localFileFilterType = tLocalFileFilterType
}

// ArgsUpdateMarge 融合上传文件参数
type ArgsUpdateMarge struct {
	//表单名称
	FormName string `json:"formName"`
	//创建来源
	CreateInfo CoreSQLFrom.FieldsFrom `json:"createInfo"`
	//描述
	Des string `json:"des" check:"des" min:"1" max:"600" empty:"true"`
	//IP地址
	IP string `json:"ip"`
	//创建用户
	UserID int64 `json:"userID" check:"id" empty:"true"`
	//创建组织
	OrgID int64 `json:"orgID" check:"id" empty:"true"`
	//是否为公开的文件
	IsPublic bool `json:"isPublic" check:"bool" empty:"true"`
	//过期时间
	ExpireAt time.Time `json:"expireAt" check:"defaultTime" empty:"true"`
	//扩展信息
	Infos []CoreSQLConfig.FieldsConfigType `json:"infos"`
	//扩展信息
	ClaimInfos []CoreSQLConfig.FieldsConfigType `json:"claimInfos"`
}

// UpdateMarge 融合上传文件
func UpdateMarge(c *gin.Context, args *ArgsUpdateMarge) (newFileClaimID int64, errCode string, err error) {
	if args.Infos == nil {
		args.Infos = []CoreSQLConfig.FieldsConfigType{}
	}
	if args.ClaimInfos == nil {
		args.ClaimInfos = []CoreSQLConfig.FieldsConfigType{}
	}
	fileSaveFrom := BaseConfig.GetDataStringNoErr("FileSaveFrom")
	if fileSaveFrom == "" {
		fileSaveFrom = "local"
	}
	var fileClaimData BaseFileSys.FieldsFileClaimType
	switch fileSaveFrom {
	case "local":
		//本地文件
		fileClaimData, errCode, err = UploadFileToFileSys(c, &ArgsUploadFileToFileSys{
			FormName:   args.FormName,
			TargetSrc:  localDefaultDir,
			MaxSize:    localFileMaxSize,
			FilterType: localFileFilterType,
			CreateInfo: args.CreateInfo,
			ExpireAt:   args.ExpireAt,
			FromInfo:   CoreSQLFrom.FieldsFrom{System: "local"},
			Infos:      args.Infos,
			ClaimInfos: args.ClaimInfos,
			Des:        args.Des,
		})
		if err != nil {
			return
		}
		newFileClaimID = fileClaimData.ID
		return
	case "qiniu":
		//七牛云
		fileQiniuDefaultBucketName := BaseConfig.GetDataStringNoErr("FileQiniuDefaultBucketName")
		fileClaimData, errCode, err = UploadFileToFileSysByQiniu(c, &ArgsUploadFileToFileSysByQiniu{
			FormName:   args.FormName,
			BucketName: fileQiniuDefaultBucketName,
			IP:         args.IP,
			CreateInfo: args.CreateInfo,
			UserID:     args.UserID,
			OrgID:      args.OrgID,
			IsPublic:   args.IsPublic,
			ExpireAt:   args.ExpireAt,
			ClaimInfos: args.ClaimInfos,
			Des:        args.Des,
		})
		if err != nil {
			return
		}
		newFileClaimID = fileClaimData.ID
		return
	default:
		//未知类型
		errCode = "err_file_no_support"
		return
	}
}

// ArgsUploadFileToFileSysByLocal 使用本地默认设置上传新的文件参数
type ArgsUploadFileToFileSysByLocal struct {
	//表单名称
	FormName string
	//创建来源
	CreateInfo CoreSQLFrom.FieldsFrom
	//过期时间
	ExpireAt time.Time
	//扩展信息
	Infos []CoreSQLConfig.FieldsConfigType
	//描述
	Des string
}

// UploadFileToFileSysByLocal 使用本地默认设置上传新的文件
func UploadFileToFileSysByLocal(c *gin.Context, args *ArgsUploadFileToFileSysByLocal) (data BaseFileSys.FieldsFileClaimType, errCode string, err error) {
	return UploadFileToFileSys(c, &ArgsUploadFileToFileSys{
		FormName:   args.FormName,
		TargetSrc:  localDefaultDir,
		MaxSize:    localFileMaxSize,
		FilterType: localFileFilterType,
		CreateInfo: args.CreateInfo,
		ExpireAt:   args.ExpireAt,
		FromInfo:   CoreSQLFrom.FieldsFrom{System: "local"},
		Infos:      args.Infos,
		ClaimInfos: args.Infos,
		Des:        args.Des,
	})
}

// ArgsUploadFileToFileSysByQiniu 上传到七牛云参数
type ArgsUploadFileToFileSysByQiniu struct {
	//表单名称
	FormName string
	//存储块名称
	BucketName string
	//IP地址
	IP string
	//创建来源
	CreateInfo CoreSQLFrom.FieldsFrom
	//创建用户
	UserID int64
	//创建组织
	OrgID int64
	//是否为公开的文件
	IsPublic bool `json:"isPublic" check:"bool" empty:"true"`
	//过期时间
	ExpireAt time.Time
	//扩展信息
	ClaimInfos []CoreSQLConfig.FieldsConfigType
	//描述
	Des string
}

// UploadFileToFileSysByQiniu 上传到七牛云
func UploadFileToFileSysByQiniu(c *gin.Context, args *ArgsUploadFileToFileSysByQiniu) (data BaseFileSys.FieldsFileClaimType, errCode string, err error) {
	//生成子目录结构
	nowTime := CoreFilter.GetNowTime()
	newDir := localDefaultDir + CoreFile.Sep + nowTime.Format("200601"+CoreFile.Sep+"02"+CoreFile.Sep+"15") + CoreFile.Sep
	//保存文件
	var newFileData DataUploadFileType
	newFileData, err = UploadFile(c, &ArgsUploadFile{
		FormName:   args.FormName,
		TargetSrc:  newDir,
		MaxSize:    localFileMaxSize,
		FilterType: localFileFilterType,
		IsRename:   true,
	})
	if err != nil {
		errCode = "save_temp_file"
		err = errors.New("cannot save upload file, " + err.Error())
		return
	}
	data, errCode, err = BaseQiniu.UploadBySrc(&BaseQiniu.ArgsUploadBySrc{
		Src:        newFileData.Src,
		BucketName: args.BucketName,
		FileType:   newFileData.Type,
		IP:         args.IP,
		CreateInfo: args.CreateInfo,
		UserID:     args.UserID,
		OrgID:      args.OrgID,
		IsPublic:   args.IsPublic,
		ExpireAt:   args.ExpireAt,
		ClaimInfos: args.ClaimInfos,
		Des:        args.Des,
	})
	return
}

// ArgsUploadFileToFileSysByQiniuBase64 通过base64形式上传七牛云文件参数
type ArgsUploadFileToFileSysByQiniuBase64 struct {
	//文件名称
	FileName string `json:"fileName"`
	//存储块名称
	BucketName string
	//IP地址
	IP string
	//创建来源
	CreateInfo CoreSQLFrom.FieldsFrom
	//创建用户
	UserID int64
	//创建组织
	OrgID int64
	//是否为公开的文件
	IsPublic bool `json:"isPublic" check:"bool" empty:"true"`
	//过期时间
	ExpireAt time.Time
	//扩展信息
	ClaimInfos []CoreSQLConfig.FieldsConfigType
	//描述
	Des string
}

// UploadFileToFileSysByQiniuBase64 通过base64形式上传七牛云文件参数
func UploadFileToFileSysByQiniuBase64(args *ArgsUploadFileToFileSysByQiniuBase64, fileData string) (data BaseFileSys.FieldsFileClaimType, errCode string, err error) {
	//生成子目录结构
	nowTime := CoreFilter.GetNowTime()
	newDir := localDefaultDir + CoreFile.Sep + nowTime.Format("200601"+CoreFile.Sep+"02"+CoreFile.Sep+"15") + CoreFile.Sep
	//保存文件
	var newFileData DataUploadFileType
	newFileData, err = UploadFileBase64(&ArgsUploadFileBase64{
		FileName:   args.FileName,
		TargetSrc:  newDir,
		MaxSize:    localFileMaxSize,
		FilterType: localFileFilterType,
		IsRename:   true,
	}, fileData)
	if err != nil {
		errCode = "save_temp_file"
		err = errors.New("cannot save upload file, " + err.Error())
		return
	}
	data, errCode, err = BaseQiniu.UploadBySrc(&BaseQiniu.ArgsUploadBySrc{
		Src:        newFileData.Src,
		BucketName: args.BucketName,
		FileType:   newFileData.Type,
		IP:         args.IP,
		CreateInfo: args.CreateInfo,
		UserID:     args.UserID,
		OrgID:      args.OrgID,
		IsPublic:   args.IsPublic,
		ExpireAt:   args.ExpireAt,
		ClaimInfos: args.ClaimInfos,
		Des:        args.Des,
	})
	return
}

// ArgsUploadFileToFileSys 上传到文件系统内新的文件
// 相关参数会自动绕过本地默认设定
type ArgsUploadFileToFileSys struct {
	//用户结构
	UserInfo *UserCore.DataUserDataType
	//表单名称
	FormName string
	//目标路径
	TargetSrc string
	//文件尺寸最大
	MaxSize int64
	//限制文件类型
	FilterType []string
	//创建来源
	CreateInfo CoreSQLFrom.FieldsFrom
	//过期时间
	ExpireAt time.Time
	//来源信息
	FromInfo CoreSQLFrom.FieldsFrom
	//扩展信息
	Infos []CoreSQLConfig.FieldsConfigType
	//引用扩展信息
	ClaimInfos []CoreSQLConfig.FieldsConfigType
	//描述
	Des string
}

func UploadFileToFileSys(c *gin.Context, args *ArgsUploadFileToFileSys) (newData BaseFileSys.FieldsFileClaimType, errCode string, err error) {
	//生成子目录结构
	nowTime := CoreFilter.GetNowTime()
	newDir := args.TargetSrc + CoreFile.Sep + nowTime.Format("200601"+CoreFile.Sep+"02"+CoreFile.Sep+"15") + CoreFile.Sep
	//保存文件
	var newFileData DataUploadFileType
	newFileData, err = UploadFile(c, &ArgsUploadFile{
		FormName:   args.FormName,
		TargetSrc:  newDir,
		MaxSize:    args.MaxSize,
		FilterType: args.FilterType,
		IsRename:   true,
	})
	if err != nil {
		errCode = "文件存储失败"
		err = errors.New("cannot save upload file, " + err.Error())
		return
	}
	//保存到数据库结构体
	newData, _, errCode, err = BaseFileSys.Create(&BaseFileSys.ArgsCreate{
		CreateIP:   c.ClientIP(),
		CreateInfo: args.CreateInfo,
		UserID:     args.UserInfo.Info.ID,
		OrgID:      args.UserInfo.Info.OrgID,
		FileSize:   newFileData.Size,
		FileType:   newFileData.Type,
		FileHash:   newFileData.SHA256,
		FileSrc:    newFileData.Src,
		ExpireAt:   args.ExpireAt,
		FromInfo:   args.FromInfo,
		Infos:      args.Infos,
		ClaimInfos: args.ClaimInfos,
		Des:        args.Des,
	})
	if err != nil {
		//失败后删除上传文件，并反馈失败
		return BaseFileSys.FieldsFileClaimType{}, errCode, errors.New("cannot create data, " + err.Error())
	}
	return
}

// ArgsUploadBase64ToFileSys 上传base64数据到文件系统参数
type ArgsUploadBase64ToFileSys struct {
	//文件名称
	FileName string
	//IP
	ClientIP string `json:"clientIP"`
	//创建来源
	CreateInfo CoreSQLFrom.FieldsFrom
	//创建用户
	UserID int64
	//创建组织
	OrgID int64
	//是否为公开的文件
	IsPublic bool `json:"isPublic" check:"bool" empty:"true"`
	//过期时间
	ExpireAt time.Time
	//扩展信息
	Infos []CoreSQLConfig.FieldsConfigType `json:"infos"`
	//扩展信息
	ClaimInfos []CoreSQLConfig.FieldsConfigType
	//描述
	Des string
}

// UploadBase64ToFileSys 上传base64数据到文件系统
func UploadBase64ToFileSys(args *ArgsUploadBase64ToFileSys, fileData string) (newData BaseFileSys.FieldsFileClaimType, errCode string, err error) {
	//生成子目录结构
	nowTime := CoreFilter.GetNowTime()
	newDir := localDefaultDir + CoreFile.Sep + nowTime.Format("200601"+CoreFile.Sep+"02"+CoreFile.Sep+"15") + CoreFile.Sep
	//保存文件
	var newFileData DataUploadFileType
	newFileData, err = UploadFileBase64(&ArgsUploadFileBase64{
		FileName:   args.FileName,
		TargetSrc:  newDir,
		MaxSize:    localFileMaxSize,
		FilterType: localFileFilterType,
		IsRename:   true,
	}, fileData)
	if err != nil {
		errCode = "文件存储失败"
		err = errors.New("cannot save upload file, " + err.Error())
		return
	}
	//保存到数据库结构体
	newData, _, errCode, err = BaseFileSys.Create(&BaseFileSys.ArgsCreate{
		CreateIP:   args.ClientIP,
		CreateInfo: args.CreateInfo,
		UserID:     args.UserID,
		OrgID:      args.OrgID,
		FileSize:   newFileData.Size,
		FileType:   newFileData.Type,
		FileHash:   newFileData.SHA256,
		FileSrc:    newFileData.Src,
		ExpireAt:   args.ExpireAt,
		FromInfo:   CoreSQLFrom.FieldsFrom{System: "local"},
		Infos:      args.Infos,
		ClaimInfos: args.ClaimInfos,
		Des:        args.Des,
	})
	if err != nil {
		//失败后删除上传文件，并反馈失败
		return BaseFileSys.FieldsFileClaimType{}, errCode, errors.New("cannot create data, " + err.Error())
	}
	return
}

// ArgsUploadFile 上传文件参数
type ArgsUploadFile struct {
	//表单名称
	FormName string
	//目标路径，末尾必须添加Sep
	TargetSrc string
	//文件最大大小，如果为0则不限制
	MaxSize int64
	//文件类别限制
	FilterType []string
	//是否重新命名文件名称
	IsRename bool
}

// UploadFile 上传文件参数
// 可利用该方法，实现任意文件、目标得上传
// 注意，目标路径不要交给用户设计，否则将出现重大系统级漏洞，影响系统安全
// 存储后，默认根据创建"Unix时间戳_"结构设计文件名称
// param c *gin.Context
// return FileUploadType 文件类型
// return error 错误信息
func UploadFile(c *gin.Context, args *ArgsUploadFile) (DataUploadFileType, error) {
	//检查前置并生成数据结构
	res, dataByte, err := checkUpload(c, args.FormName, args.MaxSize, args.FilterType)
	if err != nil {
		return res, err
	}
	//构建新的文件名称
	if args.IsRename {
		res.NewName = strconv.FormatInt(res.CreateTime, 10) + "_" + res.SHA256 + "." + res.Type
	} else {
		res.NewName = res.Name
	}
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

// ArgsUploadFileBase64 上传文件参数
type ArgsUploadFileBase64 struct {
	//文件名称
	FileName string
	//目标路径，末尾必须添加Sep
	TargetSrc string
	//文件最大大小，如果为0则不限制
	MaxSize int64
	//文件类别限制
	FilterType []string
	//是否重新命名文件名称
	IsRename bool
}

func UploadFileBase64(args *ArgsUploadFileBase64, fileData string) (DataUploadFileType, error) {
	//将base64数据写入临时文件
	tempDir := Router2SystemConfig.RootDir + CoreFile.Sep + "temp" + CoreFile.Sep + "base64"
	tempFileName, err := CoreFilter.GetRandStr3(30)
	if err != nil {
		return DataUploadFileType{}, errors.New(fmt.Sprint("get rand temp file name, ", err))
	}
	if !CoreFile.IsExist(tempDir) {
		if err := CoreFile.CreateFolder(tempDir); err != nil {
			return DataUploadFileType{}, errors.New(fmt.Sprint("create temp dir, ", err))
		}
	}
	tempSrc := tempDir + CoreFile.Sep + tempFileName
	//将base64数据写入临时文件
	if CoreFile.IsExist(tempSrc) {
		if err := CoreFile.DeleteF(tempSrc); err != nil {
			return DataUploadFileType{}, errors.New(fmt.Sprint("delete exist temp file, ", err))
		}
	}
	baseDec, err := base64.RawStdEncoding.DecodeString(fileData)
	if err != nil {
		baseDec, err = base64.StdEncoding.DecodeString(fileData)
		if err != nil {
			baseDec, err = base64.RawURLEncoding.DecodeString(fileData)
			if err != nil {
				baseDec, err = base64.URLEncoding.DecodeString(fileData)
				if err != nil {
					return DataUploadFileType{}, errors.New(fmt.Sprint("base64 decode by string, ", err))
				}
			}
		}
	}
	if err := CoreFile.WriteFile(tempSrc, baseDec); err != nil {
		return DataUploadFileType{}, errors.New(fmt.Sprint("save temp file, src: ", tempSrc, ", err: ", err))
	}
	//检查前置并生成数据结构
	res, dataByte, err := checkUploadSrc(tempSrc, args.FileName, args.MaxSize, args.FilterType)
	if err != nil {
		return res, err
	}
	//构建新的文件名称
	if args.IsRename {
		res.NewName = strconv.FormatInt(res.CreateTime, 10) + "_" + res.SHA256 + "." + res.Type
	} else {
		res.NewName = res.Name
	}
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

// 上传文件前置判断
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

// checkUploadSrc base64文件上传文件前置判断
func checkUploadSrc(fileSrc string, fileName string, maxSize int64, filterType []string) (DataUploadFileType, []byte, error) {
	//初始化
	var res DataUploadFileType
	//获取文件
	fileData, err := CoreFile.GetFileInfo(fileSrc)
	if err != nil {
		return res, []byte{}, errors.New(fmt.Sprint("get file info, ", err))
	}
	//判断文件尺寸
	if fileData.Size() > maxSize && maxSize > 0 {
		return res, []byte{}, errors.New("upload file size too lager")
	}
	res.Size = fileData.Size()
	//获取文件名
	res.Name = fileName
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
	dataByte, err := CoreFile.LoadFile(fileSrc)
	if err != nil {
		return res, []byte{}, err
	}
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
