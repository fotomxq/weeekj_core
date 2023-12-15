package BaseQiniu

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
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
			CoreLog.Error("base qiniu run, ", r)
		}
	}()
	//延迟10秒启动
	time.Sleep(time.Second * 13)
	//启动定时器
	runTimer = cron.New()
	if err := runTimer.AddFunc("0 10 * * * *", func() {
		if runExpireLock {
			return
		}
		runExpireLock = true
		//删除过期的数据
		// 过期时间以创建后的1小时为准
		_, _ = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "core_file_qiniu_wait", "create_at < :expire", map[string]interface{}{
			"expire": CoreFilter.GetNowTimeCarbon().SubHours(1).Time,
		})
		runExpireLock = false
	}); err != nil {
		CoreLog.Error("base qiniu run, cron time, ", err)
	}
	runTimer.Start()
	//卡住进程
	CoreSystemClose.Wait()
	//退出时间
	runTimer.Stop()
}
