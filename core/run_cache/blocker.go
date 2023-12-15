package CoreRunCache

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	"sync"
)

// Blocker 阻断器
type Blocker struct {
	//数据编辑次数
	EditCount int
	//锁定机制
	Lock sync.Mutex
	//上次通行时间
	LastPassTimeUnix int64
	//阻断间隔时间
	ExpireSec int64
}

// SetExpire 设置阻断器过期时间
func (t *Blocker) SetExpire(sec int64) {
	t.ExpireSec = sec
}

// NewEdit 编辑一次数据
func (t *Blocker) NewEdit() {
	t.Lock.Lock()
	t.EditCount += 1
	t.Lock.Unlock()
}

// CheckPass 检查阻断次数
func (t *Blocker) CheckPass() bool {
	//重归t.ExpireSec
	if t.ExpireSec < 1 {
		t.ExpireSec = 3
	}
	//如果时间超出ExpireSec，则直接通行
	nowTime := CoreFilter.GetNowTime().Unix()
	if nowTime-t.LastPassTimeUnix > t.ExpireSec {
		t.LastPassTimeUnix = nowTime
		return true
	}
	//检查阻断器
	if t.EditCount < 1 {
		return false
	}
	t.LastPassTimeUnix = nowTime
	t.Lock.Lock()
	t.EditCount = 0
	t.Lock.Unlock()
	return true
}
