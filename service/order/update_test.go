package ServiceOrder

import (
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	CoreSQLTime2 "gitee.com/weeekj/weeekj_core/v5/core/sql/time2"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
)

func TestInitUpdate(t *testing.T) {
	TestInitPay(t)
	TestCreate(t)
}

func TestGetList(t *testing.T) {
	//获取最新的订单
	dataList, dataCount, err := GetList(&ArgsGetList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		SystemMark:             "",
		OrgID:                  -1,
		UserID:                 -1,
		CompanyID:              -1,
		CreateFrom:             -1,
		Status:                 []int{},
		RefundStatus:           []int{},
		TransportID:            -1,
		NeedTransportAllowAuto: false,
		TransportAllowAuto:     false,
		PayStatus:              []int{},
		PayID:                  -1,
		PayFrom:                "",
		GoodFrom:               CoreSQLFrom.FieldsFrom{},
		TimeBetween:            CoreSQLTime2.DataCoreTime{},
		IsRemove:               false,
		IsHistory:              false,
		Search:                 "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
	if err == nil && len(dataList) > 0 {
		newOrderData = dataList[0]
		t.Log("newOrderData: ", newOrderData)
	} else {
		t.Error(dataList)
	}
}

func TestUpdatePost(t *testing.T) {
	err := UpdatePost(&ArgsUpdatePost{
		ID:        newOrderData.ID,
		OrgID:     0,
		UserID:    0,
		OrgBindID: 0,
		Des:       "测试post",
	})
	ToolsTest.ReportError(t, err)
	if err == nil {
		TestGetList(t)
	}
}

func TestUpdateAudit(t *testing.T) {
	err := UpdateAudit(&ArgsUpdateAudit{
		ID:        newOrderData.ID,
		OrgID:     0,
		UserID:    0,
		OrgBindID: 0,
		Des:       "测试audit",
	})
	ToolsTest.ReportError(t, err)
	if err == nil {
		TestGetList(t)
	}
}

func TestUpdateFinish(t *testing.T) {
	TestCreatePay(t)
	err := UpdateFinish(&ArgsUpdateFinish{
		ID:        newOrderData.ID,
		OrgID:     0,
		UserID:    0,
		OrgBindID: 0,
		Des:       "测试finish",
	})
	ToolsTest.ReportError(t, err)
	if err == nil {
		TestGetList(t)
	}
}

func TestUpdateFailed(t *testing.T) {
	TestCreate(t)
	err := UpdateFailed(&ArgsUpdateFailed{
		ID:        newOrderData.ID,
		OrgID:     0,
		UserID:    0,
		OrgBindID: 0,
		Des:       "测试failed",
	})
	ToolsTest.ReportError(t, err)
}

func TestUpdateCancel(t *testing.T) {
	TestCreate(t)
	err := UpdateCancel(&ArgsUpdateCancel{
		ID:        newOrderData.ID,
		OrgID:     0,
		UserID:    0,
		OrgBindID: 0,
		Des:       "测试failed",
	})
	ToolsTest.ReportError(t, err)
}

func TestClearUpdate(t *testing.T) {
	TestClear(t)
}
