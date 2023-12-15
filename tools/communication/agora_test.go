package ToolsCommunication

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

func TestInitAgora(t *testing.T) {
	TestInit(t)
}

func TestMakeAgoraToken(t *testing.T) {
	data, err := MakeAgoraToken(&ArgsMakeAgoraToken{
		RoomID:     1,
		FromSystem: 1,
		FromID:     1,
		ExpireAt:   CoreFilter.GetNowTimeCarbon().AddMinutes(3).Time,
	})
	ToolsTest.ReportData(t, err, data)
}
