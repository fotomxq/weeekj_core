package CoreExcel

//该模块用于实现便捷的excel操作
// 也可以绕过该模块，直接操作相关表格

import (
	"github.com/360EntSecGroup-Skylar/excelize"
)

//NewFile 建立新文件
func NewFile() *excelize.File {
	return excelize.NewFile()
}

//LoadFile 打开一个新的文件
func LoadFile(src string) (*excelize.File, error) {
	return excelize.OpenFile(src)
}

//SaveFile 保存修改结果
func SaveFile(excel *excelize.File, src string) error {
	return excel.SaveAs(src)
}

//GetSheetRows 获取子表
func GetSheetRows(excel *excelize.File, sheetName string) [][]string {
	if sheetName == "" {
		sheetName = "Sheet1"
	}
	return excel.GetRows(sheetName)
}

//GetCellValue 读取某列行内容
func GetCellValue(excel *excelize.File, sheetName string, axis string) string {
	return excel.GetCellValue(sheetName, axis)
}

//SetCellString 写入某列行内容
func SetCellString(excel *excelize.File, sheetName string, axis string, value string) {
	excel.SetCellStr(sheetName, axis, value)
}

//AddCell 加入列
func AddCell(excel *excelize.File, sheetName string, column string) {
	excel.InsertCol(sheetName, column)
}

//AddRow 加入行
func AddRow(excel *excelize.File, sheetName string, row int) {
	excel.InsertRow(sheetName, row)
}
