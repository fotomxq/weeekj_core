package CoreNats

import (
	"encoding/json"
	"errors"
	"fmt"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	"github.com/nats-io/nats.go"
	"time"
)

var (
	//是否启动nats
	openNC = false
	//NatsURL 链接的地址
	serverURL = ""
	//连接句柄
	nc *nats.Conn
	//subPrefix 全局订阅地址前缀
	subPrefix = ""
)

func Init(url string) (err error) {
	//赋予值
	serverURL = url
	//链接服务
	nc, err = nats.Connect(serverURL, nats.ReconnectWait(10*time.Second))
	if err != nil {
		err = errors.New("connect nats failed, " + err.Error())
		return
	}
	//启动NC
	openNC = true
	//反馈
	return
}

// SetSubPrefix 设置全局订阅前缀
func SetSubPrefix(prefix string) {
	subPrefix = prefix
}

// 是否启动了NC
func checkNC() bool {
	return openNC && nc.IsConnected()
}

// Sub 订阅消息
func Sub(serviceCode string, topic string, cb func(msg *nats.Msg)) (err error) {
	if !checkNC() {
		return
	}
	//var sub nats.Subscription
	if subPrefix != "" {
		topic = fmt.Sprint(subPrefix, topic)
	}
	_, err = nc.Subscribe(topic, cb)
	if err != nil {
		return
	}
	CoreLog.Info("nats sub, ", topic)
	//推送统计
	pushRequest(serviceCode, "sub")
	return
}

// SubData 订阅数据包
func SubData(serviceCode string, topic string, cb func(msg *nats.Msg, action string, id int64, mark string, data interface{})) (err error) {
	if !checkNC() {
		return
	}
	if subPrefix != "" {
		topic = fmt.Sprint(subPrefix, topic)
	}
	if err = Sub(serviceCode, topic, func(msg *nats.Msg) {
		var newData dataType
		err = json.Unmarshal(msg.Data, &newData)
		if err != nil {
			//抛弃异常数据
			CoreLog.Error("nats sub load error json data, ", topic, ", err: ", err)
			return
		}
		//CoreLog.Info("nats load sub topic: ", topic, ", action: ", newData.Action, ", id: ", newData.ID, ", mark: ", newData.Mark)
		cb(msg, newData.Action, newData.ID, newData.Mark, newData.Data)
	}); err != nil {
		return
	}
	//CoreLog.Info("nats sub data, ", topic)
	return
}

// SubDataByte 订阅数据包
func SubDataByte(serviceCode string, topic string, cb func(msg *nats.Msg, action string, id int64, mark string, data []byte)) (err error) {
	if !checkNC() {
		return
	}
	if subPrefix != "" {
		topic = fmt.Sprint(subPrefix, topic)
	}
	if err = Sub(serviceCode, topic, func(msg *nats.Msg) {
		var newData dataType
		err = json.Unmarshal(msg.Data, &newData)
		if err != nil {
			//抛弃异常数据
			CoreLog.Error("nats sub load error json data, ", topic, ", err: ", err)
			return
		}
		var jsonData []byte
		jsonData, err = json.Marshal(newData.Data)
		if err != nil {
			err = errors.New(fmt.Sprint("nats sub load error json data, err data json reflect: ", err))
			return
		}
		//CoreLog.Info("nats load sub topic: ", topic, ", action: ", newData.Action, ", id: ", newData.ID, ", mark: ", newData.Mark)
		cb(msg, newData.Action, newData.ID, newData.Mark, jsonData)
	}); err != nil {
		return
	}
	//CoreLog.Info("nats sub byte, ", topic)
	return
}

// SubDataByteNoErr 订阅数据包
func SubDataByteNoErr(serviceCode string, topic string, cb func(msg *nats.Msg, action string, id int64, mark string, data []byte)) {
	if err := SubDataByte(serviceCode, topic, cb); err != nil {
		CoreLog.Error("nats sub failed, topic: ", topic)
	}
}

// ReflectDataByte 反射结构体
func ReflectDataByte(mapData []byte, rawData interface{}) (err error) {
	err = json.Unmarshal(mapData, rawData)
	if err != nil {
		err = errors.New(fmt.Sprint("err: ", err, ", json data: ", string(mapData)))
		return
	}
	return
}

// ReflectData 反射结构体
func ReflectData(mapData interface{}, rawData interface{}) (err error) {
	var jsonData []byte
	jsonData, err = json.Marshal(mapData)
	if err != nil {
		err = errors.New(fmt.Sprint("nats sub load error json data, err data json reflect: ", err))
		return
	}
	err = json.Unmarshal(jsonData, rawData)
	if err != nil {
		return
	}
	return
}

// Push 发布消息
func Push(serviceCode string, topic string, data []byte) (err error) {
	if !checkNC() {
		return
	}
	if subPrefix != "" {
		topic = fmt.Sprint(subPrefix, topic)
	}
	err = nc.Publish(topic, data)
	if err != nil {
		return
	}
	//推送统计
	pushRequest(serviceCode, "push")
	//反馈
	return
}

// PushJson 发送json数据包
func PushJson(serviceCode string, topic string, data interface{}) (err error) {
	if !checkNC() {
		return
	}
	var dataByte []byte
	dataByte, err = json.Marshal(data)
	if err != nil {
		err = errors.New("json error, " + err.Error())
		return
	}
	return Push(serviceCode, topic, dataByte)
}

func PushJsonNoErr(serviceCode string, topic string, data interface{}) {
	if !checkNC() {
		return
	}
	if err := PushJson(serviceCode, topic, data); err != nil {
		CoreLog.Error("nats push json topic: ", topic, ", data: ", data)
	}
}

// PushData 发送数据包
func PushData(serviceCode string, topic string, action string, id int64, mark string, data interface{}) (err error) {
	if !checkNC() {
		return
	}
	var dataByte []byte
	dataByte, err = json.Marshal(dataType{
		Action: action,
		ID:     id,
		Mark:   mark,
		Data:   data,
	})
	if err != nil {
		err = errors.New("json error, " + err.Error())
		return
	}
	return Push(serviceCode, topic, dataByte)
}

// PushDataNoErr 发送数据包无错误
func PushDataNoErr(serviceCode string, topic string, action string, id int64, mark string, data interface{}) {
	if !checkNC() {
		return
	}
	if err := PushData(serviceCode, topic, action, id, mark, data); err != nil {
		CoreLog.Error("nats push data topic: ", topic, ", action: ", action, ", id: ", id, ", mark: ", mark, ", data: ", data)
	}
}
