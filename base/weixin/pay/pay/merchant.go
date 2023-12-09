package BaseWeixinPayPay

import (
	"encoding/json"
	"errors"
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	BaseWeixinPayClient "gitee.com/weeekj/weeekj_core/v5/base/weixin/pay/client"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
)

// ArgsMerchantChange 商户能力支持参数
type ArgsMerchantChange struct {
	//商户ID
	// 可以留空，则走平台微信小程序主体
	OrgID int64
	//支付key
	PayKey string
	//用户微信OpenID
	UserOpenID string
	//用户昵称
	UserName string
	//支付描述
	PayDes string
	//支付金额
	Price int
}

// MerchantChange 商户能力支持
func MerchantChange(args *ArgsMerchantChange) ([]byte, error) {
	//隶属关系
	if args.OrgID > 0 {
		financePayOtherInOne, _ := BaseConfig.GetDataBool("FinancePayOtherInOne")
		if financePayOtherInOne {
			args.OrgID = 0
		}
	}
	//获取操作对象
	var client BaseWeixinPayClient.ClientType
	var err error
	client, err = BaseWeixinPayClient.GetMerchantClient(0, args.OrgID)
	if err != nil {
		return nil, err
	}
	//获取是否需要验证
	needCheckUserName, err := BaseConfig.GetDataBool("WeixinXiaochengxuPayMerchantNeedCheckName")
	if err != nil {
		return []byte{}, errors.New("get config by WeixinXiaochengxuPayMerchantNeedCheckName, " + err.Error())
	}
	//验证用户姓名
	var checkName string
	checkMode := false
	if needCheckUserName {
		checkMode = true
		checkName = args.UserName
	}
	// 新建退款订单
	form := Transferer{
		// 必填 ...
		AppID:      client.ConfigData.AppID,
		MchID:      client.ConfigData.MerchantID,
		Amount:     args.Price,
		OutTradeNo: args.PayKey, // or TransactionID: "微信订单号",
		ToUser:     args.UserOpenID,
		Desc:       args.PayDes, // 若商户传入, 会在下发给用户的退款消息中体现退款原因
		// 选填 ...
		IP:        "", // 若商户传入, 会在下发给用户的退款消息中体现退款原因
		CheckName: checkMode,
		RealName:  checkName, // 如果 CheckName 设置为 true 则必填用户真实姓名
		Device:    "",
	}
	CoreLog.Info("weixin merchant change, ", form)
	// 需要证书
	res, err := form.Transfer(&client)
	if err != nil {
		return nil, err
	}
	resJSON, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	return resJSON, nil
}
