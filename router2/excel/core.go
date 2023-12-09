package Router2Excel

import (
	"fmt"
	BaseTempFile "gitee.com/weeekj/weeekj_core/v5/base/temp_file"
	CoreExcel "gitee.com/weeekj/weeekj_core/v5/core/excel"
	CoreFile "gitee.com/weeekj/weeekj_core/v5/core/file"
	Router2Mid "gitee.com/weeekj/weeekj_core/v5/router2/mid"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/360EntSecGroup-Skylar/excelize"
)

//excel快速组装方法集
// 提供一整套的方案，可用于快速构建和完成excel导出工作

type ExcelQuick struct {
	//路由头部
	C any
	//前置错误日志
	LogErr string
	//缓冲名称组
	// eg: fmt.Sprint("erp_permanent_assets_", orgID, "_end_", endAtCarbon.Time.Format("2006-01-02"))
	FileParams string
	//文件名称
	// eg: fmt.Sprint("固定资产", endAtCarbon.Time.Format("2006"), "年度变动情况.xlsx")
	FileName string
	//模板路径
	// 该路径位于builds的excel目录下
	// eg: fmt.Sprint("erp", CoreFile.Sep, "permanent_assets", CoreFile.Sep, "sort_change_analysis.xlsx")
	TemplatePath string
	//是否覆盖样式
	// 仅可用于单个表，如果是多个表，请手动调用QuickSetStyle方法
	NeedReplaceStyle bool
	//覆盖样式参考表名称
	// eg: ""，给空则采用默认表名称Sheet1
	ReplaceStyleSheet string
	//覆盖样式参考位置
	// eg: "A1"
	ReplaceStyleRef string
	//覆盖样式起点位置
	// eg: fmt.Sprint("A", 1)
	ReplaceStyleStart string
	//覆盖样式结束位置
	// eg: fmt.Sprint("O", rowStep)
	ReplaceStyleEnd string
	//Excel数据对象
	ExcelObj *excelize.File
	//缓冲数据保留秒
	CacheSaveTime int
}

// InitCache 准备数据集
func (t *ExcelQuick) InitCache() (b bool) {
	b = t.beforeLoadParamsFile()
	if b {
		return
	}
	if err := t.getTemplate(); err != nil {
		b = true
		Router2Mid.ReportWarnLog(t.C, t.LogErr+", load excel template failed, ", err, "err_excel_template")
		return
	}
	return
}

// QuickInsertCol 快入给单元格写入数据
func (t *ExcelQuick) QuickInsertCol(sheetName string, data map[string]string) {
	if sheetName == "" {
		sheetName = "Sheet1"
	}
	for k, v := range data {
		CoreExcel.SetCellString(t.ExcelObj, sheetName, k, v)
	}
}

// QuickSetStyle 快速设置样式
func (t *ExcelQuick) QuickSetStyle(sheetName string, defaultStyle string, areaStart string, areaEnd string) {
	if sheetName == "" {
		sheetName = "Sheet1"
	}
	if defaultStyle == "" {
		defaultStyle = "A1"
	}
	styleID := t.ExcelObj.GetCellStyle(sheetName, defaultStyle)
	t.ExcelObj.SetCellStyle(sheetName, areaStart, areaEnd, styleID)
}

// Done 完成处理
func (t *ExcelQuick) Done() (b bool) {
	t.autoQuickSetStyle()
	if err := t.saveTemplateExcel(); err != nil {
		return
	}
	b = true
	return
}

// 预先下载文件处理
func (t *ExcelQuick) beforeLoadParamsFile() (result bool) {
	//获取数据
	newID, hash, b := BaseTempFile.SaveFileBefore(t.FileParams)
	if !b {
		return
	}
	//反馈数据
	Router2Mid.ReportData(t.C, t.LogErr+", ", nil, "", map[string]interface{}{
		"id":   newID,
		"hash": hash,
	})
	return true
}

// 获取模版文件
func (t *ExcelQuick) getTemplate() (err error) {
	t.ExcelObj, err = CoreExcel.LoadFile(fmt.Sprint(Router2SystemConfig.RootDir, CoreFile.Sep+"data"+CoreFile.Sep, "excel", CoreFile.Sep, t.TemplatePath))
	if err != nil {
		return
	}
	return
}

// autoQuickSetStyle 速设置样式
func (t *ExcelQuick) autoQuickSetStyle() {
	if !t.NeedReplaceStyle {
		return
	}
	t.QuickSetStyle(t.ReplaceStyleSheet, t.ReplaceStyleRef, t.ReplaceStyleStart, t.ReplaceStyleEnd)
}

// 第二代保存excel文件
func (t *ExcelQuick) saveTemplateExcel() error {
	if t.CacheSaveTime < 1 {
		t.CacheSaveTime = 60
	}
	fileSrc, newID, hash, err := BaseTempFile.SaveFile(t.CacheSaveTime, t.FileParams, t.FileName, "", "xlsx")
	if err != nil {
		Router2Mid.ReportWarnLog(t.C, t.LogErr+", save temp file, ", err, "err_make_file")
		return err
	}
	if err := CoreExcel.SaveFile(t.ExcelObj, fileSrc); err != nil {
		Router2Mid.ReportWarnLog(t.C, t.LogErr+", , ", err, "err_make_file")
		return err
	}
	//反馈数据
	Router2Mid.ReportData(t.C, t.LogErr+", ", nil, "", map[string]interface{}{
		"id":   newID,
		"hash": hash,
	})
	return nil
}
