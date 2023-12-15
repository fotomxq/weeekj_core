package ERPDocument

import (
	"errors"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

type ArgsDeleteDoc struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

func DeleteDoc(args *ArgsDeleteDoc) (errCode string, err error) {
	//获取数据
	data := getDocByID(args.ID)
	if data.ID < 1 {
		errCode = "err_no_data"
		err = errors.New("no data")
		return
	}
	configData := getConfigByID(data.ConfigID)
	if configData.ID < 1 {
		errCode = "err_erp_document_config_not_exist"
		err = errors.New("no config data")
		return
	}
	//删除数据
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "erp_document_doc", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		errCode = "err_delete"
		return
	}
	//清理缓冲
	deleteDocCache(args.ID)
	//删除附加数据
	switch configData.DocType {
	case "custom":
		//修改组件列
		err = docComponentValObj.DeleteByBindID(data.ID)
		if err != nil {
			errCode = "err_delete"
			return
		}
	case "doc":
	case "excel":
		//获取excel配置
		err = deleteExcelByDocID(data.ID)
		if err != nil {
			errCode = "err_delete"
			return
		}
	}
	//反馈
	return
}
