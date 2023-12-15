package CoreMQTTSimple

import (
	"testing"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
)

var (
	testServerURL  = ""
	testUsername   = ""
	testPassword   = ""
	newMqttConnect = MQTTSimple{}
)

func TestInit(t *testing.T) {
	ToolsTest.Init(t)
	if err := newMqttConnect.Init(testServerURL, testUsername, testPassword, ""); err != nil {
		t.Error("cannot init, ", err)
		t.Fail()
	}
}

func TestSetDefaultClientID(t *testing.T) {
	newMqttConnect.SetDefaultClientID("default")
}

func TestSetServer(t *testing.T) {
	newMqttConnect.SetServer(testServerURL, testUsername, testPassword)
}

func TestConnectServer(t *testing.T) {
	opts := mqtt.ClientOptions{}
	opts.SetAutoReconnect(true)
	if token := newMqttConnect.ConnectServer(&opts); token.Error() != nil {
		t.Error(token.Error())
	}
}

func TestConnectClient(t *testing.T) {
	if token := newMqttConnect.ConnectClient(); token.Error() != "" {
		t.Error(token.Error())
	}
}

func TestPush(t *testing.T) {
	if err := newMqttConnect.Publish("top1/top2", 0, false, []byte("hello1")); err != nil {
		t.Error(err)
	}
}

func TestPushString(t *testing.T) {
	if err := newMqttConnect.Publish("top2/top3", 0, false, []byte("hello1")); err != nil {
		t.Error(err)
	}
}

// 案例
func TestEg(t *testing.T) {
	//建立mqtt连接
	opts := mqtt.ClientOptions{}
	opts.SetAutoReconnect(true)
	opts.SetDefaultPublishHandler(func(client mqtt.Client, message mqtt.Message) {
		t.Log("message default.")
	})
	opts.SetWill("mirror/line", "down line.", 1, true)
	opts.OnConnect = func(client mqtt.Client) {
		if token := client.Subscribe("mirror/phone", 0, func(client2 mqtt.Client, message mqtt.Message) {
			//解析数据
			t.Log("new message, ", message.Payload())
		}); token.Wait() && token.Error() != nil {
			t.Log("message mirror phone error.")
		}
	}
	_ = newMqttConnect.ConnectServer(&opts)
	//接收测试
	newMqttConnect.Subscribe("mirror/line", 0, func(client mqtt.Client, message mqtt.Message) {
		t.Log("message : mirror line.")
	})
}
