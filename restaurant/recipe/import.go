package RestaurantRecipe

import (
	"errors"
	"github.com/360EntSecGroup-Skylar/excelize"
	BaseFileUpload "github.com/fotomxq/weeekj_core/v5/base/fileupload"
	ClassSort "github.com/fotomxq/weeekj_core/v5/class/sort"
	CoreExcel "github.com/fotomxq/weeekj_core/v5/core/excel"
	CoreFile "github.com/fotomxq/weeekj_core/v5/core/file"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	ToolsLoadExcel "github.com/fotomxq/weeekj_core/v5/tools/load_excel"
	"github.com/gin-gonic/gin"
)

// ArgsImportData 批量导入数据参数
type ArgsImportData struct {
	//分公司ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//门店ID
	// 暂时不启用
	StoreID int64 `db:"store_id" json:"storeID" check:"id" empty:"true"`
}

// ImportData 批量导入数据
// excelData 为文件的excel文件结构体，请参考ToolsLoadExcel.UploadFileAndGetExcelData实现获取
// waitDeleteFile 为导入完成后需要删除的文件路径
func ImportData(args *ArgsImportData, excelData *excelize.File, waitDeleteFile string) (errCode string, importCount int, skipImportCount int, err error) {
	//获取子表，只拿取第一张子表，其他子表不做导入
	sheetMaps := excelData.GetSheetMap()
	var sheetName string
	for _, v := range sheetMaps {
		sheetName = v
		break
	}
	//获取相关数据值
	excelVals := CoreExcel.GetSheetRows(excelData, sheetName)
	if len(excelVals) < 1 {
		errCode = "err_excel"
		err = errors.New("excel vals is empty")
		return
	}
	//遍历获取数据并导入
	// 先构建分类，然后构建菜品
	step := 0
	for {
		step += 1
		if step >= len(excelVals) {
			//fmt.Println("step >= len(excelVals), ", step, ", ", len(excelVals))
			break
		}
		rows := excelVals[step]
		if len(rows) < 3 {
			//fmt.Println("rows len < 3, ", step)
			break
		}
		if rows[0] == "" || rows[1] == "" {
			//fmt.Println(rows, step)
			break
		}
		//检查分类是否存在
		findSortData, _ := Sort.GetByName(args.OrgID, rows[0])
		if findSortData.ID < 1 && len(rows) > 2 {
			//创建分类
			_, err = Sort.Create(&ClassSort.ArgsCreate{
				BindID:      args.OrgID,
				Mark:        "",
				ParentID:    0,
				CoverFileID: 0,
				DesFiles:    nil,
				Name:        rows[0],
				Des:         "",
				Params:      nil,
			})
			if err != nil {
				errCode = "err_insert"
				return
			}
			findSortData, _ = Sort.GetByName(args.OrgID, rows[0])
		}
		if findSortData.ID < 1 {
			CoreLog.Error("restaurant recipe import data failed, sort not found: ", rows[0])
			skipImportCount += 1
			continue
		}
		//检查菜品是否存在
		if rows[1] != "" {
			findRecipeData := GetRecipeByName(args.OrgID, -1, rows[1])
			if findRecipeData.ID < 1 {
				//梳理价格
				var vPrice int64 = 0
				vPrice = CoreFilter.GetInt64ByStringNoErr(rows[2])
				//创建菜品
				_, err = CreateRecipe(&ArgsCreateRecipe{
					CategoryID: findSortData.ID,
					Name:       rows[1],
					Unit:       "",
					UnitID:     0,
					OrgID:      args.OrgID,
					StoreID:    0,
					Price:      vPrice,
					Remark:     "",
				})
				if err != nil {
					errCode = "err_insert"
					return
				}
			} else {
				//更新菜品
				err = UpdateRecipe(&ArgsUpdateRecipe{
					ID:         findRecipeData.ID,
					CategoryID: findSortData.ID,
					Name:       rows[1],
					Unit:       findRecipeData.Unit,
					UnitID:     0,
					OrgID:      findRecipeData.OrgID,
					StoreID:    findRecipeData.StoreID,
					Price:      CoreFilter.GetInt64ByStringNoErr(rows[2]),
					Remark:     findRecipeData.Remark,
				})
				if err != nil {
					errCode = "err_update"
					return
				}
			}
			importCount += 1
		}
	}
	//删除临时文件
	err = CoreFile.DeleteF(waitDeleteFile)
	if err != nil {
		CoreLog.Error("restaurant recipe import data failed, delete temp file: ", waitDeleteFile, ", err: ", err)
		err = nil
	}
	//反馈
	return
}

// ImportDataByUpload 上传文件并导入
func ImportDataByUpload(args *ArgsImportData, c *gin.Context, argsUploadTemp *BaseFileUpload.ArgsUploadToTemp) (errCode string, importCount int, skipImportCount int, err error) {
	var excelData *excelize.File
	var waitDeleteFile string
	excelData, waitDeleteFile, errCode, err = ToolsLoadExcel.UploadFileAndGetExcelData(c, argsUploadTemp)
	if err != nil {
		return
	}
	errCode, importCount, skipImportCount, err = ImportData(args, excelData, waitDeleteFile)
	if err != nil {
		return
	}
	return
}
