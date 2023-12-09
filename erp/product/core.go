package ERPProduct

import (
	ClassSort "gitee.com/weeekj/weeekj_core/v5/class/sort"
	ClassTag "gitee.com/weeekj/weeekj_core/v5/class/tag"
)

//ERP产品信息库

var (
	//Sort 分类
	Sort = ClassSort.Sort{
		SortTableName: "erp_product_sort",
	}
	//Tags 标签
	Tags = ClassTag.Tag{
		TagTableName: "erp_product_tags",
	}
	//缓冲时间
	cacheProductTime        = 1800
	cacheProductCompanyTime = 1800
	//OpenSub 订阅
	OpenSub = false
)

// Init 初始化
func Init() {
	//nats
	if OpenSub {
		subNats()
	}
}
