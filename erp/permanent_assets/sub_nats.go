package ERPPermanentAssets

import (
	BaseSystemMission "github.com/fotomxq/weeekj_core/v5/base/system_mission"
)

func subNats() {
	//自动折旧定时任务
	BaseSystemMission.ReginSub(&runAutoExpireSysM, subNatsAutoExpire)
}
