package RestaurantWeeklyRecipeMarge

import (
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	Router2Excel "github.com/fotomxq/weeekj_core/v5/router2/excel"
	Router2Mid "github.com/fotomxq/weeekj_core/v5/router2/mid"
)

type ArgsPrintExcelRaw struct {
	//每周菜谱ID
	WeeklyRecipeID int64 `db:"weekly_recipe_id" json:"weeklyRecipeID" check:"id" index:"true"`
}

// PrintExcelRaw 导出周菜谱原材料
func PrintExcelRaw(c any, logErr string, args *ArgsPrintExcelRaw) {
	//菜谱信息
	weeklyRecipeData := getWeeklyRecipeByID(args.WeeklyRecipeID)
	if weeklyRecipeData.ID < 1 {
		Router2Mid.ReportWarnLog(c, logErr+", no data, ", nil, "err_no_data")
		return
	}
	//获取原材料列表
	rawList, _ := GetWeeklyRecipeRaw(&ArgsGetWeeklyRecipeRaw{
		WeeklyRecipeID: args.WeeklyRecipeID,
	})
	if len(rawList) < 1 {
		Router2Mid.ReportWarnLog(c, logErr+", no data, ", nil, "err_no_data")
		return
	}
	//预先调用数据
	excelTools := Router2Excel.ExcelTemplate{}
	excelTools.SetFileHash(fmt.Sprint("restaurant_weekly_recipe_raw_", args, "_", CoreFilter.GetNowTimeCarbon().Format("200601021504")))
	if excelTools.BeforeLoadParamsFile(c, logErr) {
		return
	}
	//文件名称
	excelTools.SetFileName(fmt.Sprint("周菜单分化单_", CoreFilter.GetNowTimeCarbon().Format("20060102_1504"), ".xlsx"))
	//读取模板
	excelTools.SetSubDir("restaurant_weekly_recipe_print")
	excelFile, err := excelTools.GetExcelTemplate(c, logErr, "restaurant_weekly_recipe_raw.xlsx")
	if err != nil {
		return
	}
	//主要操作表名称
	mainSheetName := "Sheet1"
	//获取样式
	styleA4, _ := excelFile.GetCellStyle(mainSheetName, "A4")
	styleC4, _ := excelFile.GetCellStyle(mainSheetName, "C4")
	styleE4, _ := excelFile.GetCellStyle(mainSheetName, "E4")
	styleG4, _ := excelFile.GetCellStyle(mainSheetName, "G4")
	//当前行
	rowIndex := 4
	//序号
	index := 1
	//遍历资产列表
	for k := 0; k < len(rawList); k++ {
		vItemData := rawList[k]
		//设置序号
		_ = excelFile.SetCellValue(mainSheetName, fmt.Sprint("A", rowIndex), index)
		//时间段
		_ = excelFile.SetCellValue(mainSheetName, fmt.Sprint("B", rowIndex), fmt.Sprint(vItemData.DiningDate))
		//菜品
		_ = excelFile.SetCellValue(mainSheetName, fmt.Sprint("C", rowIndex), vItemData.RecipeName)
		//合并单元格
		_ = excelFile.MergeCell(mainSheetName, fmt.Sprint("C", rowIndex), fmt.Sprint("D", rowIndex))
		//原材料
		_ = excelFile.SetCellValue(mainSheetName, fmt.Sprint("E", rowIndex), vItemData.MaterialName)
		//合并单元格
		_ = excelFile.MergeCell(mainSheetName, fmt.Sprint("E", rowIndex), fmt.Sprint("F", rowIndex))
		//用量
		_ = excelFile.SetCellValue(mainSheetName, fmt.Sprint("G", rowIndex), vItemData.Count)
		//合并单元格
		_ = excelFile.MergeCell(mainSheetName, fmt.Sprint("G", rowIndex), fmt.Sprint("H", rowIndex))
		//设置样式
		_ = excelFile.SetCellStyle(mainSheetName, fmt.Sprint("A", rowIndex), fmt.Sprint("B", rowIndex), styleA4)
		_ = excelFile.SetCellStyle(mainSheetName, fmt.Sprint("C", rowIndex), fmt.Sprint("D", rowIndex), styleC4)
		_ = excelFile.SetCellStyle(mainSheetName, fmt.Sprint("E", rowIndex), fmt.Sprint("F", rowIndex), styleE4)
		_ = excelFile.SetCellStyle(mainSheetName, fmt.Sprint("G", rowIndex), fmt.Sprint("H", rowIndex), styleG4)
		//行叠加
		rowIndex += 1
		//叠加序号
		index += 1
	}
	//清理单元格缓冲
	_ = excelFile.UpdateLinkedValue()
	//设置工作簿的默认工作表
	excelFile.SetActiveSheet(0)
	//保存Excel文件
	err = excelTools.SaveExcelTemplate(c, logErr)
	if err != nil {
		return
	}
	//反馈
	return
}
