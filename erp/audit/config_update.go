package ERPAudit

import (
	"errors"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
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
	//节点设置
	// 流程化的核心处理
	StepList FieldsConfigStepList `db:"step_list" json:"stepList"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateConfig 更新配置信息
func UpdateConfig(args *ArgsUpdateConfig) (errCode string, err error) {
	//检查节点
	errCode, err = checkStepList(args.StepList)
	if err != nil {
		return
	}
	//获取配置
	data := GetConfig(args.ID, args.OrgID)
	if data.ID < 1 {
		errCode = "err_erp_audit_config_not_exist"
		err = errors.New("no data")
		return
	}
	if data.Hash != args.Hash {
		errCode = "err_erp_audit_config_hash"
		err = errors.New("hash error")
		return
	}
	//修改节点值，避免缓冲异常
	for k, v := range args.StepList {
		newHash := CoreFilter.GetSha1Str(v.Key)
		for k2, v2 := range args.StepList {
			if v2.NextStepKey == v.Key {
				args.StepList[k2].NextStepKey = newHash
			}
			if v2.BanNextStepKey == v.Key {
				args.StepList[k2].BanNextStepKey = newHash
			}
		}
		args.StepList[k].Key = newHash
	}
	//获取hash
	newHash := getNewHash()
	//更新数据
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE erp_audit_config SET update_at = NOW(), hash = :hash, name = :name, step_list = :step_list, params = :params WHERE id = :id", map[string]interface{}{
		"id":        args.ID,
		"hash":      newHash,
		"name":      args.Name,
		"step_list": args.StepList,
		"params":    args.Params,
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
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE erp_audit_config SET update_at = NOW(), publish_at = NOW() WHERE id = :id AND org_id = :org_id", map[string]interface{}{
		"id":     args.ID,
		"org_id": args.OrgID,
	})
	if err != nil {
		return
	}
	deleteConfigCache(args.ID)
	return
}
