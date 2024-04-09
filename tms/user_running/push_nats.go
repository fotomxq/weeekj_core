package TMSUserRunning

import CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"

// pushNatsStatusUpdate 通知跑腿状态变更
func pushNatsStatusUpdate(action string, id int64, des string) {
	CoreNats.PushDataNoErr("tms_user_running_update", "/tms/user_running/update", action, id, "", map[string]interface{}{
		"des": des,
	})
}
