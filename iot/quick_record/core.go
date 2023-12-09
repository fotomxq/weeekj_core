package IOTQuickRecord

import (
	"github.com/robfig/cron"
)

//设备快速录入设计
// 本模块用于设备主动发起邀约，服务端确认后快速构建设备组和密码，全程确保无手工录入过程


var(
	//定时器
	runTimer *cron.Cron
	runLock = false
)