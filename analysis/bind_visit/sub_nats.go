package AnalysisBindVisit

import (
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/nats-io/nats.go"
	"github.com/tidwall/gjson"
)

func subNats() {
	CoreNats.SubDataByteNoErr("/analysis/org/bind", subNatsNewBind)
}

func subNatsNewBind(_ *nats.Msg, action string, _ int64, _ string, data []byte) {
	if action != "new" {
		return
	}
	//锁定机制
	appendLogLock.Lock()
	defer appendLogLock.Unlock()
	//获取参数
	userID := gjson.GetBytes(data, "userID").Int()
	bindSystem := gjson.GetBytes(data, "bindSystem").String()
	bindID := gjson.GetBytes(data, "bindID").Int()
	//检查是否访问过，避免重复添加
	if CheckLog(userID, bindSystem, bindID) {
		return
	}
	//添加记录
	_, err := CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO analysis_bind_visit(user_id, bind_system, bind_id) VALUES(:user_id, :bind_system, :bind_id)", map[string]interface{}{
		"user_id":     userID,
		"bind_system": bindSystem,
		"bind_id":     bindID,
	})
	CoreLog.Error("analysis bind visit sub nats new bind, ", err)
}
