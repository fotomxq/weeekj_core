package BaseExpireTip

import (
	"github.com/robfig/cron/v3"
	"sync"
)

//过期通知模块
/**
允许直接写入过期序列，到达时间后该模块会发出nats通知，广播某个内容已经过期
1. 可以杜绝外部其他模块再去申明单独的run，用于过期处理
2. 减少系统负荷
3. 本模块自带了缓冲、持久化处理，不需要关心额外的处理
4. 注意，本模块适用于中小规模应用，超大密集应用请勿使用，会存在毫秒、秒级别阻塞。如果对该阻塞不敏感，可考虑使用。
*/
var (
	//定时器
	runTimer          *cron.Cron
	runTipLock        = false
	runLoadExpireLock = false
	//OpenSub 是否启动订阅
	OpenSub = false
	//最近1小时即将过期数据集合
	// 内部缓冲
	waitExpire1HourList []FieldsTip
	waitExpire1HourLock sync.Mutex
)

func Init() {
	if OpenSub {
		subNats()
	}
}
