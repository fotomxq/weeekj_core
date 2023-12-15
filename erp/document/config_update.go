package ERPDocument

import (
	"errors"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	ERPCore "github.com/fotomxq/weeekj_core/v5/erp/core"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsUpdateConfig 更新配置信息参数
type ArgsUpdateConfig struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//hash
	// 如果hash和提交hash不同，服务端将自动拒绝更新，避免流处理异常
	Hash string `db:"hash" json:"hash" check:"sha1"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//流程名称
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

// UpdateConfig 更新配置信息
func UpdateConfig(args *ArgsUpdateConfig) (errCode string, err error) {
	//检查文档类型
	if !checkDocType(args.DocType) {
		errCode = "err_erp_document_config_doc_type"
		err = errors.New("no support doc type")
		return
	}
	//获取配置
	data := GetConfig(args.ID, args.OrgID)
	if data.ID < 1 {
		errCode = "err_erp_document_config_not_exist"
		err = errors.New("no data")
		return
	}
	if data.Hash != args.Hash {
		errCode = "err_erp_document_config_hash"
		err = errors.New("hash error")
		return
	}
	//获取hash
	newHash := getNewHash()
	//更新数据
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE erp_document_config SET update_at = NOW(), hash = :hash, name = :name, des = :des, cover_file_id = :cover_file_id, doc_type = :doc_type, component_list = :component_list, list_show = :list_show, params = :params WHERE id = :id", map[string]interface{}{
		"id":             args.ID,
		"hash":           newHash,
		"name":           args.Name,
		"des":            args.Des,
		"cover_file_id":  args.CoverFileID,
		"doc_type":       args.DocType,
		"component_list": args.ComponentList,
		"list_show":      args.ListShow,
		"params":         args.Params,
	})
	if err != nil {
		errCode = "err_update"
		return
	}
	//删除缓冲
	deleteConfigCache(args.ID)
	//反馈
	return
}

// ArgsUpdateConfigPublish 发布配置参数
type ArgsUpdateConfigPublish struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// UpdateConfigPublish 发布配置
func UpdateConfigPublish(args *ArgsUpdateConfigPublish) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE erp_document_config SET update_at = NOW(), publish_at = NOW() WHERE id = :id AND org_id = :org_id", map[string]interface{}{
		"id":     args.ID,
		"org_id": args.OrgID,
	})
	if err != nil {
		return
	}
	deleteConfigCache(args.ID)
	return
}
