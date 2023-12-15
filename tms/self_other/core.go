package TMSSelfOther

import TMSSelfOtherKuai100 "github.com/fotomxq/weeekj_core/v5/tms/self_other/kuai100"

//第三方配送单查询聚合方案
/**
1. 本接口融合了多个快递查询平台，可快速查询对应快递单的物流信息
2. 对外暴露接口将聚合和数据清洗，统一数据结构体
*/

// Init 聚合初始化
func Init() (err error) {
	//快递100
	err = TMSSelfOtherKuai100.Init()
	if err != nil {
		return
	}
	return
}
