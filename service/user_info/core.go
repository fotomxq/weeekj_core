package ServiceUserInfo

import (
	AnalysisAny2 "github.com/fotomxq/weeekj_core/v5/analysis/any2"
	ClassSort "github.com/fotomxq/weeekj_core/v5/class/sort"
	ClassTag "github.com/fotomxq/weeekj_core/v5/class/tag"
	CoreHighf "github.com/fotomxq/weeekj_core/v5/core/highf"
)

//信息服务模块
// 关联用户或非关联用户，构建信息汇总和处理中心
// 记录如该对象的姓名、年龄等

var (
	//Sort 人员分类
	Sort = ClassSort.Sort{
		SortTableName: "service_user_info_sort",
	}
	//Tag 人员标签
	Tag = ClassTag.Tag{
		TagTableName: "service_user_info_tags",
	}
	//DocSort 文档分类
	DocSort = ClassSort.Sort{
		SortTableName: "service_user_info_doc_sort",
	}
	//DocTag 文档标签
	DocTag = ClassTag.Tag{
		TagTableName: "service_user_info_doc_tags",
	}
	//TemplateSort 模版分类
	TemplateSort = ClassSort.Sort{
		SortTableName: "service_user_info_template_sort",
	}
	//TemplateTag 模版标签
	TemplateTag = ClassTag.Tag{
		TagTableName: "service_user_info_template_tags",
	}
	//缓存时间
	cacheTime = 604800
	//OpenSub 是否启动订阅
	OpenSub = false
	//OpenAnalysis 是否启动analysis
	OpenAnalysis = false
	//统计数据的拦截器
	analysisBlockerWait CoreHighf.BlockerWait
)

func Init() {
	//初始化统计混合模块
	if OpenAnalysis {
		subAnalysis()
	}
	if OpenSub {
		analysisBlockerWait.Init(5)
		subNats()
	}
}

func subAnalysis() {
	AnalysisAny2.SetConfigBeforeNoErr("service_user_info_all_count", 3, 365)
	AnalysisAny2.SetConfigBeforeNoErr("service_user_info_count", 3, 365)
	AnalysisAny2.SetConfigBeforeNoErr("service_user_info_month_count", 3, 365)
	AnalysisAny2.SetConfigBeforeNoErr("service_user_info_die_count", 3, 365)
	AnalysisAny2.SetConfigBeforeNoErr("service_user_info_out_count", 3, 365)
	AnalysisAny2.SetConfigBeforeNoErr("service_user_info_die_month_count", 3, 365)
	AnalysisAny2.SetConfigBeforeNoErr("service_user_info_out_month_count", 3, 365)
	AnalysisAny2.SetConfigBeforeNoErr("service_user_info_old_75_count", 3, 365)
	AnalysisAny2.SetConfigBeforeNoErr("service_user_info_gender_count", 3, 365)
}
