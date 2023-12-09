package FinanceTakeCut

import (
	ClassConfig "gitee.com/weeekj/weeekj_core/v5/class/config"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	FinanceDeposit "gitee.com/weeekj/weeekj_core/v5/finance/deposit"
	OrgCoreCore "gitee.com/weeekj/weeekj_core/v5/org/core"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	TestOrg "gitee.com/weeekj/weeekj_core/v5/tools/test_org"
	"testing"
)

var (
	newOrderID int64
)

func TestInitLog(t *testing.T) {
	TestInitConfig(t)
	TestSetConfig(t)
	newOrderID = int64(CoreFilter.GetRandNumber(1, 10000))
}

func TestAddLog(t *testing.T) {
	takeChannelMark, err := OrgCoreCore.Config.GetConfigVal(&ClassConfig.ArgsGetConfig{
		BindID:    TestOrg.OrgData.ID,
		Mark:      "FinanceDepositDefaultMark",
		VisitType: "admin",
	})
	if err != nil {
		t.Error("get org deposit mark config, " + err.Error())
		return
	}
	_, _, err = FinanceDeposit.SetByFrom(&FinanceDeposit.ArgsSetByFrom{
		UpdateHash: "",
		CreateInfo: CoreSQLFrom.FieldsFrom{
			System: "org",
			ID:     TestOrg.OrgData.ID,
			Mark:   "",
			Name:   "",
		},
		FromInfo: CoreSQLFrom.FieldsFrom{
			System: "org",
			ID:     TestOrg.OrgData.ID,
			Mark:   "",
			Name:   "",
		},
		ConfigMark:      takeChannelMark,
		AppendSavePrice: 100,
	})
	if err != nil {
		t.Error("set org deposit to 100, " + err.Error())
		return
	}
	cutPrice, err := AddLog(&ArgsAddLog{
		OrgID:       TestOrg.OrgData.ID,
		OrderSystem: "mall",
		OrderPrice:  100,
		OrderID:     newOrderID,
	})
	ToolsTest.ReportData(t, err, cutPrice)
}

func TestGetLogList(t *testing.T) {
	dataList, dataCount, err := GetLogList(&ArgsGetLogList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		OrgID: -1,
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetLogByOrderID(t *testing.T) {
	data, err := GetLogByOrderID(&ArgsGetLogByOrderID{
		OrgID:   TestOrg.OrgData.ID,
		OrderID: newOrderID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestClearLog(t *testing.T) {
	TestDeleteConfig(t)
	TestClearConfig(t)
}
