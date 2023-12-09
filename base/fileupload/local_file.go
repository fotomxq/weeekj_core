package BaseFileUpload

import (
	"errors"
	CoreFile "gitee.com/weeekj/weeekj_core/v5/core/file"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	"os"
	"strconv"
	"strings"
)

// ArgsUploadFileByLocalFile 上传文件附加方法参数
type ArgsUploadFileByLocalFile struct {
	//文件路径
	FileSrc string
	//目标路径
	TargetSrc string
	//文件最大尺寸
	MaxSize int64
	//限制格式
	FilterType []string
	//是否重命名
	IsRename bool
}

// UploadFileByLocalFile 上传文件附加方法
// 将本地文件按照上传文件进行处置
func UploadFileByLocalFile(args *ArgsUploadFileByLocalFile) (DataUploadFileType, error) {
	//检查前置并生成数据结构
	res, dataByte, err := checkUploadLocalFile(args.FileSrc, args.MaxSize, args.FilterType)
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

// checkUploadLocalFile 本地文件上传潜质判断
func checkUploadLocalFile(fileSrc string, maxSize int64, filterType []string) (DataUploadFileType, []byte, error) {
	//初始化
	var res DataUploadFileType
	//获取文件
	fd, err := os.Open(fileSrc)
	if err != nil {
		return DataUploadFileType{}, nil, err
	}
	defer func() {
		if err := fd.Close(); err != nil {
			CoreLog.Error(err)
		}
	}()
	fileInfo, err := CoreFile.GetFileInfo(fileSrc)
	if err != nil {
		return DataUploadFileType{}, nil, err
	}
	//判断文件尺寸
	if fileInfo.Size() > maxSize && maxSize > 0 {
		return res, []byte{}, errors.New("upload file size too lager")
	}
	res.Size = fileInfo.Size()
	//获取文件名
	res.Name = fileInfo.Name()
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
	n, err := fd.Read(buf)
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
