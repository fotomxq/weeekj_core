package ERPDocument

import (
	"fmt"
	CoreCache "gitee.com/weeekj/weeekj_core/v5/core/cache"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// GetExcelAllSheet 获取文档所有子表
func GetExcelAllSheet(docID int64) (dataList []FieldsExcelSheet) {
	var rawList []FieldsExcelSheet
	err := Router2SystemConfig.MainDB.Select(&rawList, "SELECT id FROM erp_document_excel_sheet WHERE doc_id = $1", docID)
	if err != nil {
		return
	}
	for _, v := range rawList {
		dataList = append(dataList, getExcelSheetBySheetID(docID, v.ID))
	}
	return
}

// GetExcelSheetData 获取子表
func GetExcelSheetData(docID int64, sheetID int64) (dataList []FieldsExcelRowCol) {
	dataList = getExcelSheetDataBySheetID(docID, sheetID)
	return
}

// ArgsCreateExcelSheet 创建子表参数
type ArgsCreateExcelSheet struct {
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id"`
	//文档ID
	DocID int64 `db:"doc_id" json:"docID" check:"id"`
	//名称
	Name string `db:"name" json:"name" check:"name" min:"1" max:"300"`
	//数据集合
	// 兼容性数据集合，可剔除掉数值类内容（或不需要统计的内容），放入本集合内
	Data string `db:"data" json:"data" check:"des" min:"1" max:"60000" empty:"true"`
}

// CreateExcelSheet 创建子表
func CreateExcelSheet(args *ArgsCreateExcelSheet) (data FieldsExcelSheet, err error) {
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "erp_document_excel_sheet", "INSERT INTO erp_document_excel_sheet(config_id, doc_id, name, data) VALUES(:config_id, :doc_id, :name, :data)", args, &data)
	if err != nil {
		return
	}
	return
}

// ArgsUpdateExcelSheet 修改子表参数
type ArgsUpdateExcelSheet struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id"`
	//文档ID
	DocID int64 `db:"doc_id" json:"docID" check:"id"`
	//名称
	Name string `db:"name" json:"name" check:"name" min:"1" max:"300"`
	//数据集合
	// 兼容性数据集合，可剔除掉数值类内容（或不需要统计的内容），放入本集合内
	Data string `db:"data" json:"data" check:"des" min:"1" max:"60000" empty:"true"`
}

// UpdateExcelSheet 修改子表
func UpdateExcelSheet(args *ArgsUpdateExcelSheet) (err error) {
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE erp_document_excel_sheet SET name = :name, data = :data WHERE id = :id AND config_id = :config_id AND doc_id = :doc_id", args)
	if err != nil {
		return
	}
	deleteExcelSheetCache(args.DocID, args.ID)
	return
}

// ArgsDeleteExcelSheet 删除子表参数
type ArgsDeleteExcelSheet struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id"`
	//文档ID
	DocID int64 `db:"doc_id" json:"docID" check:"id"`
}

// DeleteExcelSheet 删除子表
func DeleteExcelSheet(args *ArgsDeleteExcelSheet) (err error) {
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "erp_document_excel_sheet", "id = :id AND config_id = :config_id AND doc_id = :doc_id", args)
	if err != nil {
		return
	}
	deleteExcelSheetCache(args.DocID, args.ID)
	_ = DeleteExcelSheetData(args.DocID, args.ID)
	return
}

// ArgsSetExcelSheetData 设置子表数据集参数
type ArgsSetExcelSheetData struct {
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id"`
	//文档ID
	DocID int64 `db:"doc_id" json:"docID" check:"id"`
	//所属文档子表
	SheetID int64 `db:"sheet_id" json:"sheetID"`
	//数据集合
	DataList []ArgsSetExcelSheetDataChild `json:"dataList"`
}
type ArgsSetExcelSheetDataChild struct {
	//位置
	Row string `db:"row" json:"row"`
	Col string `db:"col" json:"col"`
	//默认样式
	ClassName string `db:"class_name" json:"className"`
	//样式
	StyleName string `db:"style_name" json:"styleName"`
	//组件默认值
	Val string `db:"val" json:"val"`
	//整数（内部记录用）
	ValInt64 int64 `db:"val_int64" json:"valInt64"`
	//浮点数（内部记录用）
	ValFloat64 float64 `db:"val_float64" json:"valFloat64"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// SetExcelSheetData 设置子表数据集
func SetExcelSheetData(args *ArgsSetExcelSheetData) (err error) {
	dataList := getExcelSheetDataBySheetID(args.DocID, args.SheetID)
	var waitCreate []ArgsSetExcelSheetDataChild
	var waitUpdate []FieldsExcelRowCol
	for _, v := range args.DataList {
		isFind := false
		for _, v2 := range dataList {
			if v.Row == v2.Row && v.Col == v2.Col {
				isFind = true
				v2 = FieldsExcelRowCol{
					ID:         v2.ID,
					ConfigID:   v2.ConfigID,
					DocID:      v2.DocID,
					SheetID:    v2.SheetID,
					Row:        v2.Row,
					Col:        v2.Col,
					ClassName:  v.ClassName,
					StyleName:  v.StyleName,
					Val:        v.Val,
					ValInt64:   v.ValInt64,
					ValFloat64: v.ValFloat64,
					Params:     v.Params,
				}
				waitUpdate = append(waitUpdate, v2)
				break
			}
		}
		if !isFind {
			waitCreate = append(waitCreate, v)
		}
	}
	if len(waitCreate) > 0 {
		for _, v := range waitCreate {
			_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO erp_document_excel_row_col(config_id, doc_id, sheet_id, row, col, class_name, style_name, val, val_int64, val_float64, params) VALUES (:config_id, :doc_id, :sheet_id, :row, :col, :class_name, :style_name, :val, :val_int64, :val_float64, :params)", map[string]interface{}{
				"config_id":   args.ConfigID,
				"doc_id":      args.DocID,
				"sheet_id":    args.SheetID,
				"row":         v.Row,
				"col":         v.Col,
				"class_name":  v.ClassName,
				"style_name":  v.StyleName,
				"val":         v.Val,
				"val_int64":   v.ValInt64,
				"val_float64": v.ValFloat64,
				"params":      v.Params,
			})
			if err != nil {
				return
			}
		}
	}
	if len(waitUpdate) > 0 {
		for _, v := range waitUpdate {
			_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE erp_document_excel_row_col SET class_name = :class_name, style_name = :style_name, val = :val, val_int64 = :val_int64, val_float64 = :val_float64, params = :params WHERE id = :id", map[string]interface{}{
				"id":          v.ID,
				"class_name":  v.ClassName,
				"style_name":  v.StyleName,
				"val":         v.Val,
				"val_int64":   v.ValInt64,
				"val_float64": v.ValFloat64,
				"params":      v.Params,
			})
			if err != nil {
				return
			}
		}
	}
	deleteExcelDataCache(args.DocID, args.SheetID)
	return
}

// DeleteExcelSheetData 删除子表所有数据集
func DeleteExcelSheetData(docID int64, sheetID int64) (err error) {
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "erp_document_excel_row_col", "doc_id = :doc_id AND sheet_id = :sheet_id", map[string]interface{}{
		"doc_id":   docID,
		"sheet_id": sheetID,
	})
	if err != nil {
		return
	}
	deleteExcelDataCache(docID, sheetID)
	return
}

// setExcelByConfigID 通过配置初始化子表及数据集
func setExcelByConfigID(configID int64, docID int64) (err error) {
	excelConfigData := GetExcelConfigByConfigID(configID)
	if excelConfigData.ID < 1 {
		return
	}
	for _, v := range excelConfigData.Sheets {
		var vSheet FieldsExcelSheet
		vSheet, err = CreateExcelSheet(&ArgsCreateExcelSheet{
			ConfigID: excelConfigData.ConfigID,
			DocID:    docID,
			Name:     v.Name,
			Data:     v.Data,
		})
		if err != nil {
			return
		}
		var vDataArgs []ArgsSetExcelSheetDataChild
		for _, v2 := range v.RowCols {
			vDataArgs = append(vDataArgs, ArgsSetExcelSheetDataChild{
				Row:        v2.Row,
				Col:        v2.Col,
				ClassName:  v2.ClassName,
				StyleName:  v2.StyleName,
				Val:        v2.Val,
				ValInt64:   v2.ValInt64,
				ValFloat64: v2.ValFloat64,
				Params:     v2.Params,
			})
		}
		if len(vDataArgs) > 0 {
			err = SetExcelSheetData(&ArgsSetExcelSheetData{
				ConfigID: excelConfigData.ConfigID,
				DocID:    docID,
				SheetID:  vSheet.ID,
				DataList: vDataArgs,
			})
			if err != nil {
				return
			}
		}
	}
	return
}

// 删除文档的所有数据
func deleteExcelByDocID(docID int64) (err error) {
	sheetList := GetExcelAllSheet(docID)
	if len(sheetList) < 1 {
		return
	}
	for _, v := range sheetList {
		err = DeleteExcelSheet(&ArgsDeleteExcelSheet{
			ID:       v.ID,
			ConfigID: v.ConfigID,
			DocID:    v.DocID,
		})
		if err != nil {
			return
		}
	}
	return
}

// getExcelSheetBySheetID 获取子表数据
func getExcelSheetBySheetID(docID, sheetID int64) (data FieldsExcelSheet) {
	cacheMark := getExcelSheetCacheMark(docID, sheetID)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, config_id, doc_id, name, data FROM erp_document_excel_sheet WHERE doc_id = $1 AND id = $2", docID, sheetID)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, CoreCache.CacheTime1Day)
	return
}

// getExcelSheetDataBySheetID 获取子表内部数据
func getExcelSheetDataBySheetID(docID, sheetID int64) (dataList []FieldsExcelRowCol) {
	cacheMark := getExcelSheetDataCacheMark(docID, sheetID)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &dataList); err == nil && len(dataList) > 0 {
		return
	}
	err := Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, config_id, doc_id, sheet_id, row, col, class_name, style_name, val, val_int64, val_float64, params FROM erp_document_excel_row_col WHERE doc_id = $1 AND sheet_id = $2", docID, sheetID)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, dataList, CoreCache.CacheTime1Hour)
	return
}

func getExcelSheetCacheMark(docID, sheetID int64) string {
	return fmt.Sprint("erp:document:excel:sheet:id:", docID, ".", sheetID)
}

func deleteExcelSheetCache(docID, sheetID int64) {
	Router2SystemConfig.MainCache.DeleteMark(getExcelSheetCacheMark(docID, sheetID))
}

func getExcelSheetDataCacheMark(docID, sheetID int64) string {
	return fmt.Sprint("erp:document:excel:data:id:", docID, ".", sheetID)
}

func deleteExcelDataCache(docID, id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getExcelSheetDataCacheMark(docID, id))
}
