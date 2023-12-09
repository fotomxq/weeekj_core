package OrgShareSpaceFileExcel

import (
	"errors"
	"fmt"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsGetTemplateList 获取模板列表参数
type ArgsGetTemplateList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织
	OrgID int64 `json:"orgID" check:"id" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetTemplateList 获取模板列表
func GetTemplateList(args *ArgsGetTemplateList) (dataList []FieldsTemplate, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.OrgID > -1 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.Search != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "(name ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	if where == "" {
		where = "true"
	}
	tableName := "org_share_space_file_excel_template"
	var rawList []FieldsTemplate
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		"SELECT id FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getTemplateByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, FieldsTemplate{
			ID:        vData.ID,
			CreateAt:  vData.CreateAt,
			OrgID:     vData.OrgID,
			Name:      vData.Name,
			SheetData: nil,
		})
	}
	return
}

// ArgsGetTemplateByID 查看模板参数
type ArgsGetTemplateByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// GetTemplateByID 查看模板
func GetTemplateByID(args *ArgsGetTemplateByID) (data FieldsTemplate, err error) {
	data = getTemplateByID(args.ID)
	if data.ID < 1 || !CoreFilter.EqID2(args.OrgID, data.OrgID) {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsCreateTemplate 创建新模板参数
type ArgsCreateTemplate struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//结构体设计
	SheetData FieldsSheetList `db:"sheet_data" json:"sheetData"`
}

// CreateTemplate 创建新模板
func CreateTemplate(args *ArgsCreateTemplate) (data FieldsTemplate, err error) {
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "org_share_space_file_excel_template", "INSERT INTO org_share_space_file_excel_template(org_id, name, sheet_data) VALUES(:org_id, :name, :sheet_data)", args, &data)
	if err != nil {
		return
	}
	return
}

// ArgsUpdateTemplate 修改模板参数
type ArgsUpdateTemplate struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//结构体设计
	SheetData FieldsSheetList `db:"sheet_data" json:"sheetData"`
}

// UpdateTemplate 修改模板
func UpdateTemplate(args *ArgsUpdateTemplate) (err error) {
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE org_share_space_file_excel_template SET name = :name, sheet_data = :sheet_data WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	deleteTemplateCache(args.ID)
	return
}

// ArgsDeleteTemplate 删除模板参数
type ArgsDeleteTemplate struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// DeleteTemplate 删除模板
func DeleteTemplate(args *ArgsDeleteTemplate) (err error) {
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "org_share_space_file_excel_template", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	deleteTemplateCache(args.ID)
	return
}

// 获取模板数据
func getTemplateByID(id int64) (data FieldsTemplate) {
	cacheMark := getTemplateCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, org_id, name, sheet_data FROM org_share_space_file_excel_template WHERE id = $1", id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 1800)
	return
}

// 缓冲
func getTemplateCacheMark(id int64) string {
	return fmt.Sprint("org:share:space:file:excel:template:id:", id)
}

func deleteTemplateCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getTemplateCacheMark(id))
}
