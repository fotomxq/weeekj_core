package BaseFileSys

import (
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

func runExpire() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("file run, ", r)
		}
	}()
	//检查是否存在即将过期的数据
	var count int64
	err := Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM core_file_claim WHERE expire_at > to_timestamp(1000000)")
	if err != nil || count < 1 {
		return
	}
	//清理过期数据
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "core_file_claim", "expire_at <= NOW() AND expire_at > to_timestamp(1000000)", nil)
	if err != nil {
		//CoreLog.Error("file expire run, ", err)
	}
}
