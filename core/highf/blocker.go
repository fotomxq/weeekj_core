package CoreHighf

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	"sync"
	"time"
)

// HighFBlocker nats拦截器
// 拦截器主要阻拦不必要的消费处理，如高频数据等，一些敏感数据请勿使用
type HighFBlocker struct {
	//间隔时间秒
	WaitSec int64
	//是否仅拦截相同数据
	OnlySameData bool
	//数据集合
	dataList []dataType
	//数据锁
	dataLock sync.Mutex
}

// 数据集合
type dataType struct {
	ID         int64
	Mark       string
	Data       string
	CreateTime int64
}

// Init 初始化
func (t *HighFBlocker) Init(waitSec int64, onlySameData bool) {
	t.WaitSec = waitSec
	t.OnlySameData = onlySameData
}

// Run 运行拦截器
func (t *HighFBlocker) Run() {
	for {
		nowTime := CoreFilter.GetNowTime().Unix()
		t.dataLock.Lock()
		var newDataList []dataType
		for _, v := range t.dataList {
			if v.CreateTime+t.WaitSec < nowTime {
				continue
			}
			newDataList = append(newDataList, v)
		}
		t.dataList = newDataList
		t.dataLock.Unlock()
		time.Sleep(time.Second * 10)
	}
}

// Check 检查是否是否可通行
func (t *HighFBlocker) Check(id int64, mark string, data string) bool {
	t.dataLock.Lock()
	defer t.dataLock.Unlock()
	nowTime := CoreFilter.GetNowTime().Unix()
	for _, v := range t.dataList {
		if v.ID == id && v.Mark == mark && v.CreateTime+t.WaitSec >= nowTime {
			if t.OnlySameData {
				if v.Data == data {
					return false
				}
			} else {
				return false
			}
			break
		}
	}
	t.dataList = append(t.dataList, dataType{
		ID:         id,
		Mark:       mark,
		Data:       data,
		CreateTime: nowTime,
	})
	return true
}
