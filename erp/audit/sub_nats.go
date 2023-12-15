package ERPAudit

import (
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/nats-io/nats.go"
)

func subNats() {
	//节点过期处理
	CoreNats.SubDataByteNoErr("/base/expire_tip/expire", subNatsStepChildExpire)
}

// subNatsStepChildExpire 节点过期处理
func subNatsStepChildExpire(_ *nats.Msg, action string, stepChildID int64, _ string, _ []byte) {
	if action != "erp_audit_child" {
		return
	}
	//日志
	appendLog := "erp audit sub nats step child expire, "
	//获取节点内容
	stepChildData := getStepChildByID(stepChildID)
	if stepChildData.ID < 1 {
		return
	}
	if CoreSQL.CheckTimeHaveData(stepChildData.AuditAt) || CoreSQL.CheckTimeHaveData(stepChildData.BanAt) {
		return
	}
	//更新节点拒绝审核
	_, err := CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE erp_audit_step_child SET ban_at = NOW() WHERE id = :id", map[string]interface{}{
		"id": stepChildData.ID,
	})
	if err != nil {
		CoreLog.Error(appendLog, "update step child expire, step child id: ", stepChildData.ID, ", err: ", err)
		return
	}
	//清理缓冲
	deleteStepChildCache(stepChildData.ID)
	//通知和收尾当前节点
	sendStepAudit(stepChildData.StepID)
	//反馈
	return
}
