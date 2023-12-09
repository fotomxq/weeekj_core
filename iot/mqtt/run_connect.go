package IOTMQTT

import (
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

func runConnect() (b bool) {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("iot device run error, ", r)
		}
	}()
	mqttOpen := BaseConfig.GetDataBoolNoErr("MQTTOpen")
	if !mqttOpen {
		return
	}
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
	MQTTClient.NeedDisConnectOnIsConnect = false
	if err := MQTTClient.Init(mqttServerURL, mqttServerUsername, mqttServerPassword, mqttPrefix); err != nil {
		if connectTryCount < 1 {
		} else if connectTryCount < 5 && connectTryCount >= 1 {
			time.Sleep(time.Second * 3)
		} else if connectTryCount < 7 && connectTryCount >= 5 {
			time.Sleep(time.Second * 5)
		} else if connectTryCount < 15 && connectTryCount >= 7 {
			time.Sleep(time.Second * 15)
		} else if connectTryCount < 20 && connectTryCount >= 15 {
			time.Sleep(time.Second * 30)
		} else {
			time.Sleep(time.Second * 60)
		}
		connectTryCount += 1
		CoreLog.Error("iot mqtt, connect server, try count: ", connectTryCount, ", ", err)
		runConnectLock = false
		return
	}
	MQTTClient.Options.SetOnConnectHandler(func(client mqtt.Client) {
		connectTryCount = 0
		if OpenBaseMission {
			if err := initSub(); err != nil {
				CoreLog.Error("iot mqtt, connect init sub, ", err)
			}
		}
		for k, _ := range subFunc {
			subFunc[k]()
		}
		MQTTIsConnect = true
		runConnectLock = false
	})
	if err = MQTTClient.ConnectClient(); err != nil {
		if connectTryCount < 1 {
		} else if connectTryCount < 5 && connectTryCount >= 1 {
			time.Sleep(time.Second * 3)
		} else if connectTryCount < 7 && connectTryCount >= 5 {
			time.Sleep(time.Second * 30)
		} else if connectTryCount < 15 && connectTryCount >= 7 {
			time.Sleep(time.Minute * 5)
		} else if connectTryCount < 20 && connectTryCount >= 15 {
			time.Sleep(time.Minute * 30)
		} else if connectTryCount < 50 && connectTryCount >= 20 {
			time.Sleep(time.Hour * 1)
		} else if connectTryCount < 100 && connectTryCount >= 50 {
			time.Sleep(time.Hour * 6)
		} else {
			time.Sleep(time.Hour * 12)
		}
		connectTryCount += 1
		CoreLog.Error("iot mqtt, connect failed, try count: ", connectTryCount, ", ", err)
		runConnectLock = false
		return
	}
	return true
}
