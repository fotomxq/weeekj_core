package BaseExpireTip

import (
	BaseService "github.com/fotomxq/weeekj_core/v5/base/service"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

func runTip() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("base expire tip run, ", r)
		}
	}()
	//推送服务初始化
	_ = BaseService.SetService(&BaseService.ArgsSetService{
		ExpireAt:     CoreFilter.GetNowTimeCarbon().AddDay().Time,
		Name:         "过期模块过期通知",
		Description:  "数据需要发生过期，立刻过期",
		EventSubType: "push",
		Code:         "base_expire_tip_expire",
		EventType:    "nats",
		EventURL:     "/base/expire_tip/expire",
		//TODO:待补充
		EventParams: "",
	})
	_ = BaseService.SetService(&BaseService.ArgsSetService{
		ExpireAt:     CoreFilter.GetNowTimeCarbon().AddDay().Time,
		Name:         "过期模块清理通知",
		Description:  "数据需要全局清理",
		EventSubType: "all",
		Code:         "base_expire_tip_expire_clear",
		EventType:    "nats",
		EventURL:     "/base/expire_tip/expire_clear",
		//TODO:待补充
		EventParams: "",
	})
	//持续执行
	for {
		//调用子模块进程
		runTipIn()
		//每秒执行1次
		time.Sleep(time.Second * 1)
	}
}

func runTipIn() {
	//当前时间
	nowAt := CoreFilter.GetNowTime()
	//锁定机制
	waitExpire1HourLock.Lock()
	defer waitExpire1HourLock.Unlock()
	//遍历数据，检查是否到期
	var newCacheList []FieldsTip
	for _, v := range waitExpire1HourList {
		if v.ExpireAt.Unix() <= nowAt.Unix() {
			//通知模块
			CoreNats.PushDataNoErr("base_expire_tip_expire", "/base/expire_tip/expire", v.SystemMark, v.BindID, v.Hash, DataGetExpireData{
				OrgID:    v.OrgID,
				UserID:   v.UserID,
				ExpireAt: v.ExpireAt,
			})
			Router2SystemConfig.MainCache.DeleteMark(getCacheMark(v.ID))
			if _, err := CoreSQL.DeleteOne(Router2SystemConfig.MainDB.DB, "core_expire_tip", "id", map[string]interface{}{
				"id": v.ID,
			}); err != nil {
				CoreLog.Error("base expire tip delete tip failed, id: ", v.ID, ", err: ", err)
				continue
			}
			continue
		}
		newCacheList = append(newCacheList, v)
	}
	waitExpire1HourList = newCacheList
}
