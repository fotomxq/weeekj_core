package Router2File

import (
	"fmt"
	BaseFileSys "gitee.com/weeekj/weeekj_core/v5/base/filesys"
	CoreFile "gitee.com/weeekj/weeekj_core/v5/core/file"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	Router2Mid "gitee.com/weeekj/weeekj_core/v5/router2/mid"
	"net/http"
)

// ViewLocalFile 下载本地文件
func ViewLocalFile(c *Router2Mid.RouterURLPublicC, fileID int64) {
	//获取文件数据
	claimData, err := BaseFileSys.GetFileClaimByID(&BaseFileSys.ArgsGetFileClaimByID{
		ClaimID: fileID,
		UserID:  -1,
		OrgID:   -1,
	})
	if err != nil {
		Router2Mid.ReportWarnLog(c, "get file claim", err, "err_file_valid")
		return
	}
	//获取实体文件数据
	fileData, err := BaseFileSys.GetFileByID(&BaseFileSys.ArgsGetFileByID{
		ID:         claimData.FileID,
		CreateInfo: CoreSQLFrom.FieldsFrom{},
	})
	if err != nil {
		Router2Mid.ReportWarnLog(c, "get file claim", err, "err_file_valid")
		return
	}
	//下载文件
	fileDataByte, err := CoreFile.LoadFile(fileData.FileSrc)
	if err != nil {
		Router2Mid.ReportWarnLog(c, "download local file", err, "err_io")
		return
	}
	fileName := claimData.Des
	if fileName == "" {
		fileName = fmt.Sprint(claimData.ID)
	}
	fileName = fmt.Sprint(fileName, ".", fileData.FileType)
	fileContentDisposition := "attachment;filename=\"" + fileName + "\""
	c.Context.Header("Content-Type", "application/octet-stream")
	c.Context.Header("Content-Disposition", fileContentDisposition)
	c.Context.Data(http.StatusOK, "application/octet-stream", fileDataByte)
}
