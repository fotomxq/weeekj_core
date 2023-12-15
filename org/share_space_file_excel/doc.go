package OrgShareSpaceFileExcel

import (
	"encoding/json"
	"errors"
	"fmt"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	OrgShareSpaceMod "github.com/fotomxq/weeekj_core/v5/org/share_space/mod"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// GetDocByID 查看文档
func GetDocByID(id int64) (data FieldsDoc, err error) {
	data = getDocByID(id)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsCreateDoc 创建新文档参数
type ArgsCreateDoc struct {
	//采用模板
	TemplateID int64 `db:"template_id" json:"templateID"`
	//结构体设计
	SheetData FieldsSheetList `db:"sheet_data" json:"sheetData"`
}

// CreateDoc 创建新文档
func CreateDoc(args *ArgsCreateDoc) (data FieldsDoc, err error) {
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "org_share_space_file_excel_doc", "INSERT INTO org_share_space_file_excel_doc(template_id, sheet_data) VALUES(:template_id, :sheet_data)", args, &data)
	if err != nil {
		return
	}
	updateDocSize(data.ID)
	return
}

// ArgsUpdateDoc 修改文档参数
type ArgsUpdateDoc struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//结构体设计
	SheetData FieldsSheetList `db:"sheet_data" json:"sheetData"`
}

// UpdateDoc 修改文档
func UpdateDoc(args *ArgsUpdateDoc) (err error) {
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE org_share_space_file_excel_doc SET sheet_data = :sheet_data WHERE id = :id", args)
	if err != nil {
		return
	}
	deleteDocCache(args.ID)
	updateDocSize(args.ID)
	return
}

// deleteDoc 删除文档
func deleteDoc(id int64) (err error) {
	_, err = CoreSQL.DeleteOne(Router2SystemConfig.MainDB.DB, "org_share_space_file_excel_doc", "id", map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return
	}
	deleteDocCache(id)
	return
}

// 获取文档数据
func getDocByID(id int64) (data FieldsDoc) {
	cacheMark := getDocCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, template_id, sheet_data FROM org_share_space_file_excel_doc WHERE id = $1", id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 1800)
	return
}

// 缓冲
func getDocCacheMark(id int64) string {
	return fmt.Sprint("org:share:space:file:excel:doc:id:", id)
}

func deleteDocCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getDocCacheMark(id))
}

// 更新文档的尺寸
func updateDocSize(id int64) {
	data := getDocByID(id)
	if data.ID < 1 {
		return
	}
	dataByte, err := json.Marshal(data.SheetData)
	if err != nil {
		return
	}
	fileSize := int64(len(dataByte))
	OrgShareSpaceMod.UpdateFileSize("excel", data.ID, fileSize)
}
