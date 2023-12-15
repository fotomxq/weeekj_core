package MapRoom

import (
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	IOTMQTT "github.com/fotomxq/weeekj_core/v5/iot/mqtt"
)

// 订阅方法集合
func subMQTT() {
	//呼叫请求变更
	if token := IOTMQTT.MQTTClient.Subscribe("map/room/service_status", 0, subRoomServiceStatus); token.Wait() && token.Error() != nil {
		CoreLog.Error("sub mqtt but token is error, " + token.Error().Error())
		return
	}
	//紧急呼叫请求
	if token := IOTMQTT.MQTTClient.Subscribe("map/room/emergency_call_need", 0, subRoomEmergencyCallNeed); token.Wait() && token.Error() != nil {
		CoreLog.Error("sub mqtt but token is error, " + token.Error().Error())
		return
	}
	if token := IOTMQTT.MQTTClient.Subscribe("map/room/emergency_call_un/org/#", 0, subRoomUnEmergencyCall); token.Wait() && token.Error() != nil {
		CoreLog.Error("sub mqtt but token is error, " + token.Error().Error())
		return
	}
}
