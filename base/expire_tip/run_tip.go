package BaseExpireTip

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"time"
)

func runTip() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("base expire tip run, ", r)
		}
	}()
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
			CoreNats.PushDataNoErr("/base/expire_tip/expire", v.SystemMark, v.BindID, v.Hash, DataGetExpireData{
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
