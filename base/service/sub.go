package BaseService

import (
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	"github.com/nats-io/nats.go"
	"sync"
	"time"
)

var (
	//阻断器
	blockSync sync.Mutex
	//堆叠列队
	blockList []waitSyncBlock
)

type waitSyncBlock struct {
	//服务ID
	ServiceID int64
	//追加时间
	// 倒计时机制
	AppendTime int64
	//是否正在等待
	IsWait bool
	//增加服务端发送消息次数
	SendCount int64 `db:"send_count" json:"sendCount" check:"intThan0"`
	//增加服务端接收次数
	ReceiveCount int64 `db:"receive_count" json:"receiveCount" check:"intThan0"`
}

func subNats() {
	//发生新的请求
	CoreNats.SubDataByteNoErr("base_service_request", "/base/service/request", subNatsRequest)
}

// action 服务code
// mark 订阅和推送类型: sub订阅; pub发布
func subNatsRequest(msg *nats.Msg, action string, _ int64, mark string, _ []byte) {
	//等待数据库连接
	if !WaitDBConnect {
		time.Sleep(time.Second * 60)
		WaitDBConnect = true
	}
	//初始化
	var sendCount, receiveCount int64
	sendCount = 0
	receiveCount = 0
	switch mark {
	case "sub":
		receiveCount += 1
	case "push":
		sendCount += 1
	default:
		return
	}
	//找到服务
	serviceData := getServiceByCode(action)
	if serviceData.ID < 1 {
		return
	}
	//更新统计
	var nowWaitData waitSyncBlock
	var nowWaitKey int
	needReturn := false
	blockSync.Lock()
	if len(blockList) < 1 {
		nowWaitData = waitSyncBlock{
			ServiceID:    serviceData.ID,
			AppendTime:   10,
			SendCount:    sendCount,
			ReceiveCount: receiveCount,
		}
		blockList = append(blockList, nowWaitData)
		nowWaitKey = 0
	} else {
		for k, v := range blockList {
			if v.ServiceID != serviceData.ID {
				continue
			}
			blockList[k].SendCount += sendCount
			blockList[k].ReceiveCount += receiveCount
			if !blockList[k].IsWait {
				blockList[k].IsWait = true
			} else {
				//说明已经有在序列的请求处理，直接跳出
				needReturn = true
			}
			nowWaitData = blockList[k]
			nowWaitKey = k
			break
		}
	}
	blockSync.Unlock()
	if needReturn {
		return
	}
	//进入时间阻塞
	for {
		if nowWaitData.AppendTime < 1 {
			break
		}
		nowWaitData.AppendTime -= 1
		time.Sleep(time.Second * 1)
	}
	//追加数据
	blockSync.Lock()
	err := appendAnalysisData(&argsAppendAnalysisData{
		ServiceID:    serviceData.ID,
		SendCount:    blockList[nowWaitKey].SendCount,
		ReceiveCount: blockList[nowWaitKey].ReceiveCount,
	})
	if err != nil {
		CoreLog.Error("base service subNatsRequest appendAnalysisData failed, ", err)
	}
	var newBlockList []waitSyncBlock
	for _, v := range blockList {
		if v.ServiceID == serviceData.ID {
			continue
		}
		newBlockList = append(newBlockList, v)
	}
	blockList = newBlockList
	blockSync.Unlock()
	//更新服务过期时间
	_ = updateServiceExpire(serviceData.ID)
}
