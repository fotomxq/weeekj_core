package ServiceHousekeeping

import (
	ClassComment "github.com/fotomxq/weeekj_core/v5/class/comment"
	"github.com/robfig/cron"
)

//家政服务系统模块设计
/**
本模块分为：
1. 家政任务分配模块，该系统类似于配送服务系统，但更独立化设计。信息结构上，核心以客户信息、服务内容为基准。
2. 服务模块，为客户提供可选择的服务内容，方便一键购选下单，系统指派上⻔服务。
3. 服务统计模块，对服务团队人员的相关数据进行统计，支持包括评分、效率、服务时⻓等内容统计。
4. 统计支持，支持对服务效果单独统计并作为考评依据。
5. 服务绩效模块，用于计算每次服务类型下抽成的比例金额，可用于团队的工资发放依据。统计周期采用月为单位结算，统计可指定任意时间范围查看。
*/

/**
家政服务借助商城商品实现服务展示和购买，订单下单后识别商品扩展参数"housekeeping"标记为"true"，当存在时，不会创建配送单，而是直接构建服务请求。
本模块仅实现服务请求结构，用于提供上门服务等内容的管理、统计支持。
*/

var (
	Comment = ClassComment.Comment{
		TableName:         "service_housekeeping_comment",
		UserMoreComment:   false,
		UserEditComment:   false,
		UserDeleteComment: false,
		OrgDeleteComment:  false,
		System:            "service_housekeeping",
	}
	//定时器
	runTimer        *cron.Cron
	runAnalysisLock = false
	//OpenSub 是否启动订阅
	OpenSub = false
)

func Init() {
	//订阅关系
	if OpenSub {
		//消息列队
		subNats()
	}
}
