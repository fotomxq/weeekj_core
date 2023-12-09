package ServiceOrder

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	FinanceDeposit "gitee.com/weeekj/weeekj_core/v5/finance/deposit"
	FinancePay "gitee.com/weeekj/weeekj_core/v5/finance/pay"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
	"time"
)

func TestInitPay(t *testing.T) {
	TestInit(t)
	TestCreate(t)
	TestGetList(t)
	//提交和审核订单
	TestUpdatePost(t)
	TestUpdateAudit(t)
	//创建储蓄账户
	_, err := FinanceDeposit.SetConfig(&FinanceDeposit.ArgsSetConfig{
		Mark:             "savings",
		Name:             "储蓄",
		Des:              "储蓄描述",
		Currency:         86,
		TakeOut:          true,
		TakeLimit:        0,
		OnceSaveMinLimit: 0,
		OnceSaveMaxLimit: 0,
		OnceTakeMinLimit: 0,
		OnceTakeMaxLimit: 0,
		Configs:          CoreSQLConfig.FieldsConfigsType{},
	})
	if err != nil {
		t.Error(err)
	}
}

func TestCreatePay(t *testing.T) {
	payData, errCode, err := CreatePay(&ArgsCreatePay{
		IDs: []int64{newOrderData.ID},
		PaymentChannel: CoreSQLFrom.FieldsFrom{
			System: "cash",
			ID:     0,
			Mark:   "",
			Name:   "",
		},
		Des: "支付订单测试",
	})
	if err != nil {
		t.Error(errCode, err)
		return
	} else {
		t.Log(payData)
	}
	//强制完成支付处理
	_, _, _, _, err = FinancePay.UpdateStatusClient(&FinancePay.ArgsUpdateStatusClient{
		CreateInfo: payData.PaymentCreate,
		ID:         payData.ID,
		Key:        "",
		Params:     []CoreSQLConfig.FieldsConfigType{},
		IP:         "0.0.0.0",
	})
	if err != nil {
		t.Error(err)
		return
	}
	_, err = FinancePay.UpdateStatusFinish(&FinancePay.ArgsUpdateStatusFinish{
		CreateInfo: payData.PaymentCreate,
		ID:         payData.ID,
		Key:        "",
		Params:     []CoreSQLConfig.FieldsConfigType{},
	})
	if err != nil {
		t.Error(err)
		return
	}
	//启动自动验证服务
	//go runPay()
	time.Sleep(time.Second * 1)
}

func TestCheckPay(t *testing.T) {
	b := CheckPay(&ArgsCheckPay{
		IDs:    []int64{newOrderData.ID},
		OrgID:  0,
		UserID: 0,
	})
	if !b {
		t.Error("订单未支付完成..", newOrderData.ID)
	} else {
		t.Log("订单支付完成")
	}
}

func TestPayFailed(t *testing.T) {
	TestCreate(t)
	TestGetList(t)
	//提交和审核订单
	TestUpdatePost(t)
	TestUpdateAudit(t)
	err := PayFailed(&ArgsPayFailed{
		ID:        newOrderData.ID,
		OrgID:     0,
		UserID:    0,
		OrgBindID: 0,
		Des:       "测试pay_failed",
	})
	ToolsTest.ReportError(t, err)
}

func TestPayFinish(t *testing.T) {
	TestCreate(t)
	TestGetList(t)
	//提交和审核订单
	TestUpdatePost(t)
	TestUpdateAudit(t)
	err := PayFinish(&ArgsPayFinish{
		ID:        newOrderData.ID,
		OrgID:     0,
		UserID:    0,
		OrgBindID: 0,
		Des:       "测试pay_finish",
	})
	ToolsTest.ReportError(t, err)
}

func TestClearPay(t *testing.T) {
	TestClear(t)
}
