package CoreFilter

import "testing"

//本模块主要测试string结构对任意类型转化的问题
func TestGetBoolByInterface(t *testing.T) {
	data, err := GetBoolByInterface("true")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data)
	}
	data2, err := GetFloat64ByInterface("0.012")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data2)
	}
	data3, err := GetIntByInterface("1")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data3)
	}
	data4, err := GetInt64ByInterface("12213213")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data4)
	}
}
