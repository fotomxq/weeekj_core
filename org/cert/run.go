package OrgCert

import (
	"fmt"
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSystemClose "gitee.com/weeekj/weeekj_core/v5/core/system_close"
	"github.com/robfig/cron"
	"time"
)

// Run 维护服务
func Run() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("org cert run, ", r)
		}
	}()
	//等待
	time.Sleep(time.Second * 10)
	//启动时间
	runTimer = cron.New()
	if err := runTimer.AddFunc("@every 3s", func() {
		if runPayLock {
			return
		}
		runPayLock = true
		runPay()
		runPayLock = false
	}); err != nil {
		CoreLog.Error("org cert pay run, cron time, ", err)
	}
	orgCertWarningTime, err := BaseConfig.GetDataString("OrgCertWarningTime")
	if err != nil || orgCertWarningTime == "" {
		orgCertWarningTime = "6h"
	} else {
		if len(orgCertWarningTime) > 3 || len(orgCertWarningTime) < 1 {
			orgCertWarningTime = "6h"
		}
	}
	if err := runTimer.AddFunc(fmt.Sprint("@every ", orgCertWarningTime), func() {
		if runWarningCreateLock {
			return
		}
		runWarningCreateLock = true
		runWarningCreate()
		runWarningCreateLock = false
	}); err != nil {
		CoreLog.Error("org cert warning create run, cron time, ", err)
	}
	if err := runTimer.AddFunc("@every 5s", func() {
		if runAutoAudit {
			return
		}
		runAutoAudit = true
		runAudit()
		runAutoAudit = false
	}); err != nil {
		CoreLog.Error("org cert auto audit run, cron time, ", err)
	}
	runTimer.Start()
	//卡住进程
	CoreSystemClose.Wait()
	//退出时间
	runTimer.Stop()
}
