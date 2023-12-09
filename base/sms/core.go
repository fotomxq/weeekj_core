package BaseSMS

import (
	"github.com/robfig/cron"
)

var (
	//定时器
	runTimer      *cron.Cron
	runSendLock   = false
	runExpireLock = false
)

//TODO: 改进短信服务，将fields改为data且记录到redis存储，而不是数据库
/**
1. 将fields废弃，改为data记录数据
2. 记录到redis而不是数据库存储数据，每条数据最大保留30天
3. 区分验证类短信和营销类、通知类短信
4. 取消run机制，改为消息列队机制
*/
