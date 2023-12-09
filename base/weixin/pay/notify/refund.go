package BaseWeixinPayNotify

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	BaseWeixinPayClientCrypto "gitee.com/weeekj/weeekj_core/v5/base/weixin/pay/client/crypto"
	BaseWeixinWXXClient "gitee.com/weeekj/weeekj_core/v5/base/weixin/wxx/client"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	FinancePay "gitee.com/weeekj/weeekj_core/v5/finance/pay"
	"strings"
)

// 退款结果通知
type refundNotify struct {
	AppID      string `xml:"appid"`       // 小程序 APPID
	MchID      string `xml:"mch_id"`      // 商户号
	NonceStr   string `xml:"nonce_str"`   // 随机字符串
	Ciphertext string `xml:"req_info"`    // 加密信息
	ReturnCode string `xml:"return_code"` // 返回状态码: SUCCESS/FAIL
	ReturnMsg  string `xml:"return_msg"`  // 返回信息: 返回信息，如非空，为错误原因
}

// 检测返回信息是否包含错误
func checkRefundNotify(res refundNotify) error {
	if res.ReturnCode != "SUCCESS" {
		return errors.New("交易失败: " + res.ReturnMsg)
	}
	return nil
}

// RefundedNotify 解密后的退款通知消息体
type RefundedNotify struct {
	AppID         string // 小程序ID
	MchID         string // 商户号
	NonceStr      string // 随机字符串
	TransactionID string `xml:"transaction_id"` // 微信支付订单号
	// 商户系统内部订单号: 要求32个字符内，只能是数字、大小写字母_-|*@ ，且在同一个商户号下唯一。
	OutTradeNo  string  `xml:"out_trade_no"`
	RefundID    string  `xml:"refund_id"`     // 微信退款单号
	OutRefundNo string  `xml:"out_refund_no"` // 商户退款单号
	TotalFee    float64 `xml:"total_fee"`     // 标价金额
	// 当该订单有使用非充值券时，返回此字段。
	// 应结订单金额=订单金额-非充值代金券金额，应结订单金额<=订单金额。
	Settlement float64 `xml:"settlement_total_fee,omitempty"`
	RefundFee  float64 `xml:"refund_fee"` // 退款总金额,单位为分
	// 退款金额
	// 退款金额=申请退款金额-非充值代金券退款金额，退款金额<=申请退款金额
	SettlementRefund float64 `xml:"settlement_refund_fee"`
	// 退款状态
	// SUCCESS 退款成功 | CHANGE 退款异常 | REFUNDCLOSE 退款关闭
	RefundStatus string `xml:"refund_status"`
	// 退款成功时间
	// 资金退款至用户帐号的时间，格式2017-12-15 09:46:01
	SuccessTime string `xml:"success_time,omitempty"`
	// 退款入账账户:取当前退款单的退款入账方
	// 1）退回银行卡:  {银行名称}{卡类型}{卡尾号}
	// 2）退回支付用户零钱: 支付用户零钱
	// 3）退还商户: 商户基本账户 商户结算银行账户
	// 4）退回支付用户零钱通: 支付用户零钱通
	ReceiveAccount string `xml:"refund_recv_accout"`
	// 退款资金来源
	// REFUND_SOURCE_RECHARGE_FUNDS 可用余额退款/基本账户
	// REFUND_SOURCE_UNSETTLED_FUNDS 未结算资金退款
	RefundAccount string `xml:"refund_account"`
	// 退款发起来源
	// API接口
	// VENDOR_PLATFORM商户平台
	Source string `xml:"refund_request_source"`
}

// ArgsHandleRefundedNotify 处理退款结果通知参数
// key: 微信支付 KEY
type ArgsHandleRefundedNotify struct {
	//商户ID
	// 可以留空，则走平台微信小程序主体
	MerchantID int64
	//通知内容
	BodyByte []byte
}

// HandleRefundedNotify 处理退款结果通知
func HandleRefundedNotify(args *ArgsHandleRefundedNotify) ([]byte, error) {
	//获取操作对象
	client, err := BaseWeixinWXXClient.GetMerchantClient(args.MerchantID)
	if err != nil {
		return nil, err
	}
	//解析数据集合
	var ref refundNotify
	if err := xml.Unmarshal(args.BodyByte, &ref); err != nil {
		err = errors.New("get xml un data, " + err.Error())
		return nil, err
	}
	if err := checkRefundNotify(ref); err != nil {
		err = errors.New("check return notify, " + err.Error())
		return nil, err
	}
	ciphertext, err := base64.StdEncoding.DecodeString(ref.Ciphertext)
	if err != nil {
		err = errors.New("get std encoding decode string, " + err.Error())
		return nil, err
	}
	key, err := BaseWeixinPayClientCrypto.MD5(client.ConfigData.Key)
	if err != nil {
		err = errors.New("get base weixin pay client cry, md5, " + err.Error())
		return nil, err
	}
	key = strings.ToLower(key)
	ntfByte, err := BaseWeixinPayClientCrypto.AesECBDecrypt(ciphertext, []byte(key))
	if err != nil {
		err = errors.New("get base weixin pay client cry, aes, " + err.Error())
		return nil, err
	}
	ntf := RefundedNotify{
		AppID:            ref.AppID,
		MchID:            ref.MchID,
		NonceStr:         ref.NonceStr,
		TransactionID:    "",
		OutTradeNo:       "",
		RefundID:         "",
		OutRefundNo:      "",
		TotalFee:         0,
		Settlement:       0,
		RefundFee:        0,
		SettlementRefund: 0,
		RefundStatus:     "",
		SuccessTime:      "",
		ReceiveAccount:   "",
		RefundAccount:    "",
		Source:           "",
	}
	if err := xml.Unmarshal(ntfByte, &ntf); err != nil {
		err = errors.New("get xml un, " + err.Error())
		return nil, err
	}
	//处理交易数据
	//生成交易日志
	CoreLog.Info("weixin pay refund notify, ", ntfByte)
	//退款参数集合
	params := []CoreSQLConfig.FieldsConfigType{
		{
			Mark: "refund-ntf",
			Val:  string(ntfByte),
		},
	}
	var failedCode string
	if ntf.RefundStatus == "SUCCESS" {
		//服务器确定支付完成
		if errCode, err := FinancePay.UpdateStatusRefundFinish(&FinancePay.ArgsUpdateStatusRefundFinish{
			CreateInfo: CoreSQLFrom.FieldsFrom{
				System: "weixin",
				ID:     0,
				Mark:   ntf.TransactionID,
				Name:   "",
			},
			ID:     0,
			Key:    ntf.OutTradeNo,
			Params: params,
		}); err != nil {
			CoreLog.Error("weixin pay payment notify, err: ", err)
			failedCode = errCode
		} else {
			CoreLog.Info("weixin refund notify success, pay trade no: ", ntf.OutTradeNo)
		}
	} else {
		//交易失败
		if errCode, err := FinancePay.UpdateStatusFailed(&FinancePay.ArgsUpdateStatusFailed{
			CreateInfo: CoreSQLFrom.FieldsFrom{
				System: "weixin",
				ID:     0,
				Mark:   ntf.TransactionID,
				Name:   "",
			},
			ID:            0,
			Key:           ntf.OutTradeNo,
			FailedCode:    "refund-weixin",
			FailedMessage: ntf.RefundStatus,
			Params:        nil,
		}); err != nil {
			CoreLog.Error("weixin pay payment notify, refund failed and update failed, pay trade no: ", ntf.OutTradeNo, ", err: ", err)
			failedCode = errCode
		} else {
			CoreLog.Error("weixin pay payment notify, refund failed, ", ntf.RefundStatus, ", pay trade no: ", ntf.OutTradeNo)
		}
	}
	//反馈数据
	isOK := failedCode == ""
	pr := newReplay(isOK, failedCode)
	resByte, err := xml.Marshal(pr)
	if err != nil {
		err = errors.New("res marshal, " + err.Error())
		return nil, err
	}
	return resByte, err
}
