package ERPPermanentAssets

import (
	BaseSystemMission "github.com/fotomxq/weeekj_core/v5/base/system_mission"
	ClassSort "github.com/fotomxq/weeekj_core/v5/class/sort"
	ClassTag "github.com/fotomxq/weeekj_core/v5/class/tag"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

//固定资产
/**
1. 记录固定资产数据
2. 可以定期盘点，设置盘点时间
*/

var (
	//Sort 分类
	Sort = ClassSort.Sort{
		SortTableName: "erp_permanent_assets_product_sort",
	}
	//Tags 标签
	Tags = ClassTag.Tag{
		TagTableName: "erp_permanent_assets_product_tag",
	}
	//数据库
	productSQL CoreSQL2.Client
	//调度任务
	runAutoExpireSysM = BaseSystemMission.Mission{
		OrgID:    0,
		Name:     "ERP固定资产自动折旧",
		Mark:     "erp.run.permanent_assets",
		NextTime: "每天05:00",
		Bind: BaseSystemMission.MissionBind{
			NatsMsg: "/erp/permanent_assets/run_auto_expire",
		},
	}
	//OpenSub 订阅
	OpenSub = false
)

// Init 初始化
func Init() {
	//初始化数据库
	productSQL.Init(&Router2SystemConfig.MainSQL, "erp_permanent_assets_product")
	//nats
	if OpenSub {
		subNats()
		BaseSystemMission.ReginWait(&runAutoExpireSysM, time.Time{})
	}
}
