package BaseTempFile

import (
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
)

func subNats() {
	//请求删除过期的文件
	CoreNats.SubDataByteNoErr("base_expire_tip_expire", "/base/expire_tip/expire", subNatsExpireID)
}
