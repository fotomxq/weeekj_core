package PedometerCore

import (
	CoreRPCX "github.com/fotomxq/weeekj_core/v5/core/rpcx"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

var (
	isRun = false
)

func TestInit(t *testing.T) {
	if isRun {
		return
	}
	isRun = true
	ToolsTest.Init(t)
}

func TestSetConfig(t *testing.T) {
	if err := SetConfig(&ArgsSetConfig{
		Mark: "user", Count: 10, ExpireAdd: "30h", MaxCount: 100,
	}); err != nil {
		t.Error(err)
	}
	if err := SetConfig(&ArgsSetConfig{
		Mark: "token", Count: 1, ExpireAdd: "10h", MinCount: 1, MaxCount: 5, IsAdd: true,
	}); err != nil {
		t.Error(err)
	}
	if err := SetConfig(&ArgsSetConfig{
		Mark: "pg", Count: 10, ExpireAdd: "11h", MinCount: 1, MaxCount: 30, IsAdd: true,
	}); err != nil {
		t.Error(err)
	}
	if err := SetConfig(&ArgsSetConfig{
		Mark: "prev", Count: 3, ExpireAdd: "12h", MinCount: 1, MaxCount: 10, IsAdd: true,
	}); err != nil {
		t.Error(err)
	}
}

func TestGetConfig(t *testing.T) {
	data, err := getConfig("user")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data)
	}
}

func TestSetData(t *testing.T) {
	newData, err := SetData(&ArgsSetData{
		CreateInfo: CoreSQLFrom.FieldsFrom{System: "user", ID: 1, Mark: "mark"}, Count: 1,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("change new data by 1: ", newData)
	}
}

func TestGetData(t *testing.T) {
	newData1, err := GetData(&CoreRPCX.ArgsFrom{
		From: CoreSQLFrom.FieldsFrom{System: "user", ID: 1, Mark: "mark"},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("get data: ", newData1)
	}
}

func TestClearData(t *testing.T) {
	if err := ClearData(&ArgsClearData{
		System: "user",
	}); err != nil {
		t.Error(err)
	}
}

func TestNextData(t *testing.T) {
	if newData, err := NextData(&CoreRPCX.ArgsFrom{
		From: CoreSQLFrom.FieldsFrom{System: "prev", ID: 1, Mark: "mark"},
	}); err != nil {
		t.Error(err)
	} else {
		t.Log("change next new data: ", newData)
	}
}

func TestPrevData(t *testing.T) {
	if newData, err := PrevData(&CoreRPCX.ArgsFrom{
		From: CoreSQLFrom.FieldsFrom{System: "prev", ID: 1, Mark: "mark"},
	}); err != nil {
		t.Error(err)
	} else {
		t.Log("change prev new data: ", newData)
	}
}

func TestGetCount(t *testing.T) {
	newData := GetCount(&CoreRPCX.ArgsFrom{
		From: CoreSQLFrom.FieldsFrom{System: "prev", ID: 1, Mark: "mark"},
	})
	t.Log("get data by mark: ", newData)
}

func TestCheckData(t *testing.T) {
	newData := CheckData(&CoreRPCX.ArgsFrom{
		From: CoreSQLFrom.FieldsFrom{System: "prev", ID: 1, Mark: "mark"},
	})
	t.Log("get data by mark: ", newData)
}

// 检查陌生IP数据
func TestCheckDataByIPNoHave(t *testing.T) {
	newData := CheckData(&CoreRPCX.ArgsFrom{
		From: CoreSQLFrom.FieldsFrom{System: "x1", ID: 1, Mark: "213.123.123.213"},
	})
	if newData {
		t.Log("get data by mark: ", newData)
	}
}

func TestReturnData(t *testing.T) {
	newData, err := ReturnData(&CoreRPCX.ArgsFrom{
		From: CoreSQLFrom.FieldsFrom{System: "prev", ID: 1, Mark: "mark"},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("get data : ", newData)
	}
}

func TestDeleteAllData(t *testing.T) {
	if _, err := CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "core_pedometer", "true", nil); err != nil {
		t.Error(err)
	}
	if _, err := CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "core_pedometer_config", "true", nil); err != nil {
		t.Error(err)
	}
}
