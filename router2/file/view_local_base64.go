package Router2File

import (
	"encoding/base64"
	BaseFileSys "gitee.com/weeekj/weeekj_core/v5/base/filesys"
	CoreFile "gitee.com/weeekj/weeekj_core/v5/core/file"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	Router2Mid "gitee.com/weeekj/weeekj_core/v5/router2/mid"
)

// ViewLocalFileBase64 下载本地文件
func ViewLocalFileBase64(c *Router2Mid.RouterURLPublicC, fileID int64) {
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
	//读取本地文件
	fileDataByte, err := CoreFile.LoadFile(fileData.FileSrc)
	if err != nil {
		Router2Mid.ReportWarnLog(c, "download local file", err, "err_io")
		return
	}
	//转化base64
	fileDataBase64 := base64.URLEncoding.EncodeToString(fileDataByte)
	//输出数据
	Router2Mid.BaseData(c, fileDataBase64)
}
