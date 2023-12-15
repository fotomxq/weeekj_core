package BaseTempFile

import (
	CoreFile "github.com/fotomxq/weeekj_core/v5/core/file"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	Router2Mid "github.com/fotomxq/weeekj_core/v5/router2/mid"
	"net/http"
)

func getFileByURL(c *Router2Mid.RouterURLPublicC) (data FieldsFile, b bool) {
	//获取参数
	id := c.Context.Param("id")
	if id == "" {
		Router2Mid.ReportBaseError(c, "report_params_lost")
		return
	}
	idInt64, err := CoreFilter.GetInt64ByString(id)
	if err != nil {
		Router2Mid.ReportBaseError(c, "report_params_lost")
		return
	}
	//获取文件名称
	hash := c.Context.Param("hash")
	if hash == "" {
		Router2Mid.ReportBaseError(c, "report_params_lost")
		return
	}
	//获取文件
	data = getFileID(idInt64)
	if data.ID < 1 {
		Router2Mid.ReportWarnLog(c, "download temp file", nil, "err_no_data")
		return
	}
	if data.FileSHA1 != hash {
		Router2Mid.ReportWarnLog(c, "download temp file", nil, "err_no_data")
		return
	}
	b = true
	return
}

func DownloadFile(c *Router2Mid.RouterURLPublicC) {
	//获取文件
	data, b := getFileByURL(c)
	if !b {
		return
	}
	//下载文件
	fileData, err := CoreFile.LoadFile(data.FileSrc)
	if err != nil {
		Router2Mid.ReportWarnLog(c, "download temp file", err, "err_io")
		return
	}
	fileContentDisposition := "attachment;filename=\"" + data.Name + "\""
	c.Context.Header("Content-Type", "application/octet-stream")
	c.Context.Header("Content-Disposition", fileContentDisposition)
	c.Context.Data(http.StatusOK, "application/octet-stream", fileData)
	//反馈
	return
}

func LoadImgFile(c *Router2Mid.RouterURLPublicC) {
	//获取文件类型
	fileType := c.Context.Param("fileType")
	if fileType == "" {
		Router2Mid.ReportBaseError(c, "report_params_lost")
		return
	}
	switch fileType {
	case "jpg":
	case "jpeg":
	case "png":
	case "gif":
	default:
		Router2Mid.ReportBaseError(c, "report_params_lost")
		return
	}
	//获取文件
	data, b := getFileByURL(c)
	if !b {
		return
	}
	//下载文件
	fileData, err := CoreFile.LoadFile(data.FileSrc)
	if err != nil {
		Router2Mid.ReportWarnLog(c, "download temp file", err, "err_io")
		return
	}
	c.Context.Header("Content-Type", "image/"+fileType)
	c.Context.Data(http.StatusOK, "image/"+fileType, fileData)
	//反馈
	return
}
