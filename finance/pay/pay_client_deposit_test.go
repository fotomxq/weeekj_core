package FinancePay

import (
	"testing"
	"time"

	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	TestFinanceDeposit "github.com/fotomxq/weeekj_core/v5/tools/test_finance_deposit"
)

// 本测试部分不进行跑通测试，而是进行和储蓄相关的存取测试
func TestInitDeposit(t *testing.T) {
	TestInit(t)
	TestFinanceDepositHave(t)
}

// 确保资金存在
func TestDepositSetPrice0(t *testing.T) {
	t.Log("TestFinanceDeposit.LocalConfig: ", TestFinanceDeposit.LocalConfig)
	TestFinanceDeposit.LocalSetPrice(
		t,
		CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     testDepositUserID2,
			Mark:   "",
			Name:   "测试用户2",
		},
		CoreSQLFrom.FieldsFrom{},
		TestFinanceDeposit.LocalConfig.Mark,
		50000,
	)
	TestFinanceDeposit.LocalSetPrice(
		t,
		CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     testDepositUserID,
			Mark:   "",
			Name:   "测试用户",
		},
		CoreSQLFrom.FieldsFrom{},
		TestFinanceDeposit.LocalConfig.Mark,
		0,
	)
	t.Log(TestFinanceDeposit.LocalDepositData)
}

// 获取存储资金量
func TestDepositGet(t *testing.T) {
	TestFinanceDeposit.LocalGetPriceByFrom(CoreSQLFrom.FieldsFrom{
		System: "user",
		ID:     testDepositUserID,
		Mark:   "",
		Name:   "",
	}, CoreSQLFrom.FieldsFrom{}, "test_mark")
}

// 储蓄存入资金的请求的测试
func TestDepositCreateTo(t *testing.T) {
	var errCode string
	var err error
	//发起测试请求
	// 从现金渠道，给储蓄账户付款
	payData, errCode, err = Create(&ArgsCreate{
		CreateInfo: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     testDepositUserID,
			Mark:   "",
			Name:   "",
		},
		PaymentCreate: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     testDepositUserID,
			Mark:   "",
			Name:   "",
		},
		PaymentChannel: CoreSQLFrom.FieldsFrom{
			System: "cash",
			ID:     0,
			Mark:   "",
			Name:   "",
		},
		PaymentFrom: CoreSQLFrom.FieldsFrom{},
		TakeCreate: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     testDepositUserID,
			Mark:   "",
			Name:   "",
		},
		TakeChannel: CoreSQLFrom.FieldsFrom{
			System: "deposit",
			ID:     0,
			Mark:   "test_mark",
			Name:   "",
		},
		TakeFrom: CoreSQLFrom.FieldsFrom{},
		Des:      "测试入账",
		ExpireAt: CoreFilter.GetNowTime().Add(time.Hour * 1),
		Currency: 86,
		Price:    1050,
		Params:   nil,
	})
	if err != nil {
		t.Error(errCode, err)
		return
	}
	t.Log(payData)
	//客户端通过请求
	if _, _, _, _, err := UpdateStatusClient(&ArgsUpdateStatusClient{
		CreateInfo: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     testDepositUserID,
			Mark:   "",
			Name:   "",
		},
		ID:     payData.ID,
		Key:    "",
		Params: nil,
		IP:     "0.0.0.0",
	}); err != nil {
		t.Error(err)
		return
	}
	//保存请求完成前的资金
	TestDepositGet(t)
	beforeSavePrice := TestFinanceDeposit.LocalDepositData.SavePrice
	//服务端通过请求
	if errCode, err = UpdateStatusFinish(&ArgsUpdateStatusFinish{
		CreateInfo: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     testDepositUserID,
			Mark:   "",
			Name:   "",
		},
		ID:     payData.ID,
		Key:    "",
		Params: nil,
	}); err != nil {
		t.Error(errCode, err)
		return
	} else {
		TestGetOne(t)
	}
	//保存请求完成后的资金
	TestDepositGet(t)
	afterSavePrice := TestFinanceDeposit.LocalDepositData.SavePrice
	//检查账户额度
	t.Log("beforeSavePrice: ", beforeSavePrice, ", afterSavePrice: ", afterSavePrice)
	//检查资金变动情况是否符合？
	if afterSavePrice-beforeSavePrice != 1050 {
		t.Error("price error")
	} else {
		t.Log("price ok and is 10.50")
	}
}

// 储蓄存入资金后退款测试
func TestDepositCreateToAndRefund(t *testing.T) {
	//发起退款申请
	TestUpdateStatusRefund(t)
	//审核通过
	errCode, err := UpdateStatusRefundAudit(&ArgsUpdateStatusRefundAudit{
		CreateInfo: CoreSQLFrom.FieldsFrom{
			System: "system",
			ID:     0,
			Mark:   "test",
			Name:   "",
		},
		ID:          payData.ID,
		Key:         "",
		Params:      nil,
		RefundPrice: 550,
	})
	if err != nil {
		t.Error(errCode, err)
	}
	//保存请求完成前的资金
	TestDepositGet(t)
	beforeSavePrice := TestFinanceDeposit.LocalDepositData.SavePrice
	//完成退款行为
	TestUpdateStatusRefundFinish(t)
	//保存请求完成后的资金
	TestDepositGet(t)
	afterSavePrice := TestFinanceDeposit.LocalDepositData.SavePrice
	//检查账户额度
	TestDepositGet(t)
	t.Log("beforeSavePrice: ", beforeSavePrice, ", afterSavePrice: ", afterSavePrice)
	//检查资金变动情况是否符合？
	if beforeSavePrice-afterSavePrice != 550 {
		t.Error("price error, result is 550")
	} else {
		t.Log("price ok and is 5.00")
	}
}

var (
	//第二个储蓄账户
	depositData1Price int64
	depositData2Price int64
)

// 获取两个储蓄账户信息
func TestDepositGetPrice(t *testing.T) {
	TestFinanceDeposit.LocalGetPriceByFrom(CoreSQLFrom.FieldsFrom{
		System: "user",
		ID:     testDepositUserID,
		Mark:   "",
		Name:   "",
	}, CoreSQLFrom.FieldsFrom{}, "test_mark")
	depositData1Price = TestFinanceDeposit.LocalDepositData.SavePrice
	TestFinanceDeposit.LocalGetPriceByFrom(CoreSQLFrom.FieldsFrom{
		System: "user",
		ID:     testDepositUserID2,
		Mark:   "",
		Name:   "",
	}, CoreSQLFrom.FieldsFrom{}, "test_mark")
	depositData2Price = TestFinanceDeposit.LocalDepositData.SavePrice
}

// 以储蓄来源存入资金的请求测试
func TestDepositCreateFrom(t *testing.T) {
	var errCode string
	var err error
	//给测试账户1写入3050
	TestFinanceDeposit.LocalSetPrice(t, CoreSQLFrom.FieldsFrom{
		System: "user",
		ID:     testDepositUserID,
		Mark:   "",
		Name:   "",
	}, CoreSQLFrom.FieldsFrom{}, "test_mark", 3050)
	//发起测试请求
	// 测试用户2的账户2，余额500.00元，给账户1付款30.5元
	payData, errCode, err = Create(&ArgsCreate{
		CreateInfo: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     testDepositUserID,
			Mark:   "",
			Name:   "",
		},
		PaymentCreate: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     testDepositUserID,
			Mark:   "",
			Name:   "",
		},
		PaymentChannel: CoreSQLFrom.FieldsFrom{
			System: "deposit",
			ID:     0,
			Mark:   "test_mark",
			Name:   "",
		},
		PaymentFrom: CoreSQLFrom.FieldsFrom{},
		TakeCreate: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     testDepositUserID2,
			Mark:   "",
			Name:   "",
		},
		TakeChannel: CoreSQLFrom.FieldsFrom{
			System: "deposit",
			ID:     0,
			Mark:   "test_mark",
			Name:   "",
		},
		TakeFrom: CoreSQLFrom.FieldsFrom{},
		Des:      "测试账户之间汇款",
		ExpireAt: CoreFilter.GetNowTime().Add(time.Hour * 1),
		Currency: 86,
		Price:    3050,
		Params:   nil,
	})
	if err != nil {
		t.Error(errCode, err)
		return
	}
	t.Log(payData)
	//客户端通过请求
	if _, _, _, _, err := UpdateStatusClient(&ArgsUpdateStatusClient{
		CreateInfo: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     testDepositUserID,
			Mark:   "",
			Name:   "",
		},
		ID:     payData.ID,
		Key:    "",
		Params: nil,
		IP:     "0.0.0.0",
	}); err != nil {
		t.Error(err)
		return
	}
	//保存请求完成前的资金
	TestDepositGetPrice(t)
	depositData1Before := depositData1Price
	depositData2Before := depositData2Price
	t.Log("depositData1Before: ", depositData1Before, ", depositData2Before: ", depositData2Before)
	//服务端通过请求
	if errCode, err = UpdateStatusFinish(&ArgsUpdateStatusFinish{
		CreateInfo: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     testDepositUserID,
			Mark:   "",
			Name:   "",
		},
		ID:     payData.ID,
		Key:    "",
		Params: nil,
	}); err != nil {
		t.Error(errCode, err)
		return
	}
	//保存请求完成后的资金
	TestDepositGetPrice(t)
	depositData1After := depositData1Price
	depositData2After := depositData2Price
	//检查账户额度
	t.Log("depositData1Before: ", depositData1Before, ", depositData1After: ", depositData1After)
	//检查资金变动情况是否符合？相差应该符合30.5
	if depositData1Before-depositData1After != 3050 {
		t.Error("price error")
	} else {
		t.Log("price ok and is 30.50")
	}
	//检查付款的账户额度
	t.Log("depositData2Before: ", depositData2Before, ", depositData2After: ", depositData2After)
	//检查资金变动情况是否符合？相差应该符合30.5
	if depositData2After-depositData2Before != 3050 {
		t.Error("price error")
	} else {
		t.Log("price ok and is 30.50")
	}
}

// 以储蓄来源退款的测试
func TestDepositCreateFromAndRefund(t *testing.T) {
	//对之前的请求发起退款申请
	//发起退款申请
	TestUpdateStatusRefund(t)
	//审核通过
	errCode, err := UpdateStatusRefundAudit(&ArgsUpdateStatusRefundAudit{
		CreateInfo: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     testDepositUserID,
			Mark:   "",
			Name:   "",
		},
		ID:          payData.ID,
		Key:         "",
		Params:      nil,
		RefundPrice: 750,
	})
	if err != nil {
		t.Error(errCode, err)
	}
	//保存请求完成前的资金
	TestDepositGetPrice(t)
	depositData1Before := depositData1Price
	depositData2Before := depositData2Price
	//检查账户额度
	t.Log("depositData1Before: ", depositData1Before, ", depositData2Before: ", depositData2Before)
	//完成退款行为
	TestUpdateStatusRefundFinish(t)
	//保存请求完成后的资金
	TestDepositGetPrice(t)
	depositData1After := depositData1Price
	depositData2After := depositData2Price
	//检查账户额度
	t.Log("depositData1Before: ", depositData1Before, ", depositData1After: ", depositData1After)
	//检查资金变动情况是否符合？相差应该符合7.50
	if depositData1After-depositData1Before != 750 {
		t.Error("price error")
	} else {
		t.Log("price ok and is 7.50")
	}
	//检查付款的账户额度
	t.Log("depositData2Before: ", depositData2Before, ", depositData2After: ", depositData2After)
	//检查资金变动情况是否符合？相差应该符合7.50
	if depositData2Before-depositData2After != 750 {
		t.Error("price error")
	} else {
		t.Log("price ok and is 7.50")
	}
}

// 清理工作
func TestDepositClear(t *testing.T) {
	//删除资金池
	TestFinanceDeposit.LocalDeleteConfig(t, TestFinanceDeposit.LocalConfig.Mark)
}
