package MarketCore

import (
	ClassSort "github.com/fotomxq/weeekj_core/v5/class/sort"
	ClassTag "github.com/fotomxq/weeekj_core/v5/class/tag"
)

//营销服务模块
// 支持营销人员管理、营销记录、排名处理、提成计算管理、人员和客户关系处理

// 营销成员采用组织成员分组进行协调管理，将分组列指定为营销团队（商户设置）即可完成。
// 系统在该商户配置下，查询对应分组的所有成员，作为营销的人员进行管理协调。

var (
	//BindSort 成员分类
	BindSort = ClassSort.Sort{
		SortTableName: "market_core_sort",
	}
	//BindTag 成员标签
	BindTag = ClassTag.Tag{
		TagTableName: "market_core_tags",
	}
	//OpenSub 是否启动订阅
	OpenSub = false
)
