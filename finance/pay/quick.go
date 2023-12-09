package FinancePay

import "sync"

// 锁定机制
var (
	createQuickPayLock          sync.Mutex
	createQuickPayAndClientLock sync.Mutex
)

// CreateQuickPay 快速创建并完成付款
// 可用于储蓄、现金转账处理，其他渠道请勿使用；本模块无法与第三方衔接
func CreateQuickPay(args *ArgsCreate) (payData FieldsPayType, errCode string, err error) {
	//锁定机制
	createQuickPayLock.Lock()
	defer createQuickPayLock.Unlock()
	//创建支付
	payData, errCode, err = Create(args)
	if err != nil {
		return
	}
	//确认支付
	_, _, _, errCode, err = UpdateStatusClient(&ArgsUpdateStatusClient{
		CreateInfo: args.PaymentCreate,
		ID:         payData.ID,
		Key:        "",
		Params:     nil,
		IP:         "0.0.0.0",
	})
	if err != nil {
		return
	}
	//完成支付
	errCode, err = UpdateStatusFinish(&ArgsUpdateStatusFinish{
		CreateInfo: args.PaymentCreate,
		ID:         payData.ID,
		Key:        "",
		Params:     nil,
	})
	if err != nil {
		return
	}
	//反馈
	return
}

// CreateQuickPayAndConfirm 快速完成支付和确认支付
func CreateQuickPayAndConfirm(args *ArgsCreate) (payData FieldsPayType, errCode string, err error) {
	//锁定机制
	createQuickPayAndClientLock.Lock()
	defer createQuickPayAndClientLock.Unlock()
	//创建支付
	payData, errCode, err = Create(args)
	if err != nil {
		return
	}
	//确认支付
	_, _, _, errCode, err = UpdateStatusClient(&ArgsUpdateStatusClient{
		CreateInfo: args.PaymentCreate,
		ID:         payData.ID,
		Key:        "",
		Params:     nil,
		IP:         "0.0.0.0",
	})
	if err != nil {
		return
	}
	//反馈
	return
}
