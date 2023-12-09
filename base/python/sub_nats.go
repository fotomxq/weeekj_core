package BasePython

import (
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"
	"github.com/nats-io/nats.go"
)

func subNats() {
	//订阅数据反馈
	CoreNats.SubDataByteNoErr("/base/python/result", func(_ *nats.Msg, _ string, resultID int64, _ string, _ []byte) {
		updateResult(resultID)
	})
	//数据过期处理
	CoreNats.SubDataByteNoErr("/base/expire_tip/expire", func(_ *nats.Msg, action string, id int64, _ string, _ []byte) {
		if action != "core_python" {
			return
		}
		if err := deleteByID(id); err != nil {
			CoreLog.Error("sub nats core python sub expire, id: ", id, ", err: ", err)
		}
	})
}
