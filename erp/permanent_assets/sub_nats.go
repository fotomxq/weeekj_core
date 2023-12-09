package ERPPermanentAssets

import (
	BaseSystemMission "gitee.com/weeekj/weeekj_core/v5/base/system_mission"
)

func subNats() {
	//自动折旧定时任务
	BaseSystemMission.ReginSub(&runAutoExpireSysM, subNatsAutoExpire)
}
