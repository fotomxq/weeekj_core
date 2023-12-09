package BaseExpireTip

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

func runLoadExpire() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("base expire tip run, ", r)
		}
	}()
	//加载最近1小时的过期数据
	var dataList []FieldsTip
	if err := Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, hash, expire_at FROM core_expire_tip WHERE expire_at <= $1", CoreFilter.GetNowTimeCarbon().AddHour().Time); err == nil {
		for _, v := range dataList {
			appendHaveNewData(v.ID, v.Hash, v.ExpireAt)
		}
	}
}
