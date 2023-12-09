package ERPDocument

import (
	"errors"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

type ArgsDeleteConfig struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

func DeleteConfig(args *ArgsDeleteConfig) (err error) {
	//获取数据
	data := getConfigByID(args.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	//删除数据
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "erp_document_config", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	//清理缓冲
	deleteConfigCache(args.ID)
	//根据类型判断处理
	switch data.DocType {
	case "custom":
	case "doc":
	case "excel":
		_ = deleteExcelConfig(data.ID)
	}
	//反馈
	return
}
