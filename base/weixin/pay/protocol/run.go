package BaseWeixinPayProtocol

import (
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSystemClose "github.com/fotomxq/weeekj_core/v5/core/system_close"
	"github.com/robfig/cron"
	"time"
)

func Run() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("base weixin pay protocol run, ", r)
		}
		runTimer.Stop()
	}()
	//等待
	time.Sleep(time.Second * 10)
	//日志归档
	runTimer = cron.New()
	if err := runTimer.AddFunc("0 0 1 * * *", func() {
		if runLock {
			return
		}
		runLock = true
		runNext()
		runLock = false
	}); err != nil {
		CoreLog.Error("base weixin pay protocol run, cron time, ", err)
	}
	runTimer.Start()
	//卡住进程
	CoreSystemClose.Wait()
	//退出时间
	runTimer.Stop()
}
