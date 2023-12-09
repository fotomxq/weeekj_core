package ERPAudit

import (
	"errors"
	"fmt"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	ERPCore "gitee.com/weeekj/weeekj_core/v5/erp/core"
	OrgCore "gitee.com/weeekj/weeekj_core/v5/org/core"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsCreateStep 创建审批参数
type ArgsCreateStep struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//流程配置
	ConfigID int64 `db:"config_id" json:"configID" check:"id"`
	//流程名称
	Name string `db:"name" json:"name" check:"name" min:"1" max:"300"`
	//创建成员
	CreateOrgBindID int64 `db:"create_org_bind_id" json:"createOrgBindID" check:"id" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateStep 创建审批
func CreateStep(args *ArgsCreateStep) (stepData FieldsStep, errCode string, err error) {
	//锁定
	auditCreateLock.Lock()
	defer auditCreateLock.Unlock()
	//获取审批配置
	configData := GetConfig(args.ConfigID, args.OrgID)
	if configData.ID < 1 || !CoreSQL.CheckTimeHaveData(configData.PublishAt) {
		errCode = "err_erp_audit_config_not_exist"
		err = errors.New("config not exist")
		return
	}
	//检查是否发布
	if !checkConfigPublish(args.ConfigID, args.OrgID) {
		errCode = "err_erp_audit_config_not_exist"
		err = errors.New("config no publish")
		return
	}
	//必须存在至少一个节点
	if len(configData.StepList) < 1 {
		errCode = "err_erp_audit_config_no_step"
		err = errors.New("config no have step")
		return
	}
	//构建审批
	// 构建SN
	var newSN int64
	newSN = getAuditCountByOrgID(configData.OrgID) + 1
	var viewOrgBinds, editOrgBinds pq.Int64Array
	// 组装人员数据
	viewOrgBinds = append(viewOrgBinds, args.CreateOrgBindID)
	editOrgBinds = append(editOrgBinds, args.CreateOrgBindID)
	for _, v := range configData.StepList {
		vBinds := OrgCore.GetBindByGroupAndRole(configData.OrgID, []int64{v.AuditOrgBindGroup}, v.AuditOrgRoleIDs)
		if len(vBinds) < 1 {
			continue
		}
		for _, v2 := range vBinds {
			viewOrgBinds = CoreFilter.MargeNoReplaceArrayInt64(viewOrgBinds, v2.ID)
			if v.AuditMode == "all" || v.AuditMode == "only" {
				editOrgBinds = CoreFilter.MargeNoReplaceArrayInt64(editOrgBinds, v2.ID)
			}
		}
		for _, v2 := range v.AuditOrgBindIDs {
			viewOrgBinds = CoreFilter.MargeNoReplaceArrayInt64(viewOrgBinds, v2)
			if v.AuditMode == "all" || v.AuditMode == "only" {
				editOrgBinds = CoreFilter.MargeNoReplaceArrayInt64(editOrgBinds, v2)
			}
		}
	}
	// 获取开始节点
	var nowStepChildKey string
	for _, v := range configData.StepList {
		isFind := false
		for _, v2 := range configData.StepList {
			if v.Key == v2.NextStepKey || v.Key == v2.BanNextStepKey {
				isFind = true
				break
			}
		}
		if !isFind {
			nowStepChildKey = v.Key
			break
		}
	}
	// 写入数据
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "erp_audit_step", "INSERT INTO erp_audit_step (finish_at, finish_status, org_id, config_id, sn, name, create_org_bind_id, can_view_org_bind_ids, can_edit_org_bind_ids, have_org_bind_ids, now_step_child_key, params) VALUES (to_timestamp(0),0,:org_id,:config_id,:sn,:name,:create_org_bind_id,:can_view_org_bind_ids,:can_edit_org_bind_ids,:have_org_bind_ids,:now_step_child_key,:params)", map[string]interface{}{
		"org_id":                configData.OrgID,
		"config_id":             configData.ID,
		"sn":                    newSN,
		"name":                  args.Name,
		"create_org_bind_id":    args.CreateOrgBindID,
		"can_view_org_bind_ids": viewOrgBinds,
		"can_edit_org_bind_ids": editOrgBinds,
		"have_org_bind_ids":     pq.Int64Array{},
		"now_step_child_key":    nowStepChildKey,
		"params":                args.Params,
	}, &stepData)
	if err != nil {
		errCode = "err_insert"
		err = errors.New(fmt.Sprint("create step by config id: ", configData.ID, ", err: ", err))
		return
	}
	//构建审批节点
	for kStepChild := 0; kStepChild < len(configData.StepList); kStepChild++ {
		//获取当前节点
		vStepChild := configData.StepList[kStepChild]
		//获取过期时间
		vExpireAt := CoreFilter.GetNowTimeCarbon().AddSeconds(vStepChild.ExpireSec).Time
		//等待审批人员
		var waitAuditOrgBinds pq.Int64Array
		for _, v2 := range vStepChild.AuditOrgBindIDs {
			waitAuditOrgBinds = append(waitAuditOrgBinds, v2)
		}
		addAuditOrgBinds := OrgCore.GetBindByGroupAndRole(stepData.OrgID, []int64{vStepChild.AuditOrgBindGroup}, vStepChild.AuditOrgRoleIDs)
		for _, v2 := range addAuditOrgBinds {
			waitAuditOrgBinds = append(waitAuditOrgBinds, v2.ID)
		}
		//写入节点
		var vStepChildData FieldsStepChild
		err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "erp_audit_step_child", "INSERT INTO erp_audit_step_child (step_id, audit_at, ban_at, expire_at, key, name, audit_mode, audit_org_bind_group, audit_org_bind_ids, audit_org_role_ids, wait_audit_org_bind_ids, finish_audit_org_binds, next_step_key, ban_next_step_key, params) VALUES (:step_id,to_timestamp(0),to_timestamp(0),:expire_at,:key,:name,:audit_mode,:audit_org_bind_group,:audit_org_bind_ids,:audit_org_role_ids,:wait_audit_org_bind_ids,:finish_audit_org_binds,:next_step_key,:ban_next_step_key,:params)", map[string]interface{}{
			"step_id":                 stepData.ID,
			"expire_at":               vExpireAt,
			"key":                     vStepChild.Key,
			"name":                    vStepChild.Name,
			"audit_mode":              vStepChild.AuditMode,
			"audit_org_bind_group":    vStepChild.AuditOrgBindGroup,
			"audit_org_bind_ids":      vStepChild.AuditOrgBindIDs,
			"audit_org_role_ids":      vStepChild.AuditOrgRoleIDs,
			"wait_audit_org_bind_ids": waitAuditOrgBinds,
			"finish_audit_org_binds":  pq.Int64Array{},
			"next_step_key":           vStepChild.NextStepKey,
			"ban_next_step_key":       vStepChild.BanNextStepKey,
			"params":                  vStepChild.Params,
		}, &vStepChildData)
		if err != nil {
			errCode = "err_insert"
			err = errors.New(fmt.Sprint("create step child, child key: ", vStepChild.Key, ", err: ", err))
			return
		}
		//批量写入组件
		err = componentValObj.SetMore(&ERPCore.ArgsSetMore{
			BindID:   vStepChildData.ID,
			DataList: vStepChild.ComponentList,
		})
		if err != nil {
			errCode = "err_insert"
			return
		}
	}
	//处理通知
	sendStepAudit(stepData.ID)
	//反馈
	return
}

// getAuditCountByOrgID 获取审批数量
func getAuditCountByOrgID(orgID int64) (count int64) {
	_ = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM erp_audit_step WHERE org_id = $1", orgID)
	return
}
