package UserIntegral

import (
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

var (
	orgID  int64 = 123
	userID int64 = 234
)

func TestInit(t *testing.T) {
	ToolsTest.Init(t)
}

func TestAddCount(t *testing.T) {
	err := AddCount(&ArgsAddCount{
		OrgID:    orgID,
		UserID:   userID,
		AddCount: 10,
		Des:      "测试调整",
	})
	ToolsTest.ReportError(t, err)
}

func TestGetList(t *testing.T) {
	dataList, dataCount, err := GetList(&ArgsGetList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:  0,
		UserID: 0,
		Min:    0,
		Max:    0,
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetLogList(t *testing.T) {
	dataList, dataCount, err := GetLogList(&ArgsGetLogList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:  0,
		UserID: 0,
		Min:    0,
		Max:    0,
		Search: "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetUser(t *testing.T) {
	data, err := GetUser(&ArgsGetUser{
		OrgID:  orgID,
		UserID: userID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetUserCount(t *testing.T) {
	data := GetUserCount(orgID, userID)
	ToolsTest.ReportData(t, nil, data)
	if data != 10 {
		t.Error("count not 10")
	}
}

// 积分增减测试
func TestAddCount2(t *testing.T) {
	TestAddCount(t)
	data := GetUserCount(orgID, userID)
	ToolsTest.ReportData(t, nil, data)
	if data != 20 {
		t.Error("count not 20")
	}
}

func TestClearUser(t *testing.T) {
	err := ClearUser(&ArgsClearUser{
		UserID: userID,
	})
	ToolsTest.ReportError(t, err)
}

func TestClearOrg(t *testing.T) {
	TestAddCount(t)
	err := ClearOrg(&ArgsClearOrg{
		OrgID: orgID,
	})
	ToolsTest.ReportError(t, err)
}
