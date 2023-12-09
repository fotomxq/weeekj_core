package WeixinPayV3

import (
	"context"
	"errors"
	"fmt"
	FinancePayMod "gitee.com/weeekj/weeekj_core/v5/finance/pay/mod"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
	"net/http"
)

// ArgsRefundReport 退款回调处理函数
type ArgsRefundReport struct {
	//组织ID
	OrgID int64 `json:"orgID"`
}

type DataRefundReport struct {
	Mchid               *string                 `json:"mchid,omitempty"`
	TransactionId       *string                 `json:"transaction_id,omitempty"`
	OutTradeNo          *string                 `json:"out_trade_no,omitempty"`
	RefundID            *string                 `json:"refund_id,omitempty"`
	OutRefundNo         *string                 `json:"out_refund_no,omitempty"`
	RefundStatus        *string                 `json:"refund_status,omitempty"`
	SuccessTime         *string                 `json:"success_time,omitempty"`
	UserReceivedAccount *string                 `json:"user_received_account,omitempty"`
	Amount              *DataRefundReportAmount `json:"amount,omitempty"`
}

type DataRefundReportAmount struct {
	Total       *int64 `json:"total,omitempty"`
	Refund      *int64 `json:"refund,omitempty"`
	PayerTotal  *int64 `json:"payer_total,omitempty"`
	PayerRefund *int64 `json:"payer_refund,omitempty"`
}

// RefundReport 退款回调处理
func RefundReport(args *ArgsRefundReport, request *http.Request) (*DataRefundReport, error) {
	//隶属关系
	args.OrgID = FinancePayMod.FixOrgID(args.OrgID)
	//获取证书
	_, clientConfig, err := getClient(args.OrgID)
	if err != nil {
		return &DataRefundReport{}, err
	}
	//获取平台证书访问器
	certVisitor := downloader.MgrInstance().GetCertificateVisitor(clientConfig.MerchantID)
	handler := notify.NewNotifyHandler(clientConfig.KeyV3, verifiers.NewSHA256WithRSAVerifier(certVisitor))
	//解密数据包
	transaction := new(DataRefundReport)
	notifyReq, err := handler.ParseNotifyRequest(context.Background(), request, transaction)
	// 如果验签未通过，或者解密失败
	if err != nil {
		err = errors.New(fmt.Sprint("parse notify request, org id: ", args.OrgID, ", merchant id: ", clientConfig.MerchantID, ", cert sn: ", clientConfig.CertSN, ", err: ", err))
		return &DataRefundReport{}, err
	}
	if notifyReq.EventType == "REFUND.SUCCESS" {
		return transaction, nil
	}
	return nil, errors.New(fmt.Sprint("failed, ", notifyReq))
}

// ArgsRefundCheck 检查退款是否完成参数
type ArgsRefundCheck struct {
	//组织ID
	OrgID int64 `json:"orgID"`
	//支付key
	PayKey string `json:"payKey,omitempty"`
}

// RefundCheck 检查退款是否完成
func RefundCheck(args *ArgsRefundCheck) (*refunddomestic.Refund, error) {
	ctx := context.Background()
	//隶属关系
	args.OrgID = FinancePayMod.FixOrgID(args.OrgID)
	//构建client
	client, _, err := getClient(args.OrgID)
	if err != nil {
		return &refunddomestic.Refund{}, err
	}
	//处理反馈
	svc := refunddomestic.RefundsApiService{Client: client}
	resp, _, err := svc.QueryByOutRefundNo(ctx,
		refunddomestic.QueryByOutRefundNoRequest{
			OutRefundNo: core.String(args.PayKey),
			SubMchid:    nil,
		},
	)
	if err != nil {
		err = errors.New(fmt.Sprint("get report data, ", err))
		return &refunddomestic.Refund{}, err
	}
	return resp, nil
}
