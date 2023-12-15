package Router2File

import (
	CoreFile "github.com/fotomxq/weeekj_core/v5/core/file"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	Router2Mid "github.com/fotomxq/weeekj_core/v5/router2/mid"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"net/http"
)

// DownloadLocalFile 下载本地文件
func DownloadLocalFile(c *Router2Mid.RouterURLPublicC) {
	//获取目录
	tempChildDir := c.Context.Param("dir")
	if tempChildDir == "" {
		Router2Mid.ReportBaseError(c, "report_params_lost")
		return
	}
	//获取文件名称
	tempFileName := c.Context.Param("temp_file")
	if tempFileName == "" {
		Router2Mid.ReportBaseError(c, "report_params_lost")
		return
	}
	//检查temp合法性和文件是否存在
	if !CoreFilter.CheckFileName(tempChildDir) || !CoreFilter.CheckFileName(tempFileName) {
		Router2Mid.ReportBaseError(c, "report_params_error")
		return
	}
	tempFile := Router2SystemConfig.RootDir + CoreFile.Sep + "temp" + CoreFile.Sep + tempChildDir + CoreFile.Sep + tempFileName
	if b := CoreFile.IsFile(tempFile); !b {
		Router2Mid.ReportWarnLog(c, "download local file,temp file:"+tempFile, nil, "err_io")
		return
	}
	//下载文件
	fileData, err := CoreFile.LoadFile(tempFile)
	if err != nil {
		Router2Mid.ReportWarnLog(c, "download local file", err, "err_io")
		return
	}
	fileContentDisposition := "attachment;filename=\"" + tempFileName + "\""
	c.Context.Header("Content-Type", "application/octet-stream")
	c.Context.Header("Content-Disposition", fileContentDisposition)
	c.Context.Data(http.StatusOK, "application/octet-stream", fileData)
}
