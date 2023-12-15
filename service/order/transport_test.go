package ServiceOrder

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
	"time"
)

func TestInitTransport(t *testing.T) {
	TestInit(t)
	TestCreate(t)
	TestGetList(t)
	//提交和审核订单
	TestUpdatePost(t)
	TestUpdateAudit(t)
}

func TestUpdateTransportID(t *testing.T) {
	err := UpdateTransportID(&ArgsUpdateTransportID{
		ID:          newOrderData.ID,
		OrgID:       0,
		OrgBindID:   0,
		Des:         "",
		TransportID: 123,
	})
	ToolsTest.ReportError(t, err)
}

func TestUpdateTransportAuto(t *testing.T) {
	err := UpdateTransportAuto(&ArgsUpdateTransportAuto{
		ID:                 newOrderData.ID,
		OrgID:              0,
		OrgBindID:          0,
		Des:                "",
		TransportAllowAuto: true,
		TransportTaskAt:    CoreFilter.GetNowTime().Add(time.Second * 100),
	})
	ToolsTest.ReportError(t, err)
}

func TestClearTransport(t *testing.T) {
	TestClear(t)
}
