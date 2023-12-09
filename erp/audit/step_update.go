package ERPAudit

import (
	"errors"
	"fmt"
	BaseExpireTip "gitee.com/weeekj/weeekj_core/v5/base/expire_tip"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	ERPCore "gitee.com/weeekj/weeekj_core/v5/erp/core"
	OrgWorkTipMod "gitee.com/weeekj/weeekj_core/v5/org/work_tip/mod"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsUpdateStepAudit 审批目标节点参数
type ArgsUpdateStepAudit struct {
	//ID
	StepID int64 `json:"stepID" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//审批人员
	OrgBindID int64 `json:"orgBindID" check:"id"`
	//是否通过审批
	IsAudit bool `json:"isAudit" check:"bool"`
	//填写内容
	DataList []ERPCore.ArgsComponentValSetOnlyUpdate `json:"dataList"`
}

// UpdateStepAudit 审批目标节点
/**
1. 只有审批类才能使用本接口，否则将反馈失败
2. 已经完成的禁止使用本接口
*/
func UpdateStepAudit(args *ArgsUpdateStepAudit) (errCode string, err error) {
	//锁定机制
	auditUpdateLock.Lock()
	defer auditUpdateLock.Unlock()
	//获取节点数据
	stepData := GetStep(args.StepID, args.OrgID, args.OrgBindID)
	if stepData.ID < 1 {
		errCode = "err_erp_audit_step_no_data"
		return
	}
	//检查是否发布
	if !checkConfigPublish(stepData.ConfigID, args.OrgID) {
		errCode = "err_erp_audit_config_not_exist"
		err = errors.New("config no publish")
		return
	}
	//找到当前节点内容
	stepChildData := getStepChildByKey(stepData.ID, stepData.NowStepChildKey)
	if stepChildData.ID < 1 {
		errCode = "err_erp_audit_step_child_no_data"
		err = errors.New(fmt.Sprint("get step child failed, step id: ", stepData.ID, ", now step child key: ", stepData.NowStepChildKey, ", err: ", err))
		return
	}
	stepChildData = getStepChildByID(stepChildData.ID)
	//检查是否可以审批该节点？
	if !CoreFilter.CheckInt64InArray(stepChildData.WaitAuditOrgBindIDs, args.OrgBindID) {
		errCode = "err_erp_audit_step_no_permission"
		err = errors.New("no permission audit")
		return
	}
	//记录所有填入内容
	err = componentValObj.SetValMoreOnlyUpdate(&ERPCore.ArgsComponentValMoreSetOnlyUpdate{
		BindID:   stepChildData.ID,
		DataList: args.DataList,
	})
	if err != nil {
		errCode = "err_erp_audit_step_child_component_update"
		err = errors.New(fmt.Sprint("set step child component, step child id: ", stepChildData.ID, ", err: ", err))
		return
	}
	//增加审批人
	stepChildData.FinishAuditOrgBinds = CoreFilter.MargeNoReplaceArrayInt64(stepChildData.FinishAuditOrgBinds, args.OrgBindID)
	if args.IsAudit {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE erp_audit_step_child SET audit_at = NOW(), finish_audit_org_binds = :finish_audit_org_binds WHERE id = :id", map[string]interface{}{
			"id":                     stepChildData.ID,
			"finish_audit_org_binds": stepChildData.FinishAuditOrgBinds,
		})
	} else {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE erp_audit_step_child SET ban_at = NOW(), finish_audit_org_binds = :finish_audit_org_binds WHERE id = :id", map[string]interface{}{
			"id":                     stepChildData.ID,
			"finish_audit_org_binds": stepChildData.FinishAuditOrgBinds,
		})
	}
	if err != nil {
		errCode = "err_update"
		return
	}
	deleteStepChildCache(stepChildData.ID)
	//处理节点
	sendStepAudit(stepData.ID)
	//反馈
	return
}

// 处理节点通知
/**
1. 如果是通知抄送类，将发送工作提醒后，更新节点到下一个
2. 如果是审批类，将固定节点并发送工作提醒
3. 自动会轮动到下一个节点，如果需要
*/
func sendStepAudit(stepID int64) {
	//日志
	appendLog := "erp audit send step audit, "
	//获取节点数据
	stepData := getStepByID(stepID)
	if stepData.ID < 1 {
		return
	}
	//找到当前节点内容
	stepChildData := getStepChildByKey(stepData.ID, stepData.NowStepChildKey)
	if stepChildData.ID < 1 {
		CoreLog.Error(appendLog, "get step child failed, step id: ", stepData.ID, ", now step child key: ", stepData.NowStepChildKey)
		return
	}
	stepChildData = getStepChildByID(stepChildData.ID)
	//是否自动进入下一个流程
	autoNext := false
	//需要抄送的人员列
	var sendOrgBinds []int64
	//需要审批通知的人员列
	var auditOrgBinds []int64
	//是否需要具体判断审批是否通过？
	needCheckAuditOK := false
	//识别节点类型
	switch stepChildData.AuditMode {
	case "none":
		//记录跳过
		autoNext = true
	case "all":
		//需要审批
		auditOrgBinds = stepChildData.WaitAuditOrgBindIDs
		//如果全部完成审批，则标记可下一步
		autoNext = true
		for _, v := range stepChildData.WaitAuditOrgBindIDs {
			isFind := false
			for _, v2 := range stepChildData.FinishAuditOrgBinds {
				if v == v2 {
					isFind = true
					break
				}
			}
			if !isFind {
				autoNext = false
				break
			}
		}
		//需判断具体通过情况
		needCheckAuditOK = true
	case "only":
		//需要审批
		auditOrgBinds = stepChildData.WaitAuditOrgBindIDs
		//如果已经存在审批人，则标记可下一步
		autoNext = len(stepChildData.FinishAuditOrgBinds) > 0
		//需判断具体通过情况
		needCheckAuditOK = true
	case "send":
		//抄送跳过
		sendOrgBinds = stepChildData.WaitAuditOrgBindIDs
		autoNext = true
	}
	//开始抄送处理
	if len(sendOrgBinds) > 0 {
		sendMsg := fmt.Sprint("有一个审批[", stepData.Name, "]等待您查阅")
		for _, v := range sendOrgBinds {
			OrgWorkTipMod.AppendTip(&OrgWorkTipMod.ArgsAppendTip{
				OrgID:     stepData.ID,
				OrgBindID: v,
				Msg:       sendMsg,
				System:    "erp_audit",
				BindID:    stepData.ID,
			})
		}
	}
	//通知需要审批人员
	if len(auditOrgBinds) > 0 {
		auditMsg := fmt.Sprint("有一个审批[", stepData.Name, "]需要您尽快审核")
		for _, v := range auditOrgBinds {
			OrgWorkTipMod.AppendTip(&OrgWorkTipMod.ArgsAppendTip{
				OrgID:     stepData.ID,
				OrgBindID: v,
				Msg:       auditMsg,
				System:    "erp_audit",
				BindID:    stepData.ID,
			})
		}
	}
	//下一个节点或当前节点的数据
	var nextStepChildData FieldsStepChild
	//自动进入下一个流程
	if autoNext {
		//检查下一个审批流程是否完成？
		nextKey := ""
		nowStepChildFinish := false
		isFinish := false
		finishStatus := 0
		//如果需要具体判断，则判断是否审批等情况
		if needCheckAuditOK {
			if CoreSQL.CheckTimeHaveData(stepChildData.AuditAt) || CoreSQL.CheckTimeHaveData(stepChildData.BanAt) {
				if CoreSQL.CheckTimeHaveData(stepChildData.AuditAt) {
					nowStepChildFinish = true
					nextKey = stepChildData.NextStepKey
				}
				if CoreSQL.CheckTimeHaveData(stepChildData.BanAt) {
					nextKey = stepChildData.BanNextStepKey
				}
				if nextKey == "" {
					if CoreSQL.CheckTimeHaveData(stepChildData.AuditAt) {
						isFinish = true
						finishStatus = 1
					}
					if CoreSQL.CheckTimeHaveData(stepChildData.BanAt) {
						isFinish = true
						finishStatus = 2
					}
				}
			}
		} else {
			//抄送等模式，直接写入下一个节点即可
			nextKey = stepChildData.NextStepKey
			if nextKey == "" {
				nowStepChildFinish = true
				isFinish = true
				finishStatus = 1
			}
		}
		//修改当前子节点完成
		if nowStepChildFinish {
			//如果是需要检查具体审核情况的，则跳过，因为之前已经更新了审批情况
			if !needCheckAuditOK {
				_, err := CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE erp_audit_step_child SET audit_at = NOW() WHERE id = :id", map[string]interface{}{
					"id": stepChildData.ID,
				})
				if err != nil {
					CoreLog.Error(appendLog, "update step child finish, step id: ", stepData.ID, ", step child id: ", stepChildData.ID, ", err: ", err)
					return
				}
				deleteStepChildCache(stepChildData.ID)
			}
		}
		//覆盖节点内容
		if nextKey != "" {
			nextStepChildData = getStepChildByKey(stepData.ID, nextKey)
		}
		//更新节点
		if isFinish {
			//完成，更新节点基本审核信息
			_, err := CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE erp_audit_step SET update_at = NOW(), finish_at = NOW(), finish_status = :finish_status WHERE id = :id", map[string]interface{}{
				"id":            stepData.ID,
				"finish_status": finishStatus,
			})
			if err != nil {
				CoreLog.Error(appendLog, "update step finish, step id: ", stepData.ID, ", err: ", err)
				return
			}
			deleteStepCache(stepData.ID)
		} else {
			//没有完成，只更新节点位置
			_, err := CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE erp_audit_step SET update_at = NOW(), now_step_child_key = :now_step_child_key WHERE id = :id", map[string]interface{}{
				"id":                 stepData.ID,
				"now_step_child_key": nextKey,
			})
			if err != nil {
				CoreLog.Error(appendLog, "update step, step id: ", stepData.ID, ", next key: ", nextKey, ", err: ", err)
				return
			}
			deleteStepCache(stepData.ID)
		}
	} else {
		nextStepChildData = stepChildData
	}
	//如果下一轮节点存在，则更新过期处理
	if nextStepChildData.ID > 0 {
		//根据节点已经存在的过期时间，递减step创建时间计算出实际的过期秒
		expireSec := nextStepChildData.ExpireAt.Unix() - stepData.CreateAt.Unix()
		//强制修正过期时间，至少确保1分钟
		if expireSec < 1 {
			expireSec = 60
		}
		//更新过期时间
		expireAt := CoreFilter.GetNowTimeCarbon().AddSeconds(int(expireSec)).Time
		_, err := CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE erp_audit_step_child SET expire_at = :expire_at WHERE id = :id", map[string]interface{}{
			"id":        nextStepChildData.ID,
			"expire_at": expireAt,
		})
		if err != nil {
			CoreLog.Error(appendLog, "update step child expire, step child id: ", nextStepChildData.ID, ", new expire at: ", expireAt, ", err: ", err)
			return
		}
		deleteStepChildCache(nextStepChildData.ID)
		//发送通知
		BaseExpireTip.AppendTipNoErr(&BaseExpireTip.ArgsAppendTip{
			OrgID:      stepData.OrgID,
			UserID:     0,
			SystemMark: "erp_audit_child",
			BindID:     nextStepChildData.ID,
			Hash:       "",
			ExpireAt:   expireAt,
		})
	}
	//完成反馈
	return
}
