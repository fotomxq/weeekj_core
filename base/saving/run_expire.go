package BaseSaving

import (
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

func runExpire() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("user saving expire run, ", r)
		}
	}()
	//检查阻拦器
	if !runExpireBlocker.CheckPass() {
		return
	}
	//删除数据
	_, _ = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "core_saving", "expire_at < NOW()", nil)
}
