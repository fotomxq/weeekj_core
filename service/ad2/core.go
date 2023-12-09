package ServiceAD2

import AnalysisAny2 "gitee.com/weeekj/weeekj_core/v5/analysis/any2"

//第二代广告模块
/**
1. 没有分区的概念，简单设置即可使用
2. 投放区域直接约定为不同的mark，mark全系统交给前端固定
*/

var (
	//OpenAnalysis 是否启动通用框架体系的统计支持
	OpenAnalysis = false
)

func Init() {
	if OpenAnalysis {
		AnalysisAny2.SetConfigBeforeNoErr("service_ad2_ad_put", 0, 30)
		AnalysisAny2.SetConfigBeforeNoErr("service_ad2_ad_click", 0, 30)
	}
}
