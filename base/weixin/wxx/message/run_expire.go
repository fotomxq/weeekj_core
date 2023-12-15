package BaseWeixinWXXMessage

import (
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

func runExpire() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("weixin message template run, ", r)
		}
	}()
	//删除过期数据
	//删除失败或没有数据，不做记录
	_, _ = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "core_weixin_wxx_message", "expire_at < NOW()", nil)
}
