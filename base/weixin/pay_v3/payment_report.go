package WeixinPayV3

import (
	"context"
	"errors"
	"fmt"
	FinancePayMod "gitee.com/weeekj/weeekj_core/v5/finance/pay/mod"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"net/http"
)

// ArgsPaymentReport 支付回调处理参数
type ArgsPaymentReport struct {
	//组织ID
	OrgID int64 `json:"orgID"`
}

// PaymentReport 支付回调处理
func PaymentReport(args *ArgsPaymentReport, request *http.Request) (*payments.Transaction, error) {
	//隶属关系
	args.OrgID = FinancePayMod.FixOrgID(args.OrgID)
	//获取证书
	_, clientConfig, err := getClient(args.OrgID)
	if err != nil {
		return &payments.Transaction{}, err
	}
	//获取平台证书访问器
	certVisitor := downloader.MgrInstance().GetCertificateVisitor(clientConfig.MerchantID)
	handler := notify.NewNotifyHandler(clientConfig.KeyV3, verifiers.NewSHA256WithRSAVerifier(certVisitor))
	//解密数据包
	transaction := new(payments.Transaction)
	notifyReq, err := handler.ParseNotifyRequest(context.Background(), request, transaction)
	// 如果验签未通过，或者解密失败
	if err != nil {
		err = errors.New(fmt.Sprint("parse notify request, org id: ", args.OrgID, ", merchant id: ", clientConfig.MerchantID, ", cert sn: ", clientConfig.CertSN, ", err: ", err))
		return &payments.Transaction{}, err
	}
	if notifyReq.EventType == "TRANSACTION.SUCCESS" {
		return transaction, nil
	}
	return nil, errors.New(fmt.Sprint("failed, ", notifyReq))
}
