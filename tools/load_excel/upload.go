package ToolsLoadExcel

import (
	"errors"
	BaseFileUpload "gitee.com/weeekj/weeekj_core/v5/base/fileupload"
	CoreExcel "gitee.com/weeekj/weeekj_core/v5/core/excel"
	CoreFile "gitee.com/weeekj/weeekj_core/v5/core/file"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
)

// UploadFileAndGetExcelData 上传到临时文件并读取为数据包
func UploadFileAndGetExcelData(c *gin.Context, args *BaseFileUpload.ArgsUploadToTemp) (excelData *excelize.File, waitDeleteFile string, errCode string, err error) {
	//上传文件
	var fileData BaseFileUpload.DataUploadFileType
	fileData, err = BaseFileUpload.UploadToTemp(c, args)
	if err != nil {
		errCode = "err_upload"
		return
	}
	//读取excel文件
	excelData, err = CoreExcel.LoadFile(fileData.Src)
	if err != nil {
		errCode = "err_io"
		return
	}
	//读取数据
	sheetMaps := excelData.GetSheetMap()
	if len(sheetMaps) < 1 {
		_ = CoreFile.DeleteF(fileData.Src)
		errCode = "err_excel"
		err = errors.New("excel type have err")
		return
	}
	waitDeleteFile = fileData.Src
	//反馈数据包
	return
}
