package TMSTransport

import (
	"encoding/json"
	"fmt"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	IOTMQTT "github.com/fotomxq/weeekj_core/v5/iot/mqtt"
)

// MQTT订阅
func initSub() {
	//暂无订阅处理
}

// 推送给用户消息更新了
func pushMQTTTransportUpdate(orgID, bindID, transportID int64, action int) {
	//构建数据
	type newDataType struct {
		//成员ID
		BindID int64 `json:"bindID"`
		//配送单ID
		TransportID int64 `json:"transportID"`
		//行为特征
		// 0 新配送单; 1 配送单发生变更; 2 配送单被删除; 3 配送单被迁移走; 4 配送单被迁移入
		Action int `json:"action"`
	}
	newData := newDataType{
		BindID:      bindID,
		TransportID: transportID,
		Action:      action,
	}
	//获取和打包数据
	dataByte, err := json.Marshal(newData)
	if err != nil {
		CoreLog.Error("user chat push mqtt, push update by user, json err: ", err)
		return
	}
	topic := fmt.Sprint("tms/transport/org/", orgID)
	if err = IOTMQTT.MQTTClient.PublishWait(topic, 0, false, dataByte); err != nil {
		CoreLog.Warn("tms transport, push update: ", err)
		return
	}
}
