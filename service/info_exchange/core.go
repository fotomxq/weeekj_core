package ServiceInfoExchange

import (
	AnalysisAny2 "github.com/fotomxq/weeekj_core/v5/analysis/any2"
	ClassSort "github.com/fotomxq/weeekj_core/v5/class/sort"
	ClassTag "github.com/fotomxq/weeekj_core/v5/class/tag"
	"sync"
)

//信息交互模块
/**
1. 商户或用户均可以发布该信息
2. 提供基本的信息展示部分
3. 提供酬金方式展示，沟通确认后线上进行付款，平台担保
*/

var (
	//Sort 分类
	Sort = ClassSort.Sort{
		SortTableName: "service_info_exchange_sort",
	}
	//Tags 标签
	Tags = ClassTag.Tag{
		TagTableName: "service_info_exchange_tag",
	}
	//OpenSub 是否启动订阅
	OpenSub = false
	//OpenAnalysis 是否启动analysis
	OpenAnalysis = false
	//确认订单锁定机制
	orderFinishLock sync.Mutex
)

// Init 初始化
func Init() {
	//初始化统计
	if OpenAnalysis {
		AnalysisAny2.SetConfigBeforeNoErr("service_info_exchange_order_price", 1, 30)
		AnalysisAny2.SetConfigBeforeNoErr("service_info_exchange_order_count", 1, 30)
	}
	//消息订阅
	if OpenSub {
		//消息列队
		subNats()
	}
}
