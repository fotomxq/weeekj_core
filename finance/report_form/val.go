package FinanceReportForm

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetValByFile 获取制定文件的值数据列参数
type ArgsGetValByFile struct {
	//文件ID
	FileID int64 `db:"file_id" json:"fileID" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// GetValByFile 获取制定文件的值数据列
func GetValByFile(args *ArgsGetValByFile) (dataList []FieldsVal, err error) {
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, create_at, update_at, delete_at, org_id, file_id, col_id, mark, val, val_float, val_int, params FROM finance_report_form_val WHERE org_id = $1 AND file_id = $2 AND delete_at < to_timestamp(1000000)", args.OrgID, args.FileID)
	return
}

// ArgsSetVal 设置值参数
type ArgsSetVal struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//所属文件
	FileID int64 `db:"file_id" json:"fileID" check:"id"`
	//列ID
	ColID int64 `db:"col_id" json:"colID" check:"id"`
	//标识码
	// 位置标识码：A1\B1，字母代表列、数字代表行
	Mark string `db:"mark" json:"mark" check:"mark"`
	//值
	Val      string  `db:"val" json:"val"`
	ValFloat float64 `db:"val_float" json:"valFloat"`
	ValInt   int64   `db:"val_int" json:"valInt"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// SetVal 设置值
func SetVal(args *ArgsSetVal) (data FieldsVal, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, file_id, col_id, mark, val, val_float, val_int, params FROM finance_report_form_val WHERE file_id = $1 AND org_id = $2 AND col_id = $3 AND delete_at < to_timestamp(1000000)", args.FileID, args.OrgID, args.ColID)
	if err == nil && data.ID > 0 {
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE finance_report_form_val SET val = :val, val_float = :val_float, val_int = :val_int, params = :params WHERE id = :id", map[string]interface{}{
			"id":        data.ID,
			"val":       args.Val,
			"val_float": args.ValFloat,
			"val_int":   args.ValInt,
			"params":    args.Params,
		})
		if err != nil {
			return
		}
		data.Val = args.Val
		data.ValFloat = args.ValFloat
		data.ValInt = args.ValInt
		data.Params = args.Params
		return
	} else {
		err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "finance_report_form_val", "INSERT INTO finance_report_form_val (org_id, file_id, col_id, mark, val, val_float, val_int, params) VALUES (:org_id,:file_id,:col_id,:mark,:val,:val_float,:val_int,:params)", args, &data)
		return
	}
}

// ArgsDeleteVal 删除值参数
type ArgsDeleteVal struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// DeleteVal 删除值
func DeleteVal(args *ArgsDeleteVal) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "finance_report_form_val", "id = :id AND org_id = :org_id", args)
	return
}
