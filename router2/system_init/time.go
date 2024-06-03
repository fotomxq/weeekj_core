package Router2SystemInit

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	"time"
)

// RunStartPrint 启动时间处理器
type RunStartPrint struct {
	//前置语言
	Pre string
	//后置语言
	Suf string
	//开始时间
	startAt time.Time
	//结束时间
	endAt time.Time
}

func (t *RunStartPrint) Start() {
	t.startAt = CoreFilter.GetNowTime()
}

func (t *RunStartPrint) End() {
	t.endAt = CoreFilter.GetNowTime()
}

func (t *RunStartPrint) Print() {
	if t.Pre != "" {
		t.Pre = t.Pre + ", "
	}
	if t.Suf == "" {
		t.Suf = "."
	} else {
		t.Suf = ", " + t.Suf
	}
	CoreLog.Info(t.Pre, t.endAt.Unix()-t.startAt.Unix(), "s/", t.endAt.UnixMilli()-t.startAt.UnixMilli(), "ms", t.Suf)
}

func (t *RunStartPrint) EndPrint() {
	t.End()
	t.Print()
}
