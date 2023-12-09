package BaseFileSys

import (
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"time"
)

// 定时删除超过时间的访问记录
func runVisit() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("file visit run, ", r)
		}
	}()
	//获取配置
	configData, err := BaseConfig.GetDataString("FileVisitExpireTime")
	if err != nil {
		CoreLog.Error("file visit run, get config by FileVisitExpireTime, ", err)
		configData = "-15120h"
	}
	//获取过期时间
	expireAt, err := CoreFilter.GetTimeByAdd(configData)
	if err != nil {
		CoreLog.Error("file visit run, get expire time, ", err)
		expireAt = CoreFilter.GetNowTime().Add(0 - time.Hour*15120)
	}
	//直接删除旧的数据
	if _, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "core_file_claim_visit", "create_at < :expire_at", map[string]interface{}{
		"expire_at": expireAt,
	}); err != nil {
		//CoreLog.Error("file visit run, delete expire data, ", err)
	}
}
