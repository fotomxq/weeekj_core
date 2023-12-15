package FinanceReportForm

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetFileList 获取文件列表参数
type ArgsGetFileList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//采用模版
	TemplateID int64 `db:"template_id" json:"templateID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetFileList 获取文件列表
func GetFileList(args *ArgsGetFileList) (dataList []FieldsFile, dataCount int64, err error) {
	maps := map[string]interface{}{}
	where := ""
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > 0 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.TemplateID > 0 {
		where = where + " AND template_id = :template_id"
		maps["template_id"] = args.TemplateID
	}
	if args.Search != "" {
		where = where + " AND (name ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "finance_report_form_file"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, update_at, delete_at, org_id, name, des, template_id, params FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	return
}

// ArgsCreateFile 创建新的文件参数
type ArgsCreateFile struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"300" empty:"true"`
	//采用模版
	TemplateID int64 `db:"template_id" json:"templateID" check:"id" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateFile 创建新的文件
func CreateFile(args *ArgsCreateFile) (data FieldsFile, err error) {
	var templateData FieldsTemplate
	if args.TemplateID > 0 {
		err = Router2SystemConfig.MainDB.Get(&templateData, "SELECT id, col_ids, params FROM finance_report_form_template WHERE id = $1 AND delete_at < to_timestamp(1000000)", args.TemplateID)
		if err != nil {
			return
		}
	}
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "finance_report_form_file", "INSERT INTO finance_report_form_file (org_id, name, des, template_id, params) VALUES (:org_id,:name,:des,:template_id,:params)", args, &data)
	if err != nil {
		return
	}
	if templateData.ID > 0 {
		//自动为该文件构建val数据集
		var colList []FieldsCol
		err = Router2SystemConfig.MainDB.Select(&colList, "SELECT id, mark, params FROM finance_report_form_col WHERE id = ANY($1) AND delete_at < to_timestamp(1000000)", templateData.ColIDs)
		var appendVals []interface{}
		for _, v := range colList {
			appendVals = append(appendVals, map[string]interface{}{
				"org_id":  args.OrgID,
				"file_id": data.ID,
				"col_id":  v.ID,
				"mark":    v.Mark,
				"params":  CoreSQLConfig.FieldsConfigsType{},
			})
		}
		//依次建立val
		err = CoreSQL.CreateMore(Router2SystemConfig.MainDB.DB, "INSERT INTO finance_report_form_val (org_id, file_id, col_id, mark, val, val_float, val_int, params) VALUES (:org_id,:file_id,:col_id,:mark,'',0,0,:params)", appendVals)
		if err != nil {
			return
		}
	}
	return
}

// ArgsUpdateFile 修改文件信息参数
type ArgsUpdateFile struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"300" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateFile 修改文件信息
func UpdateFile(args *ArgsUpdateFile) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE finance_report_form_file SET name = :name, des = :des, params = :params WHERE id = :id AND org_id = :org_id", args)
	return
}

// ArgsDeleteFile 删除文件参数
type ArgsDeleteFile struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// DeleteFile 删除文件
func DeleteFile(args *ArgsDeleteFile) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "finance_report_form_file", "id = :id AND org_id = :org_id", args)
	return
}
