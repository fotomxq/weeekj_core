package ERPProductMall

import (
	ClassSort "github.com/fotomxq/weeekj_core/v5/class/sort"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

//产品商城服务模块
/**
1. 用途企业内部的产品商城，展示公司的产品，供员工查看、申请采购需求
*/

var (
	//Sort 分类
	Sort = ClassSort.Sort{
		SortTableName: "erp_product_mall_sort",
	}
	//缓冲时间
	cacheProductMallTime = 1800
	//数据库句柄
	productMallDB CoreSQL2.Client
	//OpenSub 订阅
	OpenSub = false
)

func Init() {
	//初始化数据库
	productMallDB.Init(&Router2SystemConfig.MainSQL, "erp_product_mall")
	//nats
	if OpenSub {
		subNats()
	}
}
