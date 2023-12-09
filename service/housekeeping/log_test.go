package ServiceHousekeeping

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQLAddress "gitee.com/weeekj/weeekj_core/v5/core/sql/address"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	CoreSQLTime "gitee.com/weeekj/weeekj_core/v5/core/sql/time"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	TestOrg "gitee.com/weeekj/weeekj_core/v5/tools/test_org"
	"testing"
	"time"
)

var (
	newLogData FieldsLog
)

func TestInitLog(t *testing.T) {
	TestInit(t)
	TestSetBind(t)
}

func TestCreateLog(t *testing.T) {
	var err error
	newLogData, _, err = CreateLog(&ArgsCreateLog{
		UserID:        TestOrg.UserInfo.ID,
		NeedAt:        CoreFilter.GetNowTime(),
		OrgID:         TestOrg.OrgData.ID,
		BindID:        TestOrg.BindData.ID,
		OtherBinds:    []int64{},
		MallProductID: 1,
		OrderID:       0,
		Currency:      86,
		Price:         100,
		PayAt:         time.Time{},
		Des:           "测试服务",
		Address:       CoreSQLAddress.FieldsAddress{},
		Params:        nil,
	})
	ToolsTest.ReportData(t, err, newLogData)
}

func TestGetLogList(t *testing.T) {
	dataList, dataCount, err := GetLogList(&ArgsGetLogList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		UserID:       -1,
		OrgID:        -1,
		BindID:       -1,
		OrderID:      -1,
		NeedIsNeed:   false,
		IsNeed:       false,
		NeedIsFinish: false,
		IsFinish:     false,
		NeedIsPay:    false,
		IsPay:        false,
		TimeBetween:  CoreSQLTime.DataCoreTime{},
		IsRemove:     false,
		Search:       "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetLogID(t *testing.T) {
	data, err := GetLogID(&ArgsGetLogID{
		ID:     newLogData.ID,
		OrgID:  TestOrg.OrgData.ID,
		BindID: TestOrg.BindData.ID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestUpdateLogBind(t *testing.T) {
	err := UpdateLogBind(&ArgsUpdateLogBind{
		ID:         newLogData.ID,
		OrgID:      TestOrg.OrgData.ID,
		BindID:     TestOrg.BindData.ID,
		NewBindID:  TestOrg.BindData.ID,
		OtherBinds: []int64{},
	})
	ToolsTest.ReportError(t, err)
}

func TestUpdateLogNeedAt(t *testing.T) {
	err := UpdateLogNeedAt(&ArgsUpdateLogNeedAt{
		ID:     newLogData.ID,
		OrgID:  TestOrg.OrgData.ID,
		BindID: TestOrg.BindData.ID,
		NeedAt: CoreFilter.GetISOByTime(CoreFilter.GetNowTime()),
	})
	ToolsTest.ReportError(t, err)
}

func TestUpdateLogPay(t *testing.T) {
	err := UpdateLogPay(&ArgsUpdateLogPay{
		ID:     newLogData.ID,
		OrgID:  TestOrg.OrgData.ID,
		BindID: TestOrg.BindData.ID,
	})
	ToolsTest.ReportError(t, err)
}

func TestUpdateLogPrice(t *testing.T) {
	err := UpdateLogPrice(&ArgsUpdateLogPrice{
		ID:     newLogData.ID,
		OrgID:  TestOrg.OrgData.ID,
		BindID: TestOrg.BindData.ID,
		Price:  6000,
	})
	ToolsTest.ReportError(t, err)
}

func TestPayLog(t *testing.T) {
	payData, errCode, err := PayLog(&ArgsPayLog{
		ID:     newLogData.ID,
		OrgID:  TestOrg.OrgData.ID,
		BindID: TestOrg.BindData.ID,
		PaymentChannel: CoreSQLFrom.FieldsFrom{
			System: "cash",
			ID:     0,
			Mark:   "",
			Name:   "",
		},
	})
	ToolsTest.ReportError(t, err)
	if err != nil {
		t.Error(errCode)
	} else {
		t.Log(payData)
	}
}

func TestUpdateLogFinish(t *testing.T) {
	err := UpdateLogFinish(&ArgsUpdateLogFinish{
		ID:     newLogData.ID,
		OrgID:  TestOrg.OrgData.ID,
		BindID: TestOrg.BindData.ID,
	})
	ToolsTest.ReportError(t, err)
}

func TestCloseLog(t *testing.T) {
	err := CloseLog(&ArgsCloseLog{
		ID:    newLogData.ID,
		OrgID: TestOrg.OrgData.ID,
	})
	ToolsTest.ReportError(t, err)
}

func TestClearLog(t *testing.T) {
	TestDeleteBind(t)
	TestClear(t)
}
