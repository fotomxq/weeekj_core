package FinanceReportForm

import (
	"errors"
	"fmt"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLIDs "github.com/fotomxq/weeekj_core/v5/core/sql/ids"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetColList 获取列的列表参数
type ArgsGetColList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetColList 获取列的列表
func GetColList(args *ArgsGetColList) (dataList []FieldsCol, dataCount int64, err error) {
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
	tableName := "finance_report_form_col"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, update_at, delete_at, org_id, mark, name, des, params FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	return
}

// ArgsGetColByTemplate 通过模版获取列参数
type ArgsGetColByTemplate struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetColByTemplate 通过模版获取列
func GetColByTemplate(args *ArgsGetColByTemplate) (dataList []FieldsCol, err error) {
	var templateData FieldsTemplate
	err = Router2SystemConfig.MainDB.Get(&templateData, "SELECT id, col_ids FROM finance_report_form_template WHERE id = $1 AND org_id = $2 AND delete_at < to_timestamp(1000000)", args.ID, args.OrgID)
	if err != nil || templateData.ID < 1 {
		err = errors.New(fmt.Sprint("template not exist, ", err))
		return
	}
	return GetColMore(&ArgsGetColMore{
		IDs:        templateData.ColIDs,
		HaveRemove: false,
	})
}

// ArgsGetColMore 获取多个列参数
type ArgsGetColMore struct {
	//ID列
	IDs pq.Int64Array `json:"ids" check:"ids"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove" check:"bool"`
}

// GetColMore 获取多个列
func GetColMore(args *ArgsGetColMore) (dataList []FieldsCol, err error) {
	err = CoreSQLIDs.GetIDsAndDelete(&dataList, "finance_report_form_col", "id, create_at, update_at, delete_at, org_id, mark, name, des, params", args.IDs, args.HaveRemove)
	return
}

// ArgsSetCol 设置列参数
type ArgsSetCol struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//标识码
	Mark string `db:"mark" json:"mark" check:"mark"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"300" empty:"true"`
	//值类型
	ValType int `db:"val_type" json:"valType"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// SetCol 设置列
func SetCol(args *ArgsSetCol) (data FieldsCol, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, mark, name, des, params FROM finance_report_form_col WHERE org_id = $1 AND mark = $2 AND delete_at < to_timestamp(1000000)", args.OrgID, args.Mark)
	if err == nil && data.ID > 0 {
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE finance_report_form_col SET update_at = NOW(), name = :name, des = :des, val_type = :val_type, params = :params WHERE id = :id", map[string]interface{}{
			"id":       data.ID,
			"name":     args.Name,
			"des":      args.Des,
			"val_type": args.ValType,
			"params":   args.Params,
		})
		if err != nil {
			return
		}
		data.Name = args.Name
		data.Des = args.Des
		data.Params = args.Params
		return
	} else {
		err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "finance_report_form_col", "INSERT INTO finance_report_form_col (org_id, mark, name, des, val_type, params) VALUES (:org_id,:mark,:name,:des,:val_type,:params)", args, &data)
		return
	}
}

// ArgsDeleteCol 删除列参数
type ArgsDeleteCol struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// DeleteCol 删除列
func DeleteCol(args *ArgsDeleteCol) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "finance_report_form_col", "id = :id AND org_id = :org_id", args)
	return
}
