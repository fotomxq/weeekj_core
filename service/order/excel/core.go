package ServiceOrderExcel

import (
	"fmt"
	BaseTempFile "gitee.com/weeekj/weeekj_core/v5/base/temp_file"
	CoreExcel "gitee.com/weeekj/weeekj_core/v5/core/excel"
	CoreFile "gitee.com/weeekj/weeekj_core/v5/core/file"
	Router2Mid "gitee.com/weeekj/weeekj_core/v5/router2/mid"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/360EntSecGroup-Skylar/excelize"
)

// 预先下载文件处理
func beforeLoadParamsFile(c any, logErr string, fileParams string) (result bool) {
	//获取数据
	newID, hash, b := BaseTempFile.SaveFileBefore(fileParams)
	if !b {
		return
	}
	//反馈数据
	Router2Mid.ReportData(c, logErr+", ", nil, "", map[string]interface{}{
		"id":   newID,
		"hash": hash,
	})
	return true
}

// 快入给单元格写入数据
func quickInsertCol(excelData *excelize.File, sheetName string, data map[string]string) {
	if sheetName == "" {
		sheetName = "Sheet1"
	}
	for k, v := range data {
		CoreExcel.SetCellString(excelData, sheetName, k, v)
	}
}

// 获取模版文件
func getTemplate(filename string) (excelData *excelize.File, err error) {
	excelData, err = CoreExcel.LoadFile(fmt.Sprint(Router2SystemConfig.RootDir, CoreFile.Sep+"data"+CoreFile.Sep, "excel", CoreFile.Sep, filename))
	if err != nil {
		return
	}
	return
}

// 快速设置样式
func quickSetStyle(excelData *excelize.File, sheetName string, defaultStyle string, areaStart string, areaEnd string) {
	if sheetName == "" {
		sheetName = "Sheet1"
	}
	if defaultStyle == "" {
		defaultStyle = "A1"
	}
	styleID := excelData.GetCellStyle(sheetName, defaultStyle)
	excelData.SetCellStyle(sheetName, areaStart, areaEnd, styleID)
}

// 第二代保存excel文件
func saveTemplateExcel2(c any, logErr string, fileParams string, fileName string, excelData *excelize.File) error {
	fileSrc, newID, hash, err := BaseTempFile.SaveFile(60, fileParams, fileName, "", "xlsx")
	if err != nil {
		Router2Mid.ReportWarnLog(c, logErr+", save temp file, ", err, "err_make_file")
		return err
	}
	if err := CoreExcel.SaveFile(excelData, fileSrc); err != nil {
		Router2Mid.ReportWarnLog(c, logErr+", , ", err, "err_make_file")
		return err
	}
	//反馈数据
	Router2Mid.ReportData(c, logErr+", ", nil, "", map[string]interface{}{
		"id":   newID,
		"hash": hash,
	})
	return nil
}
