package ToolsShortURL

import (
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSystemClose "github.com/fotomxq/weeekj_core/v5/core/system_close"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/robfig/cron"
	"time"
)

func Run() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("tools short url run, ", r)
		}
	}()
	time.Sleep(time.Minute * 10)
	//启动定时器
	runTimer = cron.New()
	if err := runTimer.AddFunc("0 30 * * * *", func() {
		if runExpireLock {
			return
		}
		runExpireLock = true
		runExpire()
		runExpireLock = false
	}); err != nil {
		CoreLog.Error("base holiday season run, cron time, ", err)
	}
	runTimer.Start()
	//卡住进程
	CoreSystemClose.Wait()
	//退出时间
	runTimer.Stop()
}

func runExpire() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("tools short url run, ", r)
		}
	}()
	//归档处理
	if _, err := CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "tools_short_url", "expire_at < NOW()", map[string]interface{}{}); err != nil {
	}
}
