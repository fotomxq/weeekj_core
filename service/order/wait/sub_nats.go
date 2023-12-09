package ServiceOrderWait

import (
	CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/nats-io/nats.go"
)

func subNats() {
	//等待订单过期处理
	CoreNats.SubDataByteNoErr("/base/expire_tip/expire", subNatsOrderExpire)
}

// 等待订单过期处理
func subNatsOrderExpire(_ *nats.Msg, action string, id int64, _ string, _ []byte) {
	//如果系统不符合，跳出
	if action != "service_order_wait" {
		return
	}
	//获取订单数据
	_, _ = CoreSQL.DeleteOne(Router2SystemConfig.MainDB.DB, "service_order_wait", "id", map[string]interface{}{
		"id": id,
	})
	//清理缓冲
	deleteWaitCache(id)
}
