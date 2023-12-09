package ERPAudit

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

type ArgsCreateConfig struct {
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

func CreateConfig(args *ArgsCreateConfig) (errCode string, err error) {
	//检查节点
	errCode, err = checkStepList(args.StepList)
	if err != nil {
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
	//获取新的值
	newHash := getNewHash()
	//写入数据
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO erp_audit_config (publish_at, hash, org_id, name, step_list, params) VALUES (to_timestamp(0), :hash, :org_id, :name, :step_list, :params)", map[string]interface{}{
		"hash":      newHash,
		"org_id":    args.OrgID,
		"name":      args.Name,
		"step_list": args.StepList,
		"params":    args.Params,
	})
	if err != nil {
		errCode = "err_insert"
		return
	}
	//反馈
	return
}
