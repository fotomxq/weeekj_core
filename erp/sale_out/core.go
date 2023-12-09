package ERPSaleOut

//销售出库单
/**
1. 订单构建后，将所有产品拆分为不同单据发送给本模块
2. 本模块用于记录每个产品的出库记录
3. 同时将核算原价、优惠抵扣金额、实际售价、产品供货商成本价、产品修正成本价、实际毛利
*/

var (
	//OpenSub 是否启动订阅
	OpenSub = false
)

func Init() {
	if OpenSub {
		//消息列队
		subNats()
	}
}

func checkSaleType(saleType int) bool {
	switch saleType {
	case 0:
	case 1:
	case 2:
	default:
		return false
	}
	return true
}
