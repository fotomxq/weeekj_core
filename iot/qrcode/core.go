package IOTQrcode

import (
	"encoding/json"
	"errors"
	"fmt"

	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	BaseQRCode "github.com/fotomxq/weeekj_core/v5/base/qrcode"
	IOTDevice "github.com/fotomxq/weeekj_core/v5/iot/device"
)

//二维码生成和存储模块

// ArgsMakeQrcode 生成指定设备的二维码参数
type ArgsMakeQrcode struct {
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id"`
	//二维码类型
	// normal 普通设备group+code二维码; normal_id 普通设备ID二维码; normal_key 带key的二维码; normal_connect 用于快速连接的二维码（带所有必要信息）
	QrcodeType string `db:"qrcode_type" json:"qrcodeType" check:"mark"`
	//尺寸
	// 10-5000px之间
	Size int `db:"size" json:"size" check:"intThan0"`
}

// MakeQrcode 生成指定设备的二维码参数
func MakeQrcode(args *ArgsMakeQrcode) (data string, err error) {
	//尺寸不能太小
	if args.Size < 10 {
		err = errors.New("size too small")
		return
	}
	if args.Size > 5000 {
		err = errors.New("size too big")
		return
	}
	//获取设备信息
	var deviceData IOTDevice.FieldsDevice
	deviceData, err = IOTDevice.GetDeviceByID(&IOTDevice.ArgsGetDeviceByID{
		ID:    args.DeviceID,
		OrgID: -1,
	})
	if err != nil {
		err = errors.New("device not exist, " + err.Error())
		return
	}
	//设备组数据
	if deviceData.GroupID < 1 {
		err = errors.New("group not have group")
		return
	}
	var groupData IOTDevice.FieldsGroup
	groupData, err = IOTDevice.GetGroupByID(&IOTDevice.ArgsGetGroupByID{
		ID: deviceData.GroupID,
	})
	if err != nil {
		err = errors.New("group not exist, " + err.Error())
		return
	}
	//根据类型输出数据
	switch args.QrcodeType {
	case "normal":
		type dataType struct {
			GroupMark string `json:"groupMark"`
			Code      string `json:"code"`
		}
		d := dataType{
			GroupMark: groupData.Mark,
			Code:      deviceData.Code,
		}
		var dJson []byte
		dJson, err = json.Marshal(d)
		if err != nil {
			return
		}
		return BaseQRCode.GetQRCode(&BaseQRCode.ArgsGetQRCode{
			Param: string(dJson),
			Size:  args.Size,
		})
	case "normal_id":
		type dataType struct {
			ID int64 `json:"id"`
		}
		d := dataType{
			ID: deviceData.ID,
		}
		var dJson []byte
		dJson, err = json.Marshal(d)
		if err != nil {
			return
		}
		return BaseQRCode.GetQRCode(&BaseQRCode.ArgsGetQRCode{
			Param: string(dJson),
			Size:  args.Size,
		})
	case "normal_key":
		type dataType struct {
			GroupMark string `json:"groupMark"`
			Code      string `json:"code"`
			Key       string `json:"key"`
		}
		var key string
		key, err = IOTDevice.GetDeviceKey(&IOTDevice.ArgsGetDeviceKey{
			ID: deviceData.ID,
		})
		if err != nil {
			return
		}
		d := dataType{
			GroupMark: groupData.Mark,
			Code:      deviceData.Code,
			Key:       key,
		}
		var dJson []byte
		dJson, err = json.Marshal(d)
		if err != nil {
			return
		}
		return BaseQRCode.GetQRCode(&BaseQRCode.ArgsGetQRCode{
			Param: string(dJson),
			Size:  args.Size,
		})
	case "normal_connect":
		type dataType struct {
			GroupMark    string `json:"m"`
			Code         string `json:"c"`
			Key          string `json:"k"`
			MQTTURL      string `json:"mu"`
			MQTTUsername string `json:"mus"`
			MQTTPassword string `json:"mpa"`
		}
		var key string
		key, err = IOTDevice.GetDeviceKey(&IOTDevice.ArgsGetDeviceKey{
			ID: deviceData.ID,
		})
		if err != nil {
			return
		}
		var mqttServerURL, mqttServerUsername, mqttServerPassword string
		mqttServerURL, err = BaseConfig.GetDataString("MQTTServerURLMQTT")
		//绕过url为空的时候，如测试或特殊业务场景，不需要mqtt服务的
		if err != nil {
			err = errors.New(fmt.Sprint("iot mqtt, get mqtt config server url, ", err))
			return
		}
		if mqttServerURL == "" {
			err = errors.New(fmt.Sprint("iot mqtt, get mqtt config server url, ", err))
			return
		}
		mqttServerUsername, err = BaseConfig.GetDataString("MQTTServerUsername")
		if err != nil {
			err = errors.New(fmt.Sprint("iot mqtt, get mqtt config server username, ", err))
			return
		}
		mqttServerPassword, err = BaseConfig.GetDataString("MQTTServerPassword")
		if err != nil {
			err = errors.New(fmt.Sprint("iot mqtt, get mqtt config server password, ", err))
			return
		}
		d := dataType{
			GroupMark:    groupData.Mark,
			Code:         deviceData.Code,
			Key:          key,
			MQTTURL:      mqttServerURL,
			MQTTUsername: mqttServerUsername,
			MQTTPassword: mqttServerPassword,
		}
		var dJson []byte
		dJson, err = json.Marshal(d)
		if err != nil {
			return
		}
		return BaseQRCode.GetQRCode(&BaseQRCode.ArgsGetQRCode{
			Param: string(dJson),
			Size:  args.Size,
		})
	default:
		err = errors.New("qrcode type error")
		return
	}
}
