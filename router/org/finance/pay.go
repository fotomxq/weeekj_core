package RouterOrgFinance

import (
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	FinanceDeposit "github.com/fotomxq/weeekj_core/v5/finance/deposit"
	FinancePay "github.com/fotomxq/weeekj_core/v5/finance/pay"
	RouterFinance "github.com/fotomxq/weeekj_core/v5/router/finance"
	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
	"github.com/gin-gonic/gin"
	"time"
)

// ArgsCreatePayToOrg 发起给组织付款的请求参数
type ArgsCreatePayToOrg struct {
	//付款渠道系统
	PaySystem string `json:"paySystem,omitempty"`
	//付款标识码
	SaveMark string `json:"saveMark,omitempty"`
	//组织ID
	OrgID int64 `json:"orgID,omitempty"`
	//货币
	Currency int `json:"currency,omitempty"`
	//金额
	Price int64 `json:"price,omitempty"`
	//备注
	Des string `json:"des,omitempty"`
}

// CreatePayToOrg 发起给组织付款的请求
func CreatePayToOrg(c *gin.Context, userData *UserCore.DataUserDataType, args *ArgsCreatePayToOrg) (payData FinancePay.FieldsPayType, failedCode string, failedMsg string, err error) {
	//如果付款方为储蓄，则获取储蓄数据
	var depositID int64
	if args.PaySystem == "deposit" {
		var depositData FinanceDeposit.FieldsDepositType
		depositData, failedCode, err = RouterFinance.DepositGetByUser(userData, CoreSQLFrom.FieldsFrom{
			System: "org",
			ID:     args.OrgID,
			Mark:   "",
			Name:   "",
		}, args.SaveMark)
		if err != nil {
			failedCode = "user-not-deposit"
			failedMsg = ""
			return
		}
		depositID = depositData.ID
	}
	//获取组织的收款账户信息
	var orgDepositData FinanceDeposit.FieldsDepositType
	var defaultDepositMark string
	orgDepositData, defaultDepositMark, err = GetDepositDataAndDefaultMark(args.OrgID)
	if err != nil {
		failedCode = "org-not-deposit"
		failedMsg = ""
		return
	}
	//创建交易请求
	var b bool
	payData, b = RouterFinance.PayCreate(c, userData, &RouterFinance.ArgsPayCreate{
		PaymentChannel: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     userData.Info.ID,
			Mark:   "",
			Name:   userData.Info.Name,
		},
		PaymentFrom: CoreSQLFrom.FieldsFrom{
			System: args.PaySystem,
			ID:     depositID,
			Mark:   args.SaveMark,
			Name:   "",
		},
		TakeCreate: CoreSQLFrom.FieldsFrom{
			System: "org",
			ID:     args.OrgID,
			Mark:   "",
			Name:   "",
		},
		TakeChannel: CoreSQLFrom.FieldsFrom{
			System: "deposit",
			ID:     orgDepositData.ID,
			Mark:   defaultDepositMark,
			Name:   "",
		},
		TakeFrom:   CoreSQLFrom.FieldsFrom{},
		Des:        args.Des,
		ExpireTime: time.Time{},
		Currency:   args.Currency,
		Price:      args.Price,
		Params:     nil,
	})
	if !b {
		failedCode = "pay-failed"
		failedMsg = ""
		return
	}
	//反馈成功
	return
}
