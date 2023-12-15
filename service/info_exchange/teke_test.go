package ServiceInfoExchange

import (
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

func TestInitTake(t *testing.T) {
	TestInitInfo(t)
	TestCreateInfo(t)
}

func TestTakeInfo(t *testing.T) {
	_, err := TakeInfo(&ArgsTakeInfo{
		UserID: newInfoData.UserID,
		InfoID: newInfoData.ID,
		Des:    "测试",
		Params: nil,
	})
	ToolsTest.ReportError(t, err)
}

func TestGetAnalysisTake(t *testing.T) {
	data, err := GetAnalysisTake(&ArgsGetAnalysisTake{
		UserID: newInfoData.UserID,
		InfoID: 0,
	})
	ToolsTest.ReportData(t, err, data)
	data2, err := GetAnalysisTake(&ArgsGetAnalysisTake{
		UserID: newInfoData.UserID,
		InfoID: newInfoData.ID,
	})
	ToolsTest.ReportData(t, err, data2)
}

func TestGetTakeList(t *testing.T) {
	dataList, dataCount, err := GetTakeList(&ArgsGetTakeList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		UserID: -1,
		InfoID: -1,
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestDeleteTake(t *testing.T) {
	err := DeleteTake(&ArgsDeleteTake{
		ID:         newInfoData.ID,
		UserID:     newInfoData.UserID,
		TakeUserID: newInfoData.UserID,
	})
	ToolsTest.ReportError(t, err)
}

func TestClearTake(t *testing.T) {
	TestDeleteInfo(t)
	TestClearInfo(t)
}
