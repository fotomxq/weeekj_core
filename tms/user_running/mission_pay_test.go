package TMSUserRunning

import (
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
)

func TestMissionPayInit(t *testing.T) {
	TestMissionInit(t)
	TestCreateMission(t)
}

func TestUpdateRunPayPrice(t *testing.T) {
	err := UpdateRunPayPrice(&ArgsUpdateRunPayPrice{
		ID:     newMissionData.ID,
		UserID: newMissionData.UserID,
		Price:  100,
	})
	ToolsTest.ReportError(t, err)
}

func TestPayRunPay(t *testing.T) {
	payData, errCode, err := PayRunPay(&ArgsPayRunPay{
		ID:     newMissionData.ID,
		UserID: newMissionData.UserID,
		PaymentChannel: CoreSQLFrom.FieldsFrom{
			System: "cash",
			ID:     0,
			Mark:   "",
			Name:   "",
		},
	})
	if err != nil {
		t.Error(errCode)
		t.Error(err)
	} else {
		t.Log("pay data: ", payData)
	}
}

func TestUpdateOrderPrice(t *testing.T) {
	err := UpdateOrderPrice(&ArgsUpdateOrderPrice{
		ID:            newMissionData.ID,
		RoleID:        newMissionData.RoleID,
		OrderPrice:    100,
		OrderDesFiles: newMissionData.OrderDesFiles,
		OrderDes:      newMissionData.Des,
	})
	ToolsTest.ReportError(t, err)
}

func TestPayOrder(t *testing.T) {
	payData, errCode, err := PayOrder(&ArgsPayOrder{
		ID:     newMissionData.ID,
		UserID: newMissionData.UserID,
		PaymentChannel: CoreSQLFrom.FieldsFrom{
			System: "cash",
			ID:     0,
			Mark:   "",
			Name:   "",
		},
	})
	if err != nil {
		t.Error(errCode)
		t.Error(err)
	} else {
		t.Log("pay data: ", payData)
	}
}

func TestMissionPayClear(t *testing.T) {
	TestDeleteMission(t)
	TestMissionClear(t)
}
