package UserChat

import (
	IOTMQTT "github.com/fotomxq/weeekj_core/v5/iot/mqtt"
	"sync"
)

//用户聊天工具模块
/**
1. 可以快速进行聊天处理，采用简化的MQTT和API轮训机制实现
2. 建立频道，完成双人、多人的聊天室功能
*/

var (
	//领取红包或票据锁定机制
	takeMoneyOrTicketLock sync.Mutex
)

func Init() {
	//初始化mqtt订阅
	IOTMQTT.AppendSubFunc(initSub)
}
