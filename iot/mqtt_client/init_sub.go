package IOTMQTTClient

import (
	"errors"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

//订阅主题
func initSub() (err error) {
	//设备在线情况更正
	if token := mqttClient.Subscribe("device/online", 0, func(client mqtt.Client, message mqtt.Message) {

	}); token.Wait() && token.Error() != nil {
		err = errors.New("sub mqtt but token is error, " + token.Error().Error())
		return
	}
	return
}