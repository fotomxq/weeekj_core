package BaseFileUpload

import (
	CoreFile "github.com/fotomxq/weeekj_core/v5/core/file"
	"github.com/gin-gonic/gin"
)

// ArgsUploadToTemp 将上传文件存储到临时目录中参数
type ArgsUploadToTemp struct {
	//表单名称
	FormName string
	//文件尺寸限制
	MaxSize int64
	//限制格式
	FilterType []string
	//是否重命名
	IsRename bool
}

// UploadToTemp 将上传文件存储到临时目录中
func UploadToTemp(c *gin.Context, args *ArgsUploadToTemp) (DataUploadFileType, error) {
	tempDir := CoreFile.BaseDir() + CoreFile.Sep + "upload_temp"
	if b := CoreFile.IsExist(tempDir); !b {
		if err := CoreFile.CreateFolder(tempDir); err != nil {
			return DataUploadFileType{}, err
		}
	}
	return UploadFile(c, &ArgsUploadFile{
		FormName:   args.FormName,
		TargetSrc:  tempDir,
		MaxSize:    args.MaxSize,
		FilterType: args.FilterType,
		IsRename:   args.IsRename,
	})
}
