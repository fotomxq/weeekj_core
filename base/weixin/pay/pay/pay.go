package BaseWeixinPayPay

import (
	"encoding/json"
	"errors"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	BaseWeixinPayClient "github.com/fotomxq/weeekj_core/v5/base/weixin/pay/client"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	"time"
)

// ArgsPayCreate 支付请求模块参数
type ArgsPayCreate struct {
	//商户ID
	// 可以留空，则走平台微信小程序主体
	OrgID int64
	//支付金额
	Price int
	//用户微信OpenID
	UserOpenID string
	//描述
	Des string
	//支付key
	PayKey string
	//IP
	IP string
	//过期时间
	ExpireTime time.Time
}

// PayCreate 支付请求模块
// 二次封装支付类请求，避免接口变化
func PayCreate(args *ArgsPayCreate) (params Params, err error) {
	//隶属关系
	if args.OrgID > 0 {
		financePayOtherInOne, _ := BaseConfig.GetDataBool("FinancePayOtherInOne")
		if financePayOtherInOne {
			args.OrgID = 0
		}
	}
	//获取操作对象
	var client BaseWeixinPayClient.ClientType
	client, err = BaseWeixinPayClient.GetMerchantClient(0, args.OrgID)
	if err != nil {
		return
	}
	//获取反馈接口参数
	var payNotifyURL string
	payNotifyURL, err = getPayNotifyURL()
	if err != nil {
		err = errors.New("get config by WeixinXiaochengxuPayNotifyURL, " + err.Error())
		return
	}
	//获取是否允许信用卡参数
	var noCredit bool
	noCredit, err = getPayNoCredit()
	if err != nil {
		err = errors.New("get config by WeixinXiaochengxuPayNoCredit, " + err.Error())
		return
	}
	//组件form请求表单
	form := Order{
		// 必填
		AppID:      client.ConfigData.AppID,
		MchID:      client.ConfigData.MerchantID,
		TotalFee:   args.Price,
		NotifyURL:  payNotifyURL,
		OpenID:     args.UserOpenID,
		Body:       args.Des,
		OutTradeNo: args.PayKey,
		// 选填 ...
		IP:        args.IP,
		NoCredit:  noCredit,
		StartedAt: CoreFilter.GetNowTime(),
		ExpiredAt: args.ExpireTime,
		//订单优惠标记
		Tag: "",
		//商品详情
		Detail: "",
		//附加数据
		Attach: "",
	}
	var pres PaidResponse
	pres, err = form.Unify(&client)
	if err != nil {
		formLog, err2 := json.Marshal(form)
		if err2 != nil {
			// 致命性错误，叠加反馈
			err = errors.New("get weixin pres by form merchant key, " + err.Error() + ", params json err: " + err2.Error())
			return
		}
		err = errors.New("get weixin pres by form merchant key, " + err.Error() + ", params: " + string(formLog))
		return
	}
	// 获取小程序前点调用支付接口所需参数
	params, err = GetParams(&client, pres.NonceStr, pres.PrePayID)
	if err != nil {
		err = errors.New("get weixin payment params, " + err.Error())
		return
	}
	//反馈成功
	return params, nil
}

// ArgsPayRefund 向微信服务器发出退款申请参数
type ArgsPayRefund struct {
	//商户ID
	// 可以留空，则走平台微信小程序主体
	OrgID int64
	//总金额
	Price int
	//退款金额
	RefundPrice int
	//支付key
	PayKey string
}

// PayRefund 向微信服务器发出退款申请
func PayRefund(args *ArgsPayRefund) (Refunder, RefundedResponse, error) {
	//获取操作对象
	var client BaseWeixinPayClient.ClientType
	var err error
	client, err = BaseWeixinPayClient.GetMerchantClient(0, args.OrgID)
	if err != nil {
		return Refunder{}, RefundedResponse{}, err
	}
	//获取反馈接口参数
	payNotifyURL, err := getPayRefundNotifyURL()
	if err != nil {
		return Refunder{}, RefundedResponse{}, errors.New("get config by WeixinXiaochengxuPayNotifyURL, " + err.Error())
	}
	//生成请求参数
	refund := Refunder{
		AppID:         client.ConfigData.AppID,
		MchID:         client.ConfigData.MerchantID,
		TotalFee:      args.Price,
		RefundFee:     args.RefundPrice,
		TransactionID: "",
		OutTradeNo:    args.PayKey,
		OutRefundNo:   args.PayKey,
		RefundDesc:    "",
		NotifyURL:     payNotifyURL,
	}
	res, err := refund.Refund(&client)
	if err != nil {
		return Refunder{}, RefundedResponse{}, errors.New("post weixin refund, " + err.Error())
	}
	//反馈
	return refund, res, nil
}

// 获取支付反馈接口
func getPayNotifyURL() (string, error) {
	return BaseConfig.GetDataString("WeixinXiaochengxuPayNotifyURL")
}

// 获取退款反馈接口
func getPayRefundNotifyURL() (string, error) {
	return BaseConfig.GetDataString("WeixinXiaochengxuRefundNotifyURL")
}

// 是否允许信用卡支付
func getPayNoCredit() (bool, error) {
	return BaseConfig.GetDataBool("WeixinXiaochengxuPayNoCredit")
}
