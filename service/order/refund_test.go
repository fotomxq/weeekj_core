package ServiceOrder

import (
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

func TestInitRefund(t *testing.T) {
	TestInitPay(t)
	TestCreatePay(t)
}

func TestRefundPost(t *testing.T) {
	errCode, err := RefundPost(&ArgsRefundPost{
		ID:             newOrderData.ID,
		OrgID:          0,
		UserID:         0,
		OrgBindID:      0,
		RefundWay:      "",
		RefundDes:      "",
		RefundFileIDs:  nil,
		RefundHaveGood: false,
	})
	ToolsTest.ReportError(t, err)
	t.Log(errCode)
}

func TestRefundAudit(t *testing.T) {
	_, err := RefundAudit(&ArgsRefundAudit{
		ID:        newOrderData.ID,
		OrgID:     0,
		UserID:    0,
		OrgBindID: 0,
		Des:       "test refund audit",
	})
	ToolsTest.ReportError(t, err)
}

func TestRefundFinish(t *testing.T) {
	err := RefundFinish(&ArgsRefundFinish{
		ID:        newOrderData.ID,
		OrgID:     0,
		UserID:    0,
		OrgBindID: 0,
		Des:       "test refund finish",
	})
	ToolsTest.ReportError(t, err)
}

func TestRefundPay(t *testing.T) {
	errCode, err := RefundPay(&ArgsRefundPay{
		ID:          newOrderData.ID,
		OrgID:       0,
		UserID:      0,
		OrgBindID:   0,
		RefundPrice: newOrderData.Price,
		Des:         "测试退款请求",
	})
	if err != nil {
		t.Error(errCode, ", ", err)
	}
}

func TestRefundFailed(t *testing.T) {
	TestInitPay(t)
	TestCreatePay(t)
	TestPayFinish(t)
	TestRefundPost(t)
	TestRefundAudit(t)
	err := RefundFailed(&ArgsRefundFailed{
		ID:        newOrderData.ID,
		OrgID:     0,
		UserID:    0,
		OrgBindID: 0,
		Des:       "test refund failed",
	})
	ToolsTest.ReportError(t, err)
}

func TestRefundCancel(t *testing.T) {
	TestInitPay(t)
	TestCreatePay(t)
	TestPayFinish(t)
	TestRefundPost(t)
	TestRefundAudit(t)
	err := RefundCancel(&ArgsRefundCancel{
		ID:        newOrderData.ID,
		OrgID:     0,
		UserID:    0,
		OrgBindID: 0,
		Des:       "test refund cancel",
	})
	ToolsTest.ReportError(t, err)
}

func TestClearRefund(t *testing.T) {
	TestClearPay(t)
}
