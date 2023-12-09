package TMSUserRunning

import CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"

// pushNatsStatusUpdate 通知跑腿状态变更
func pushNatsStatusUpdate(action string, id int64, des string) {
	CoreNats.PushDataNoErr("/tms/user_running/update", action, id, "", map[string]interface{}{
		"des": des,
	})
}
