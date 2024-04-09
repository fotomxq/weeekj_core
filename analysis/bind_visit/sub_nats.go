package AnalysisBindVisit

import (
	BaseService "github.com/fotomxq/weeekj_core/v5/base/service"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/nats-io/nats.go"
	"github.com/tidwall/gjson"
)

func subNats() {
	_ = BaseService.SetService(&BaseService.ArgsSetService{
		ExpireAt:     CoreFilter.GetNowTimeCarbon().AddDay().Time,
		Name:         "组织绑定访问统计",
		Description:  "统计组织绑定访问情况",
		EventSubType: "sub",
		Code:         "analysis_bind_visit",
		EventType:    "nats",
		EventURL:     "/analysis/org/bind",
		EventParams:  "<<action>>:[new]:预设添加动作;<<data>>:json:{'userID':{'val_default':0,'val_type':'int64','val_enum':[],'val_desc':'用户ID','val_mod':'user_id_select'},'bindSystem':{'val_default':'','val_type':'string','val_enum':[],'val_desc':'绑定模块标识码','val_mod':''},'bindID':{'val_default':0,'val_type':'int64','val_enum':[],'val_desc':'绑定模块ID','val_mod':''}}",
	})
	CoreNats.SubDataByteNoErr("analysis_bind_visit", "/analysis/org/bind", subNatsNewBind)
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
