package ServiceOrder

import (
	AnalysisAny2 "github.com/fotomxq/weeekj_core/v5/analysis/any2"
	"github.com/robfig/cron/v3"
	"sync"
)

//订单服务模块
// 外部模块禁止直接和本模块发生互动，如需互动请使用Mod子模块

var (
	//定时器
	runTimer       *cron.Cron
	runHistoryLock = false
	//等待订单锁定
	createWaitLock sync.Mutex
	//OpenSub 是否启动订阅
	OpenSub = false
	//OpenAnalysis 是否启动analysis
	OpenAnalysis = false
)

func Init() {
	if OpenAnalysis {
		//统计初始化
		AnalysisAny2.SetConfigBeforeNoErr("service_order_create_count", 1, 90)
		AnalysisAny2.SetConfigBeforeNoErr("service_order_finish_count", 1, 90)
		AnalysisAny2.SetConfigBeforeNoErr("service_order_finish_pay_price", 1, 90)
		AnalysisAny2.SetConfigBeforeNoErr("service_order_company_client_count", 1, 90)
		AnalysisAny2.SetConfigBeforeNoErr("service_order_company_client_price", 1, 90)
		AnalysisAny2.SetConfigBeforeNoErr("service_order_refund_create_count", 1, 90)
		AnalysisAny2.SetConfigBeforeNoErr("service_order_refund_finish_count", 1, 90)
		AnalysisAny2.SetConfigBeforeNoErr("service_order_refund_finish_price", 1, 90)
	}
	//消息订阅
	if OpenSub {
		//消息列队
		subNats()
	}
}

// 获取统计用的订单系统渠道
func getOrderSystemMarkKey(mark string) (orderSystemMarkKey int64) {
	switch mark {
	case "mall":
		orderSystemMarkKey = 0
	case "user_sub":
		orderSystemMarkKey = 1
	case "org_sub":
		orderSystemMarkKey = 2
	}
	return
}
