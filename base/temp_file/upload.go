package BaseTempFile

import (
	CoreFile "github.com/fotomxq/weeekj_core/v5/core/file"
	"github.com/gin-gonic/gin"
)

// UploadFileToTemp 将上传文件存储到临时文件
func UploadFileToTemp(c *gin.Context, formName string, maxSize int64, filterType []string) (result CoreFile.DataGetUploadFileData, err error) {
	//TODO: 上传文件转临时文件
	return
}
