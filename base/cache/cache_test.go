package BaseCache

import (
	"testing"

	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
)

func TestSetData(t *testing.T) {
	SetData(&DataCache{
		CreateTime:     CoreFilter.GetNowTime().Unix(),
		ExpireTime:     CoreFilter.GetNowTime().Unix() + 15,
		Mark:           "test_mark",
		Value:          "",
		ValueInt64:     0,
		ValueFloat64:   0,
		ValueBool:      false,
		ValueByte:      nil,
		ValueInterface: nil,
	})
}

func TestGetByMark(t *testing.T) {
	data, b := GetByMark(&ArgsGetByMark{
		Mark: "test_mark",
	})
	t.Log(b, data)
}

func TestGetByInterface(t *testing.T) {
	type newDataType struct {
		Data string `json:"data"`
	}
	newData := newDataType{
		Data: "abc",
	}
	SetData(&DataCache{
		CreateTime:     CoreFilter.GetNowTime().Unix(),
		ExpireTime:     CoreFilter.GetNowTime().Unix() + 15,
		Mark:           "test_mark",
		Value:          "",
		ValueInt64:     0,
		ValueFloat64:   0,
		ValueBool:      false,
		ValueByte:      nil,
		ValueInterface: newData,
	})
	newData2 := newDataType{}
	b := GetByMarkInterface("test_mark", newData2)
	t.Log(b, newData2)
	data, b := GetByMark(&ArgsGetByMark{
		Mark: "test_mark",
	})
	newData2 = data.ValueInterface.(newDataType)
	t.Log("step2: ", newData2)
	newData3, b := GetByMarkInterfaceReturn("test_mark")
	t.Log("step3: ", newData3.(newDataType))
}
