package UserLogin

import (
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

func runSave() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("user login save run, ", r)
		}
	}()
	_, _ = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "user_login_save", "expire_at < NOW()", nil)
}
