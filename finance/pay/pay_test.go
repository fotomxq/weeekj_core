package FinancePay

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	FinanceDeposit "github.com/fotomxq/weeekj_core/v5/finance/deposit"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	TestFinanceDeposit "github.com/fotomxq/weeekj_core/v5/tools/test_finance_deposit"
	"testing"
	"time"
)

var (
	isInit  = false
	payData FieldsPayType
	//测试用的用户ID
	testDepositUserID  = int64(CoreFilter.GetRandNumber(1, 9999))
	testDepositUserID2 = int64(CoreFilter.GetRandNumber(1, 9999))
	//通用测试来源
	infos CoreSQLFrom.FieldsFrom
)

func TestInit(t *testing.T) {
	if isInit {
		return
	}
	isInit = true
	ToolsTest.Init(t)
}

// 确保储蓄账户存在
func TestFinanceDepositHave(t *testing.T) {
	//查询储蓄账户是否存在
	configData, err := FinanceDeposit.GetConfigByMark(&FinanceDeposit.ArgsGetConfigByMark{
		Mark: "test_mark",
	})
	if err != nil {
		//构建储蓄账户
		TestFinanceDeposit.LocalCreateConfig(t, "test_mark", true, "测试资金池", "测试资金池描述")
		infos = CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     testDepositUserID,
			Mark:   "",
			Name:   "测试用户",
		}
	} else {
		TestFinanceDeposit.LocalConfig = configData
		t.Log("deposit config data: ", configData)
	}
}

func TestRun(t *testing.T) {
	time.Sleep(time.Second * 3)
}

func TestCreateToCash(t *testing.T) {
	var errCode string
	var err error
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
		Des:      "测试创建转账",
		ExpireAt: CoreFilter.GetNowTime().Add(time.Hour * 1),
		Currency: 86,
		Price:    int64(CoreFilter.GetRandNumber(1, 3000) * 3),
		Params:   nil,
	})
	if err != nil {
		t.Error(errCode, err)
	} else {
		t.Log(payData)
	}
	//转账创建测试
	payData2, errCode, err := Create(&ArgsCreate{
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
			ID:     testDepositUserID,
			Mark:   "",
			Name:   "",
		},
		TakeChannel: CoreSQLFrom.FieldsFrom{
			System: "cash",
			ID:     0,
			Mark:   "",
			Name:   "",
		},
		TakeFrom: CoreSQLFrom.FieldsFrom{},
		Des:      "测试创建转账",
		ExpireAt: CoreFilter.GetNowTime().Add(time.Hour * 1),
		Currency: 86,
		Price:    int64(CoreFilter.GetRandNumber(1, 3000) * 3),
		Params:   nil,
	})
	if err != nil {
		t.Error(errCode, err)
	} else {
		t.Log(payData2)
	}
}

func TestCreate(t *testing.T) {
	var errCode string
	var err error
	//资金转入测试
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
		Des:      "资金转入测试",
		ExpireAt: CoreFilter.GetNowTime().Add(time.Hour * 1),
		Currency: 86,
		Price:    int64(CoreFilter.GetRandNumber(1, 3000) * 3),
		Params:   nil,
	})
	if err != nil {
		t.Error(errCode, err)
	} else {
		t.Log(payData)
	}
}

func TestGetOne(t *testing.T) {
	data, err := GetOne(&ArgsGetOne{
		ID:  payData.ID,
		Key: "",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data)
	}
}

func TestCheckPaymentFrom(t *testing.T) {
	err := CheckPaymentFrom(&ArgsCheckPaymentFrom{
		ID: payData.ID,
		TakeFrom: CoreSQLFrom.FieldsFrom{
			System: payData.PaymentCreate.System,
			ID:     payData.PaymentCreate.ID,
			Mark:   "",
			Name:   "",
		},
	})
	if err != nil {
		t.Error(err)
	}
}

func TestCheckTakeFrom(t *testing.T) {
	err := CheckTakeFrom(&ArgsCheckTakeFrom{
		ID: payData.ID,
		TakeFrom: CoreSQLFrom.FieldsFrom{
			System: payData.TakeCreate.System,
			ID:     payData.TakeCreate.ID,
			Mark:   "",
			Name:   "",
		},
	})
	if err != nil {
		t.Error(err)
	}
}

func TestGetList(t *testing.T) {
	data, count, err := GetList(&ArgsGetList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		Status:         []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		PaymentCreate:  CoreSQLFrom.FieldsFrom{},
		PaymentChannel: CoreSQLFrom.FieldsFrom{},
		PaymentFrom:    CoreSQLFrom.FieldsFrom{},
		TakeCreate:     CoreSQLFrom.FieldsFrom{},
		TakeChannel:    CoreSQLFrom.FieldsFrom{},
		TakeFrom:       CoreSQLFrom.FieldsFrom{},
		MinPrice:       0,
		MaxPrice:       0,
		Params:         CoreSQLConfig.FieldsConfigType{},
		IsHistory:      false,
		Search:         "",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data, count)
	}
}

func TestCheckFinishByIDsIsFalse(t *testing.T) {
	checkList, err := CheckFinishByIDs(&ArgsCheckFinishByIDs{
		IDs: []int64{payData.ID},
	})
	if err != nil {
		t.Error(err, ", payData id: ", payData.ID, ", status: ", payData.Status)
	} else {
		if len(checkList) < 1 {
			return
		}
		if checkList[0].IsFinish {
			t.Error("check list is false")
		} else {
			t.Log(checkList)
			//加载数据，进行交叉检查
			TestGetOne(t)
		}
	}
}

func TestUpdateStatusClient(t *testing.T) {
	_, _, _, errCode, err := UpdateStatusClient(&ArgsUpdateStatusClient{
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
	})
	if err != nil {
		t.Error(errCode, err)
	}
}

func TestUpdateStatusFinish(t *testing.T) {
	errCode, err := UpdateStatusFinish(&ArgsUpdateStatusFinish{
		CreateInfo: CoreSQLFrom.FieldsFrom{
			System: "system",
			ID:     0,
			Mark:   "test",
			Name:   "",
		},
		ID:     payData.ID,
		Key:    "",
		Params: nil,
	})
	if err != nil {
		t.Error(errCode, err)
	}
}

func TestCheckFinishByIDsIsFinish(t *testing.T) {
	checkList, err := CheckFinishByIDs(&ArgsCheckFinishByIDs{
		IDs: []int64{payData.ID},
	})
	if err != nil {
		t.Error(err)
	} else {
		if len(checkList) < 1 {
			return
		}
		if checkList[0].IsFinish {
			t.Log(checkList)
		} else {
			t.Error("check list is true, checkList: ", checkList)
			//加载数据，进行交叉检查
			TestGetOne(t)
		}
	}
}

func TestUpdateStatusRefund(t *testing.T) {
	errCode, err := UpdateStatusRefund(&ArgsUpdateStatusRefund{
		CreateInfo: CoreSQLFrom.FieldsFrom{
			System: "system",
			ID:     0,
			Mark:   "test",
			Name:   "",
		},
		ID:     payData.ID,
		Key:    "",
		Params: nil,
	})
	if err != nil {
		t.Error(errCode, err)
	}
}

func TestUpdateStatusRefundAudit(t *testing.T) {
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
		RefundPrice: 10,
	})
	if err != nil {
		t.Error(errCode, err)
	}
}

func TestUpdateStatusRefundFinish(t *testing.T) {
	errCode, err := UpdateStatusRefundFinish(&ArgsUpdateStatusRefundFinish{
		CreateInfo: CoreSQLFrom.FieldsFrom{
			System: "system",
			ID:     0,
			Mark:   "test",
			Name:   "",
		},
		ID:     payData.ID,
		Key:    "",
		Params: nil,
	})
	if err != nil {
		t.Error(errCode, err)
	}
}

func TestUpdateStatusRemove(t *testing.T) {
	TestCreate(t)
	errCode, err := UpdateStatusRemove(&ArgsUpdateStatusRemove{
		CreateInfo: CoreSQLFrom.FieldsFrom{
			System: "system",
			ID:     0,
			Mark:   "test",
			Name:   "",
		},
		ID:     payData.ID,
		Key:    "",
		Params: nil,
	})
	if err != nil {
		t.Error(errCode, err)
	}
	TestCheckFinishByIDsIsFalse(t)
}

func TestUpdateStatusFailed(t *testing.T) {
	TestCreate(t)
	errCode, err := UpdateStatusFailed(&ArgsUpdateStatusFailed{
		CreateInfo: CoreSQLFrom.FieldsFrom{
			System: "system",
			ID:     0,
			Mark:   "test",
			Name:   "",
		},
		ID:            payData.ID,
		Key:           "",
		FailedCode:    "user-operate",
		FailedMessage: "公司拒绝收费测试",
		Params:        nil,
	})
	if err != nil {
		t.Error(errCode, err)
	}
	TestCheckFinishByIDsIsFalse(t)
}

func TestDepositClear2(t *testing.T) {
	TestDepositClear(t)
}
