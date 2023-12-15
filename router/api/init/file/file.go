package RouterAPIInitFile

import (
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	CoreFileUpload "github.com/fotomxq/weeekj_core/v5/base/fileupload"
	CoreFile "github.com/fotomxq/weeekj_core/v5/core/file"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	RouterReport "github.com/fotomxq/weeekj_core/v5/router/report"
	"github.com/gin-gonic/gin"
	"strings"
)

// UpdateLocal 更新本地设置
func UpdateLocal(c *gin.Context) bool {
	//更新本地设置
	localDefaultDir, err := BaseConfig.GetDataString("FileLocalDefaultDir")
	if err != nil {
		RouterReport.BaseError(c, "config_error", "get local default dir")
		CoreLog.Error("get config data, localDefaultDir, ", err)
		return false
	}
	localDefaultDir = CoreFile.BaseSrc + CoreFile.Sep + localDefaultDir
	localFileMaxSize, err := BaseConfig.GetDataInt64("FileLocalFileMaxSize")
	if err != nil {
		RouterReport.BaseError(c, "config_error", "get local file max size")
		CoreLog.Error("get config data, localFileMaxSize, ", err)
		return false
	}
	localFileFilterTypeStr, err := BaseConfig.GetDataString("FileLocalFileFilterType")
	if err != nil {
		RouterReport.BaseError(c, "config_error", "get file filter type")
		CoreLog.Error("get config data, localFileFilterType, ", err)
		return false
	}
	CoreFileUpload.SetLocalConfig(localDefaultDir, localFileMaxSize, strings.Split(localFileFilterTypeStr, ","))
	return true
}
