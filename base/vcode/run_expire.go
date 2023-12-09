package VCodeImageCore

import (
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"time"
)

func runExpire() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("vcode run error, ", r)
		}
	}()
	var err error
	//更新配置
	expiredTime, err = BaseConfig.GetDataInt64("VerificationCodeImageExpireTime")
	if err != nil {
		expiredTime = 120
	}
	intervalTime, err = BaseConfig.GetDataInt64("VerificationCodeImageIntervalTime")
	if err != nil {
		intervalTime = 1
	}
	//删除过期数据
	tExpiredTime := CoreFilter.GetNowTime().Add(time.Second * time.Duration(expiredTime))
	//不记录删除数据的错误信息
	_, _ = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "core_vcode_image", "create_at>=:create_at", map[string]interface{}{
		"create_at": tExpiredTime,
	})
}
