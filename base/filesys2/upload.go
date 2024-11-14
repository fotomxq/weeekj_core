package BaseFileSys2

import (
	"github.com/gin-gonic/gin"
	"time"
)

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
	//newFileData, err = CoreFile.SaveUploadFileToTemp(c, newDir, "file", args.MaxSize, args.FilterType)
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
