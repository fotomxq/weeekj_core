package CoreMQTTSimple

//mqtt物联网
// 该模块用于提供物联网专用的mqtt中间件服务
// 可通过该模块，实现标准的mqtt动作或组合动作
import (
	"errors"
	"fmt"
	"time"

	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTSimple struct {
	//全局client连接
	Client mqtt.Client
	//是否启动info追踪
	OpenInfoLog bool
	//默认特征指纹
	defaultClientID string
	//服务器地址\用户密码
	serverURL, username, password string
	//是否经过连接
	isConnect bool
	//前缀部分
	prefix string
	//选项扩展
	Options *mqtt.ClientOptions
	//必须断开重连
	NeedDisConnectOnIsConnect bool
}

// Init 初始化基本设计
// 本方法不会连接服务，需手动触发后续连接处理机制
func (t *MQTTSimple) Init(tServerURL, tUsername, tPassword, prefix string) (err error) {
	if tServerURL != t.serverURL || tUsername != t.username || tPassword != t.password {
		if t.Client != nil {
			t.Client.Disconnect(0)
		}
	}
	if t.Client != nil {
		if t.Client.IsConnected() {
			t.isConnect = true
			return nil
		}
	}
	t.prefix = prefix
	if t.defaultClientID == "" {
		t.defaultClientID, err = t.getClientID(t.prefix)
		if err != nil {
			return
		}
	}
	t.isConnect = false
	t.SetServer(tServerURL, tUsername, tPassword)
	//初始化ID指纹
	t.GetByIDOptions(t.defaultClientID)
	//默认启动日志
	t.OpenInfoLog = false
	return
}

// SetServer 设置服务器信息
func (t *MQTTSimple) SetServer(tServerURL, tUsername, tPassword string) {
	t.serverURL = tServerURL
	t.username = tUsername
	t.password = tPassword
}

// SetDefaultClientID 设置默认特征指纹
func (t *MQTTSimple) SetDefaultClientID(tDefaultClientID string) {
	t.defaultClientID = tDefaultClientID
}

// ConnectServer 连接到服务
func (t *MQTTSimple) ConnectServer(options *mqtt.ClientOptions) (token mqtt.Token) {
	t.Client = mqtt.NewClient(options)
	token = t.Client.Connect()
	return
}

// ConnectClient 链接到服务
func (t *MQTTSimple) ConnectClient() (err error) {
	if t.Client != nil {
		if t.Client.IsConnected() {
			if t.NeedDisConnectOnIsConnect {
				t.Client.Disconnect(10)
			} else {
				return nil
			}
		}
	}
	//自动重新连接
	t.Options.SetCleanSession(true)
	t.Options.SetAutoReconnect(true)
	t.Options.SetConnectTimeout(55 * time.Second)
	t.Options.SetKeepAlive(50 * time.Second)
	token := t.ConnectServer(t.Options)
	if token.Wait() && token.Error() != nil {
		t.isConnect = false
		err = errors.New("cannot connect client, " + token.Error().Error())
		return
	}
	t.isConnect = true
	return
}

// Publish 推送一个新的任务
// 异步进行
func (t *MQTTSimple) Publish(topic string, qos byte, retained bool, data []byte) (token mqtt.Token) {
	if !t.isConnect {
		return
	}
	token = t.Client.Publish(topic, qos, retained, data)
	t.logInfo("mqtt publish topic: ", topic)
	return
}

// PublishWait 推送一个新的任务
func (t *MQTTSimple) PublishWait(topic string, qos byte, retained bool, data []byte) (err error) {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint("mqtt publish topic error, ", r))
			return
		}
	}()
	//如果MQTT没有连接退出
	if !t.isConnect {
		return
	}
	//推送数据
	token := t.Client.Publish(topic, qos, retained, data)
	if token == nil {
		err = errors.New("mqtt token nil")
		return
	}
	//等待反馈
	token.Wait()
	if token.Error() != nil {
		err = errors.New(fmt.Sprint("mqtt push data, ", token.Error()))
		return
	}
	t.logInfo("mqtt publish topic: ", topic)
	return
}

// Subscribe 订阅消息
func (t *MQTTSimple) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) (token mqtt.Token) {
	token = t.Client.Subscribe(topic, qos, func(client mqtt.Client, message mqtt.Message) {
		message.Ack()
		t.logInfo("client send msg: ", topic, ", params: ", string(message.Payload()))
		go callback(client, message)
	})
	//t.logInfo("mqtt subscribe topic: ", topic)
	return
}

// SubscribeCancel 去掉订阅
func (t *MQTTSimple) SubscribeCancel(topic string) (token mqtt.Token) {
	token = t.Client.Unsubscribe(topic)
	return
}

// GetByIDOptions 构建带有id的选项
func (t *MQTTSimple) GetByIDOptions(id string) {
	t.Options = mqtt.NewClientOptions().AddBroker(t.serverURL).SetClientID(id)
	t.Options.SetUsername(t.username)
	t.Options.SetPassword(t.password)
}

// getClientID 生成带有前缀的随机clientID
func (t *MQTTSimple) getClientID(prefix string) (clientID string, err error) {
	var clientRand string
	clientRand, err = CoreFilter.GetRandStr3(5)
	if err != nil {
		return
	}
	clientID = fmt.Sprint(prefix, clientRand)
	return
}

// logInfo 推送日志
func (t *MQTTSimple) logInfo(args ...interface{}) {
	if !t.OpenInfoLog {
		return
	}
	CoreLog.Mqtt(args)
}
