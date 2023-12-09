package IOTMQTT

import (
	"encoding/json"
	IOTDevice "gitee.com/weeekj/weeekj_core/v5/iot/device"
	"reflect"
	"testing"
)

func TestCheckDeviceAndOrg(t *testing.T) {
	type dataType struct {
		//配对密钥
		Keys IOTDevice.ArgsCheckDeviceKey `json:"keys"`
		//组织ID
		OrgID int64 `db:"org_id" json:"orgID"`
	}
	resultData := dataType{
		Keys: IOTDevice.ArgsCheckDeviceKey{
			GroupMark: "gMark",
			Code:      "dCode",
			NowTime:   0,
			Rand:      "",
			Key:       "",
		},
		OrgID: 3,
	}
	resultDataJSON, err := json.Marshal(resultData)
	if err != nil {
		t.Error(err)
		return
	}
	var resultData2 dataType
	err = json.Unmarshal(resultDataJSON, &resultData2)
	if err != nil {
		t.Error(err)
		return
	}
	var deviceKeys IOTDevice.ArgsCheckDeviceKey
	var orgID int64
	for key := 0; key < reflect.TypeOf(resultData2).Elem().NumField(); key += 1 {
		field := reflect.TypeOf(resultData2).Elem().Field(key)
		val := reflect.ValueOf(resultData2).Elem().Field(key)
		fieldMark := field.Tag.Get("json")
		switch fieldMark {
		case "keys":
			deviceKeys = val.Interface().(IOTDevice.ArgsCheckDeviceKey)
			break
		case "orgID":
			orgID = val.Interface().(int64)
			break
		}
	}
	t.Log("keys: ", deviceKeys)
	t.Log("orgID: ", orgID)
}
