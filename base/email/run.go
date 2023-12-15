package BaseEmail

import (
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSystemClose "github.com/fotomxq/weeekj_core/v5/core/system_close"
	"github.com/robfig/cron"
)

// Run 自动维护工具
func Run() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("base email run error, ", r)
		}
	}()
	//设置阻断器
	runBlocker.SetExpire(15)
	//启动时间
	runTimer = cron.New()
	if err := runTimer.AddFunc("@every 1s", func() {
		if runSendEmailLock {
			return
		}
		runSendEmailLock = true
		runSendEmail()
		runSendEmailLock = false
	}); err != nil {
		CoreLog.Error("base email core run, cron time, ", err)
	}
	runTimer.Start()
	//卡住进程
	CoreSystemClose.Wait()
	//退出时间
	runTimer.Stop()
}
