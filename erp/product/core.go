package ERPProduct

import (
	ClassSort "github.com/fotomxq/weeekj_core/v5/class/sort"
	ClassTag "github.com/fotomxq/weeekj_core/v5/class/tag"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
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
	cacheBrandTime          = 1800
	cacheBrandBindTime      = 1800
	cacheTemplateTime       = 1800
	cacheTemplateBindTime   = 1800
	//OpenSub 订阅
	OpenSub = false
	//数据表
	brandDB        CoreSQL2.Client
	brandBindDB    CoreSQL2.Client
	templateDB     CoreSQL2.Client
	templateBindDB CoreSQL2.Client
	productValsDB  CoreSQL2.Client
)

// Init 初始化
func Init() {
	//初始化数据表
	brandDB.Init(&Router2SystemConfig.MainSQL, "erp_product_brand")
	brandBindDB.Init(&Router2SystemConfig.MainSQL, "erp_product_brand_bind")
	templateDB.Init(&Router2SystemConfig.MainSQL, "erp_product_template")
	templateBindDB.Init(&Router2SystemConfig.MainSQL, "erp_product_template_bind")
	productValsDB.Init(&Router2SystemConfig.MainSQL, "erp_product_product_vals")
	//nats
	if OpenSub {
		subNats()
	}
}
