package RouterOrgFinance

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	FinancePay "gitee.com/weeekj/weeekj_core/v5/finance/pay"
)

// ArgsPayCreateOrgToUser 商户给用户付款参数
type ArgsPayCreateOrgToUser struct {
	//商户ID
	OrgID int64 `json:"orgID"`
	//商户名称
	OrgName string `json:"orgName"`
	//目标用户ID
	UserID int64 `json:"userID"`
	//目标用户昵称
	UserName string `json:"userName"`
	//收款渠道
	// system: cash 现金 ; deposit 存储模块 ; weixin 微信支付 ; alipay 支付宝
	// mark: 子渠道信息，例如 weixin 的wxx/merchant
	TakeChannel CoreSQLFrom.FieldsFrom `db:"take_channel" json:"takeChannel"`
	//交易备注
	Des string `json:"des"`
	//货币
	Currency int `json:"currency"`
	//金额
	Price int64 `json:"price"`
	//扩展信息
	Params []CoreSQLConfig.FieldsConfigType `json:"params"`
}

// PayCreateOrgToUser 商户给用户付款
func PayCreateOrgToUser(args *ArgsPayCreateOrgToUser) (payData FinancePay.FieldsPayType, errCode string, err error) {
	var configMark string
	_, configMark, err = GetDepositDataAndDefaultMark(args.OrgID)
	if err != nil {
		errCode = "no_default_mark"
		return
	}
	payData, errCode, err = FinancePay.Create(&FinancePay.ArgsCreate{
		CreateInfo: CoreSQLFrom.FieldsFrom{
			System: "org",
			ID:     args.OrgID,
			Mark:   "",
			Name:   args.OrgName,
		},
		PaymentCreate: CoreSQLFrom.FieldsFrom{
			System: "org",
			ID:     args.OrgID,
			Mark:   "",
			Name:   args.OrgName,
		},
		PaymentChannel: CoreSQLFrom.FieldsFrom{
			System: "deposit",
			ID:     0,
			Mark:   "",
			Name:   configMark,
		},
		PaymentFrom: CoreSQLFrom.FieldsFrom{
			System: "org",
			ID:     args.OrgID,
			Mark:   "",
			Name:   args.OrgName,
		},
		TakeCreate: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     args.UserID,
			Mark:   "",
			Name:   args.UserName,
		},
		TakeChannel: args.TakeChannel,
		TakeFrom: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     args.UserID,
			Mark:   "",
			Name:   args.UserName,
		},
		Des:      "配送单缴费",
		ExpireAt: CoreFilter.GetNowTimeCarbon().AddMinutes(30).Time,
		Currency: args.Currency,
		Price:    args.Price,
		Params:   args.Params,
	})
	return
}
