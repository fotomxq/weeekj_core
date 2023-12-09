package UserLogin

import (
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newData DataQrcode
)

func TestInitQrcode(t *testing.T) {
	TestInit(t)
}

func TestMakeQrcode(t *testing.T) {
	var err error
	newData, err = MakeQrcode(&ArgsMakeQrcode{
		TokenID:    1,
		SystemMark: "system_mark",
	})
	ToolsTest.ReportData(t, err, newData)
	if err == nil {
		data, err := newData.GetJSON()
		ToolsTest.ReportData(t, err, data)
		if err == nil {
			err := newData.GetData(data)
			ToolsTest.ReportError(t, err)
		}
	}
}

func TestFinishQrcode(t *testing.T) {
	err := FinishQrcode(&ArgsFinishQrcode{
		ID:     newData.ID,
		Key:    newData.Key,
		UserID: 2,
	})
	ToolsTest.ReportError(t, err)
}

func TestCheckQrcode(t *testing.T) {
	userID, err := CheckQrcode(&ArgsCheckQrcode{
		ID:  newData.ID,
		Key: newData.Key,
	})
	ToolsTest.ReportError(t, err)
	if err == nil {
		t.Log("userID: ", userID)
	}
}
