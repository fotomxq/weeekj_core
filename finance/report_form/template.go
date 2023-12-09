package FinanceReportForm

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetTemplateList 获取模板列表参数
type ArgsGetTemplateList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetTemplateList 获取模板列表
func GetTemplateList(args *ArgsGetTemplateList) (dataList []FieldsTemplate, dataCount int64, err error) {
	maps := map[string]interface{}{}
	where := ""
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.Search != "" {
		where = where + " AND (name ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "finance_report_form_template"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, update_at, delete_at, org_id, name, des, cover_files, col_ids, params FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	return
}

// ArgsCreateTemplate 创建新的模板参数
type ArgsCreateTemplate struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"300" empty:"true"`
	//封面图
	CoverFiles pq.Int64Array `db:"cover_files" json:"coverFiles" check:"ids" empty:"true"`
	//列数据列
	ColIDs pq.Int64Array `db:"col_ids" json:"colIDs" check:"ids" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateTemplate 创建新的模板
func CreateTemplate(args *ArgsCreateTemplate) (data FieldsTemplate, err error) {
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "finance_report_form_template", "INSERT INTO finance_report_form_template (org_id, name, des, cover_files, col_ids, params) VALUES (:org_id,:name,:des,:cover_files,:col_ids,:params)", args, &data)
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
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"300" empty:"true"`
	//封面图
	CoverFiles pq.Int64Array `db:"cover_files" json:"coverFiles" check:"ids" empty:"true"`
	//列数据列
	ColIDs pq.Int64Array `db:"col_ids" json:"colIDs" check:"ids" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateTemplate 修改模板
func UpdateTemplate(args *ArgsUpdateTemplate) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE finance_report_form_template SET name = :name, des = :des, cover_files = :cover_files, col_ids = :col_ids, params = :params WHERE id = :id AND org_id = :org_id", args)
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
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "finance_report_form_template", "id = :id AND org_id = :org_id", args)
	return
}
