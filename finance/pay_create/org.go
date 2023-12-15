package FinancePayCreate

import (
	"errors"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	ClassConfig "github.com/fotomxq/weeekj_core/v5/class/config"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	FinancePay "github.com/fotomxq/weeekj_core/v5/finance/pay"
	FinancePayMod "github.com/fotomxq/weeekj_core/v5/finance/pay/mod"
	OrgCoreCore "github.com/fotomxq/weeekj_core/v5/org/core"
	"time"
)

// ArgsCreateUserToOrg 给商户付款创建支付请求参数
type ArgsCreateUserToOrg struct {
	//服务用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//是否为退款
	IsRefund bool `db:"is_refund" json:"isRefund" check:"bool"`
	//货币
	Currency int `db:"currency" json:"currency" check:"currency"`
	//单位报价
	// 服务负责人在收款前可以协商变更
	Price int64 `db:"price" json:"price" check:"price" empty:"true"`
	//支付方式
	// 如果为退单，则为付款方式
	PaymentChannel CoreSQLFrom.FieldsFrom `json:"paymentChannel"`
	//交易过期时间
	// 如果提交空的时间，将直接按照过期处理
	ExpireAt time.Time `db:"expire_at" json:"expireAt" check:"isoTime"`
	//备注
	Des string `db:"des" json:"des" check:"des" min:"1" max:"500" empty:"true"`
}

// CreateUserToOrg 给商户付款创建支付请求
func CreateUserToOrg(args *ArgsCreateUserToOrg) (payData FinancePay.FieldsPayType, errCode string, err error) {
	//修正orgID
	var fromInfoOrgID int64 = 0
	fromInfoOrgID = FinancePayMod.FixOrgID(args.OrgID)
	//检查储蓄
	var takeChannelMark string
	if fromInfoOrgID > 0 {
		takeChannelMark, err = OrgCoreCore.Config.GetConfigVal(&ClassConfig.ArgsGetConfig{
			BindID:    args.OrgID,
			Mark:      "FinanceDepositDefaultMark",
			VisitType: "admin",
		})
		if err != nil {
			errCode = "org_deposit_not_exist"
			err = errors.New("get org deposit mark config, " + err.Error())
			return
		}
	} else {
		takeChannelMark = BaseConfig.GetDataStringNoErr("FinancePayAllDeposit")
		if takeChannelMark == "" {
			errCode = "err_finance_pay_to"
			err = errors.New("get all deposit mark config, " + err.Error())
			return
		}
	}
	//是否为退款请求，将发起转账给用户的处理
	if args.IsRefund {
		payData, errCode, err = FinancePay.Create(&FinancePay.ArgsCreate{
			CreateInfo: CoreSQLFrom.FieldsFrom{
				System: "org",
				ID:     args.OrgID,
				Mark:   "",
				Name:   "",
			},
			PaymentCreate: CoreSQLFrom.FieldsFrom{
				System: "org",
				ID:     args.OrgID,
				Mark:   "",
				Name:   "",
			},
			PaymentChannel: CoreSQLFrom.FieldsFrom{
				System: "deposit",
				ID:     0,
				Mark:   takeChannelMark,
				Name:   "",
			},
			PaymentFrom: CoreSQLFrom.FieldsFrom{
				System: "org",
				ID:     fromInfoOrgID,
				Mark:   "",
				Name:   "",
			},
			TakeCreate: CoreSQLFrom.FieldsFrom{
				System: "user",
				ID:     args.UserID,
				Mark:   "",
				Name:   "",
			},
			TakeChannel: args.PaymentChannel,
			TakeFrom: CoreSQLFrom.FieldsFrom{
				System: "org",
				ID:     fromInfoOrgID,
				Mark:   "",
				Name:   "",
			},
			Des:      args.Des,
			ExpireAt: args.ExpireAt,
			Currency: args.Currency,
			Price:    args.Price,
			Params:   []CoreSQLConfig.FieldsConfigType{},
		})
	} else {
		//非退款模式下，正常发起支付
		payData, errCode, err = FinancePay.Create(&FinancePay.ArgsCreate{
			CreateInfo: CoreSQLFrom.FieldsFrom{
				System: "user",
				ID:     args.UserID,
				Mark:   "",
				Name:   "",
			},
			PaymentCreate: CoreSQLFrom.FieldsFrom{
				System: "user",
				ID:     args.UserID,
				Mark:   "",
				Name:   "",
			},
			PaymentChannel: args.PaymentChannel,
			PaymentFrom: CoreSQLFrom.FieldsFrom{
				System: "org",
				ID:     fromInfoOrgID,
				Mark:   "",
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
				ID:     0,
				Mark:   takeChannelMark,
				Name:   "",
			},
			TakeFrom: CoreSQLFrom.FieldsFrom{
				System: "org",
				ID:     fromInfoOrgID,
				Mark:   "",
				Name:   "",
			},
			Des:      args.Des,
			ExpireAt: args.ExpireAt,
			Currency: args.Currency,
			Price:    args.Price,
			Params:   []CoreSQLConfig.FieldsConfigType{},
		})
	}
	return
}
