package CoreSQL2

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	"github.com/golang-module/carbon"
	"reflect"
	"sync"
	"time"
)

var (
	//WaitLog 为减少和外部模块的侵入性，此设计将记录一些必要的元素内容，外部抽取后释放即可
	// 等待处理的请求记录
	WaitLog        []WaitLogType
	WaitLogLock    sync.Mutex
	waitLogNew     []WaitLogType
	waitLogNewLock sync.Mutex
)

type WaitLogType struct {
	//动作类型
	// select/get/insert/update/delete/analysis
	Action string
	//消息内容
	Msg string
	//是否为事务关系
	IsBegin bool
	//开始时间
	StartAt time.Time
	//结束时间
	EndAt time.Time
	//执行时间
	RunSec int64
	//反馈尺寸 单位字节
	ResultSize int64
	//是否存在报错
	Err error
}

// appendLog 记录日志
func appendLog(action string, msg string, isBegin bool, startAt carbon.Carbon, result any, err error) {
	endAt := CoreFilter.GetNowTimeCarbon()
	waitLogNewLock.Lock()
	waitLogNew = append(waitLogNew, WaitLogType{
		Action:     action,
		Msg:        msg,
		IsBegin:    isBegin,
		StartAt:    startAt.Time,
		EndAt:      endAt.Time,
		RunSec:     endAt.DiffInSecondsWithAbs(startAt),
		ResultSize: int64(reflect.TypeOf(reflect.ValueOf(result)).Size()),
		Err:        err,
	})
	if len(waitLogNew) > 10 && len(WaitLog) < 1 {
		WaitLogLock.Lock()
		copy(WaitLog, waitLogNew)
		waitLogNew = []WaitLogType{}
		WaitLogLock.Unlock()
	}
	waitLogNewLock.Unlock()
}
