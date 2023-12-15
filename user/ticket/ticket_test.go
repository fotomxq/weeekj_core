package UserTicket

import (
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
	"testing"
	"time"
)

func TestInitTicket(t *testing.T) {
	TestInitConfig(t)
	TestOrg.LocalCreateUser(t)
	TestCreateConfig(t)
}

func TestAddTicket(t *testing.T) {
	err := AddTicket(&ArgsAddTicket{
		ConfigID: newConfigData.ID,
		UserID:   TestOrg.UserInfo.ID,
		Count:    2,
	})
	ToolsTest.ReportError(t, err)
	time.Sleep(time.Second * 1)
	err = AddTicket(&ArgsAddTicket{
		ConfigID: newConfigData.ID,
		UserID:   TestOrg.UserInfo.ID,
		Count:    5,
	})
	ToolsTest.ReportError(t, err)
}

func TestGetTicketCount(t *testing.T) {
	data, err := GetTicketCount(&ArgsGetTicketCount{
		ConfigID: newConfigData.ID,
		UserID:   TestOrg.UserInfo.ID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetTicketList(t *testing.T) {
	dataList, dataCount, err := GetTicketList(&ArgsGetTicketList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:        -1,
		ConfigID:     -1,
		UserID:       -1,
		NeedIsExpire: false,
		IsExpire:     false,
		IsRemove:     false,
		NeedAgg:      false,
		Search:       "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
	dataList, dataCount, err = GetTicketList(&ArgsGetTicketList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  100,
			Sort: "id",
			Desc: true,
		},
		OrgID:        TestOrg.OrgData.ID,
		ConfigID:     -1,
		UserID:       TestOrg.UserInfo.ID,
		NeedIsExpire: true,
		IsExpire:     false,
		IsRemove:     false,
		NeedAgg:      true,
		Search:       "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestUseTicket(t *testing.T) {
	err := UseTicket(&ArgsUseTicket{
		ID:       0,
		ConfigID: newConfigData.ID,
		UserID:   TestOrg.UserInfo.ID,
		Count:    1,
	})
	ToolsTest.ReportError(t, err)
}

func TestClearTicket(t *testing.T) {
	err := ClearTicket(&ArgsClearTicket{
		ConfigD: newConfigData.ID,
		OrgID:   newConfigData.OrgID,
	})
	ToolsTest.ReportError(t, err)
	TestDeleteConfig(t)
}
