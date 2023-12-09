package BaseWeixinPayNotify

import (
	"encoding/xml"
	"errors"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	FinancePay "gitee.com/weeekj/weeekj_core/v5/finance/pay"
)

// response 基础返回数据
type response struct {
	ReturnCode string `xml:"return_code"` // 返回状态码: SUCCESS/FAIL
	ReturnMsg  string `xml:"return_msg"`  // 返回信息: 返回信息，如非空，为错误原因
	ResultCode string `xml:"result_code"`
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`
}

// Check 检测返回信息是否包含错误
func (res response) Check() error {
	if res.ReturnCode != "SUCCESS" {
		return errors.New("交易失败: " + res.ReturnMsg)
	}
	if res.ResultCode != "SUCCESS" {
		return errors.New("发生错误: " + res.ErrCodeDes)
	}
	return nil
}

// PaidNotify 支付结果返回数据
type PaidNotify struct {
	AppID         string  `xml:"appid"`               // 小程序ID
	MchID         string  `xml:"mch_id"`              // 商户号
	TotalFee      int     `xml:"total_fee"`           // 标价金额
	NonceStr      string  `xml:"nonce_str"`           // 随机字符串
	Sign          string  `xml:"sign"`                // 签名
	SignType      string  `xml:"sign_type,omitempty"` // 签名类型: 目前支持HMAC-SHA256和MD5，默认为MD5
	OpenID        string  `xml:"openid"`
	TradeType     string  `xml:"trade_type"`                     // 交易类型 JSAPI
	Bank          string  `xml:"bank_type"`                      // 银行类型，采用字符串类型的银行标识
	Settlement    float64 `xml:"settlement_total_fee,omitempty"` // 应结订单金额=订单金额-非充值代金券金额，应结订单金额<=订单金额。
	FeeType       string  `xml:"fee_type,omitempty"`             // 货币种类: 符合ISO4217标准的三位字母代码，默认人民币: CNY
	CashFee       float64 `xml:"cash_fee"`                       // 现金支付金额订单的现金支付金额
	CashFeeType   string  `xml:"cash_fee_type,omitempty"`        // 现金支付货币类型: 符合ISO4217标准的三位字母代码，默认人民币: CNY
	CouponFee     float64 `xml:"coupon_fee,omitempty"`           // 总代金券金额: 代金券金额<=订单金额，订单金额-代金券金额=现金支付金额
	CouponCount   int     `xml:"coupon_count,omitempty"`         // 代金券使用数量
	TransactionID string  `xml:"transaction_id"`                 // 微信支付订单号
	Attach        string  `xml:"attach,omitempty"`               // 商家数据包，原样返回
	// 商户系统内部订单号: 要求32个字符内，只能是数字、大小写字母_-|*@ ，且在同一个商户号下唯一。
	OutTradeNo string `xml:"out_trade_no"`
	// 支付完成时间，格式为yyyyMMddHHmmss，如2009年12月25日9点10分10秒表示为20091225091010
	Timeend string `xml:"time_end"`
}

type paidNotify struct {
	response
	PaidNotify
}

// HandlePaidNotify 处理支付结果通知
type ArgsHandlePaidNotify struct {
	//通知内容主体
	BodyByte []byte
}

func HandlePaidNotify(args *ArgsHandlePaidNotify) ([]byte, error) {
	//解析数据
	var ntf paidNotify
	if err := xml.Unmarshal(args.BodyByte, &ntf); err != nil {
		return nil, err
	}
	if err := ntf.Check(); err != nil {
		return nil, err
	}

	//处理交易
	var failedCode string
	//生成params
	params := []CoreSQLConfig.FieldsConfigType{
		{
			Mark: "pay-TotalFee",
			Val:  CoreFilter.GetStringByInt(ntf.TotalFee),
		},
		{
			Mark: "pay-Bank",
			Val:  ntf.Bank,
		},
		{
			Mark: "pay-Settlement",
			Val:  CoreFilter.GetStringByFloat64(ntf.Settlement),
		},
		{
			Mark: "pay-FeeType",
			Val:  ntf.FeeType,
		},
		{
			Mark: "pay-CashFee",
			Val:  CoreFilter.GetStringByFloat64(ntf.CashFee),
		},
		{
			Mark: "pay-CashFeeType",
			Val:  ntf.CashFeeType,
		},
		{
			Mark: "pay-CouponFee",
			Val:  CoreFilter.GetStringByFloat64(ntf.CouponFee),
		},
		{
			Mark: "pay-CouponCount",
			Val:  CoreFilter.GetStringByInt(ntf.CouponCount),
		},
		{
			Mark: "pay-TransactionID",
			Val:  ntf.TransactionID,
		},
		{
			Mark: "pay-Timeend",
			Val:  ntf.Timeend,
		},
		{
			Mark: "pay-ntf",
			Val:  string(args.BodyByte),
		},
	}
	//获取数据集合，确定金额是否符合
	payData, err := FinancePay.GetOne(&FinancePay.ArgsGetOne{
		ID:  0,
		Key: ntf.OutTradeNo,
	})
	if err != nil {
		CoreLog.Error("weixin pay payment notify, get pay data by key, ", err)
	}
	// 获取手续费开关 如果有手续费则不需要验证价格
	NeedCommission := payData.Params.GetValNoBool("NeedCommission")
	if NeedCommission == "true" {
		actualAmountReceived := payData.Params.GetValInt64NoBool("ActualAmountReceived")
		if actualAmountReceived != int64(ntf.TotalFee) {
			CoreLog.Error("weixin pay payment notify, price not equal, ", err)
		}
	} else {
		if payData.Price > int64(ntf.TotalFee) {
			CoreLog.Error("weixin pay payment notify, price not equal, ", err)
		}
	}
	//服务器确定支付完成
	if errCode, err := FinancePay.UpdateStatusFinish(&FinancePay.ArgsUpdateStatusFinish{
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
	}
	//反馈结果
	isOK := failedCode == ""
	replay := newReplay(isOK, failedCode)
	resByte, err := xml.Marshal(replay)
	if err != nil {
		return nil, err
	}
	return resByte, err
}
