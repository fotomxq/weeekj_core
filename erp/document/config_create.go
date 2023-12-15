package ERPDocument

import (
	"errors"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	ERPCore "github.com/fotomxq/weeekj_core/v5/erp/core"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

type ArgsCreateConfig struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//名称
	Name string `db:"name" json:"name" check:"name" min:"1" max:"300"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"6000" empty:"true"`
	//封面ID
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//文档类型
	// custom 自定义; doc 普通文稿; excel 表格
	DocType string `db:"doc_type" json:"docType"`
	//节点组件
	ComponentList ERPCore.FieldsComponentDefineList `db:"component_list" json:"componentList"`
	//列表展示数据
	ListShow FieldsConfigListShows `db:"list_show" json:"listShow"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

func CreateConfig(args *ArgsCreateConfig) (errCode string, err error) {
	//检查文档类型
	if !checkDocType(args.DocType) {
		errCode = "err_erp_document_config_doc_type"
		err = errors.New("no support doc type")
		return
	}
	//获取新的值
	newHash := getNewHash()
	//写入数据
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO erp_document_config (publish_at, hash, org_id, name, des, cover_file_id, doc_type, component_list, list_show, params) VALUES (to_timestamp(0), :hash, :org_id, :name, :des, :cover_file_id, :doc_type, :component_list, :list_show, :params)", map[string]interface{}{
		"hash":           newHash,
		"org_id":         args.OrgID,
		"name":           args.Name,
		"des":            args.Des,
		"cover_file_id":  args.CoverFileID,
		"doc_type":       args.DocType,
		"component_list": args.ComponentList,
		"list_show":      args.ListShow,
		"params":         args.Params,
	})
	if err != nil {
		errCode = "err_insert"
		return
	}
	//反馈
	return
}
