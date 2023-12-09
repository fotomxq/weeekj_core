package BaseTempFile

import (
	CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"
)

func subNats() {
	//请求删除过期的文件
	CoreNats.SubDataByteNoErr("/base/expire_tip/expire", subNatsExpireID)
}
