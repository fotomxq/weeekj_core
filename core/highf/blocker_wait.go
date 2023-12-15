package CoreHighf

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	"sync"
	"time"
)

//BlockerWait 等待拦截器
/**
1. 本模块用于持久性拦截一些重复请求，减少服务器压力。
2. 可用于统计数据拦截，当统计请求收到后，用此模块拦截处理，避免重复触发统计。
*/
type BlockerWait struct {
	//拦截数据列
	// 相同的数据将被列入拦截列
	blockerList []blockerWaitData
	//列队锁定
	blockerLock sync.Mutex
	//等待时间
	// 支持1-300秒范围，超出范围将自动回归数据
	// 时间不能超过5分钟，避免主线程序线程耗尽
	// 如果没有初始化，默认按5秒计算
	WaitTime int
}

type blockerWaitData struct {
	//模块ID
	ModID int64
	//模块标识码
	ModMark string
	//过期时间
	ExpireAt int64
}

// Init 初始化
func (t *BlockerWait) Init(waitTime int) {
	t.WaitTime = waitTime
	if t.WaitTime < 1 {
		t.WaitTime = 1
	}
	if t.WaitTime > 300 {
		t.WaitTime = 300
	}
}

// Check 检查是否通行？
// isNewData 用于说明是新的数据，方便拦截器识别处理
// b 用于说明是否通行，如果通行则需执行后续业务逻辑，否则请勿执行
func (t *BlockerWait) Check(modID int64, modMark string) (isNewData bool, b bool) {
	if t.WaitTime < 1 {
		t.WaitTime = 5
	}
	t.blockerLock.Lock()
	defer t.blockerLock.Unlock()
	var newBlockerList []blockerWaitData
	//查询是否存在，并重组数据
	isFind := false
	for _, v := range t.blockerList {
		if v.ModID == modID && v.ModMark == modMark {
			isFind = true
			if v.ExpireAt <= CoreFilter.GetNowTime().Unix() {
				b = true
				continue
			} else {
				newBlockerList = append(newBlockerList, v)
			}
		} else {
			//检查是否超出1天，则去掉遗留数据
			if v.ExpireAt <= CoreFilter.GetNowTimeCarbon().SubDay().Time.Unix() {
				continue
			}
			//其他数据加入数据集合
			newBlockerList = append(newBlockerList, v)
		}
	}
	//覆盖数据
	if isFind {
		t.blockerList = newBlockerList
	} else {
		//没有发现则添加拦截记录
		t.blockerList = append(t.blockerList, blockerWaitData{
			ModID:    modID,
			ModMark:  modMark,
			ExpireAt: CoreFilter.GetNowTimeCarbon().AddSeconds(t.WaitTime).Time.Unix(),
		})
		isNewData = true
	}
	//反馈
	return
}

// CheckWait 阻塞检查程序
// 将程序永久化阻塞，满足条件时释放
// 释放如果b=false则说明已经存在数据，需跳出，避免线程过多
func (t *BlockerWait) CheckWait(modID int64, modMark string, handle func(modID int64, modMark string)) {
	//第一次检查拦截器
	isNewData, b := t.Check(modID, modMark)
	if b {
		//通行释放
		handle(modID, modMark)
		return
	} else {
		//如果已经存在数据，则退出，因为有其他拦截器可能正在运行
		if !isNewData {
			return
		}
	}
	//持久化执行
	for {
		//检查拦截器
		_, b = t.Check(modID, modMark)
		if b {
			//通行释放
			handle(modID, modMark)
			return
		}
		//等待1秒继续
		time.Sleep(time.Second * 1)
	}
}
