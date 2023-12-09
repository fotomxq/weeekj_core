package OrgCert

import "github.com/robfig/cron"

//商户证件集合
// 可配置和管理子商户、用户、其他外围模块的证件集合信息

var(
	//定时器
	runTimer  *cron.Cron
	runPayLock = false
	runWarningCreateLock = false
	runAutoAudit = false
)