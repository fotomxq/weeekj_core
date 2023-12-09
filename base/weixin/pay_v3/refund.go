package WeixinPayV3

import (
	"context"
	"errors"
	"fmt"
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	FinancePayMod "gitee.com/weeekj/weeekj_core/v5/finance/pay/mod"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
)

// ArgsCreateRefund 创建退款请求参数
type ArgsCreateRefund struct {
	//组织ID
	OrgID int64 `json:"orgID"`
	//订单描述
	Des string `json:"des"`
	//支付Key
	PayKey string `json:"payKey"`
	//退款单号ID
	RefundKey string `json:"refundKey"`
	//微信支付ID
	TransactionId string `json:"transactionId"`
	//要退款的金额
	PriceRefund int64 `json:"priceRefund"`
	//订单总的金额
	PriceTotal int64 `json:"priceTotal"`
}

// CreateRefund 创建退款请求
func CreateRefund(args *ArgsCreateRefund) (CoreSQLConfig.FieldsConfigsType, error) {
	ctx := context.Background()
	//隶属关系
	args.OrgID = FinancePayMod.FixOrgID(args.OrgID)
	//构建client
	client, _, err := getClient(args.OrgID)
	if err != nil {
		err = errors.New(fmt.Sprint("get org client, ", err))
		return CoreSQLConfig.FieldsConfigsType{}, err
	}
	//获取反馈通知接口
	notifyUrl, err := BaseConfig.GetDataString("AppAPI")
	if err != nil {
		err = errors.New(fmt.Sprint("get AppAPI config, ", err))
		return CoreSQLConfig.FieldsConfigsType{}, err
	}
	notifyUrl = fmt.Sprint(notifyUrl, "/v2/base/weixin/public/refund/v3/notify/", args.OrgID)
	//构建退款支付请求
	svc := refunddomestic.RefundsApiService{Client: client}
	resp, _, err := svc.Create(ctx,
		refunddomestic.CreateRequest{
			SubMchid:      nil,
			TransactionId: core.String(args.TransactionId),
			OutTradeNo:    core.String(args.PayKey),
			OutRefundNo:   core.String(args.RefundKey),
			Reason:        core.String(args.Des),
			NotifyUrl:     core.String(notifyUrl),
			FundsAccount:  nil,
			Amount: &refunddomestic.AmountReq{
				Currency: core.String("CNY"),
				From:     nil,
				Refund:   core.Int64(args.PriceRefund),
				Total:    core.Int64(args.PriceTotal),
			},
			GoodsDetail: nil,
		},
	)
	if err != nil {
		err = errors.New(fmt.Sprint("create refund post, ", err))
		return CoreSQLConfig.FieldsConfigsType{}, err
	}
	//反馈数据集合
	refundId := CoreFilter.DerefString(resp.RefundId)
	amountRefund := CoreFilter.DerefInt64(resp.Amount.Refund)
	return CoreSQLConfig.FieldsConfigsType{
		{
			Mark: "refundId",
			Val:  refundId,
		},
		{
			Mark: "amountRefund",
			Val:  fmt.Sprint(amountRefund),
		},
	}, nil
}
