package CoreSQLConfig

import "testing"

func TestFieldsConfigSortsType_Scan(t *testing.T) {
	//构建虚拟化数据
	data1 := FieldsConfigsType{
		{
			Mark: "a1",
			Val:  "v1",
		},
		{
			Mark: "c2",
			Val:  "v2",
		},
		{
			Mark: "b3",
			Val:  "v3",
		},
	}
	data1json, err := data1.Value()
	if err != nil{
		t.Error(err)
	} else {
		t.Log("data1json: ", data1json)
	}
	data2 := FieldsConfigSortsType{
		{
			Mark: "a1",
			Val:  "v1",
		},
		{
			Mark: "c2",
			Val:  "v2",
		},
		{
			Mark: "b3",
			Val:  "v3",
		},
	}
	data2json, err := data2.Value()
	if err != nil{
		t.Error(err)
	} else {
		t.Log("data2json: ", data2json)
	}
}