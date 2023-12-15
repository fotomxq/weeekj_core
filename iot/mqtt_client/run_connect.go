package IOTMQTTClient

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
)

func runConnect() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("iot device run error, ", r)
		}
	}()
	mqttServerURL, err := BaseConfig.GetDataString("MQTTServerURLTCP")
	//绕过url为空的时候，如测试或特殊业务场景，不需要mqtt服务的
	if err != nil {
		CoreLog.Error("iot mqtt, get mqtt config server url, ", err)
		return
	}
	if mqttServerURL == "" {
		//CoreLog.Error("iot mqtt, get mqtt config server url, ", err)
		return
	}
	mqttServerUsername, err := BaseConfig.GetDataString("MQTTServerUsername")
	if err != nil {
		CoreLog.Error("iot mqtt, get mqtt config server username, ", err)
		return
	}
	mqttServerPassword, err := BaseConfig.GetDataString("MQTTServerPassword")
	if err != nil {
		CoreLog.Error("iot mqtt, get mqtt config server password, ", err)
		return
	}
	mqttClient.NeedDisConnectOnIsConnect = false
	if err := mqttClient.Init(mqttServerURL, mqttServerUsername, mqttServerPassword, mqttPrefix); err != nil {
		CoreLog.Error("iot mqtt, connect server, ", err)
		return
	}
	mqttClient.Options.SetOnConnectHandler(func(client mqtt.Client) {
		if err := initSub(); err != nil {
			CoreLog.Error("iot mqtt, connect init sub, ", err)
		}
	})
	if err = mqttClient.ConnectClient(); err != nil {
		CoreLog.Error("iot mqtt, connect failed, ", err)
		return
	}
}
