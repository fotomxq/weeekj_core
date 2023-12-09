package BaseWeixinWXXMessage

import (
	"errors"
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	BaseWeixinWXXMessageTemplate "gitee.com/weeekj/weeekj_core/v5/base/weixin/wxx/message/template"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
)

// ArgsSendMessageTemplate 推送模版消息参数
type ArgsSendMessageTemplate struct {
	//商户ID
	// 可以留空，则走平台微信小程序主体
	MerchantID int64
	//组织ID
	OrgID int64
	//用户ID
	UserID int64
	//用户OpenID
	UserOpenID string
	//模版ID
	TemplateID string
	//页数
	Page string
	//表单ID
	FromID string
	//推送参数
	Data map[string]interface{}
	//关键词
	EmphasisKeyword string
}

// SendMessageTemplate 推送模版消息
func SendMessageTemplate(args *ArgsSendMessageTemplate) error {
	var err error
	//如果fromID不存在，则从数据库获取
	if args.FromID == "" {
		args.FromID, err = getByOpenID(args.OrgID, args.UserID, args.UserOpenID)
		if err != nil {
			return errors.New("cannot send weixin message, weixin form id is not exisit, open id: " + args.UserOpenID + ", template id: " + args.TemplateID)
		}
	}
	if args.TemplateID == "" {
		return errors.New("template id is empty")
	}
	bindData := BaseWeixinWXXMessageTemplate.Message{}
	for k, v := range args.Data {
		bindData[k] = v
	}
	err = BaseWeixinWXXMessageTemplate.Send(args.MerchantID, args.UserOpenID, args.TemplateID, args.Page, args.FromID, bindData, args.EmphasisKeyword)
	return err
}

// ArgsSendMessageTemplateByOrderCreate 推送订单创建成功参数
type ArgsSendMessageTemplateByOrderCreate struct {
	//商户ID
	// 可以留空，则走平台微信小程序主体
	MerchantID int64
	//组织ID
	OrgID int64
	//用户ID
	UserID int64
	//用户OpenID
	UserOpenID string
	//订单ID
	OrderID string
	//价格
	Price float64
	//姓名
	Name string
	//电话
	Phone string
	//地址
	Address string
}

// SendMessageTemplateByOrderCreate 推送订单创建成功
func SendMessageTemplateByOrderCreate(args *ArgsSendMessageTemplateByOrderCreate) error {
	WeixinXiaochengxuMessageOrderCreateON, err := BaseConfig.GetDataBool("WeixinXiaochengxuMessageOrderCreateON")
	if err != nil {
		return errors.New("cannot get config, " + err.Error())
	}
	WeixinXiaochengxuMessageOrderCreateTemplateID, err := BaseConfig.GetDataString("WeixinXiaochengxuMessageOrderCreateTemplateID")
	if err != nil {
		return errors.New("cannot get config, " + err.Error())
	}
	if !WeixinXiaochengxuMessageOrderCreateON {
		return nil
	}
	if args.UserOpenID == "" {
		return nil
	}
	data := map[string]interface{}{
		"keyword1": args.OrderID,
		"keyword2": CoreFilter.GetStringByFloat64(args.Price),
		"keyword3": args.Name,
		"keyword4": args.Phone,
		"keyword5": args.Address,
	}
	return SendMessageTemplate(&ArgsSendMessageTemplate{
		MerchantID:      args.MerchantID,
		OrgID:           args.OrgID,
		UserID:          args.UserID,
		UserOpenID:      args.UserOpenID,
		TemplateID:      WeixinXiaochengxuMessageOrderCreateTemplateID,
		Page:            "",
		FromID:          "",
		Data:            data,
		EmphasisKeyword: "keyword2",
	})
}

// ArgsSendMessageTemplateByOrderPay 推送订单支付成功参数
type ArgsSendMessageTemplateByOrderPay struct {
	//商户ID
	// 可以留空，则走平台微信小程序主体
	MerchantID int64
	//组织ID
	OrgID int64
	//用户ID
	UserID int64
	//用户OpenID
	UserOpenID string
	//订单ID
	OrderID string
	//支付来源
	PayFrom string
	//价格
	Price float64
	//订单备注
	OrderDes string
}

// SendMessageTemplateByOrderPay 推送订单支付成功
func SendMessageTemplateByOrderPay(args *ArgsSendMessageTemplateByOrderPay) error {
	WeixinXiaochengxuMessagePaySuccessON, err := BaseConfig.GetDataBool("WeixinXiaochengxuMessagePaySuccessON")
	if err != nil {
		return errors.New("cannot get config, " + err.Error())
	}
	WeixinXiaochengxuMessagePaySuccessTemplateID, err := BaseConfig.GetDataString("WeixinXiaochengxuMessagePaySuccessTemplateID")
	if err != nil {
		return errors.New("cannot get config, " + err.Error())
	}
	if !WeixinXiaochengxuMessagePaySuccessON {
		return nil
	}
	if args.UserOpenID == "" {
		return nil
	}
	data := map[string]interface{}{
		"keyword1": args.OrderID,
		"keyword2": args.PayFrom,
		"keyword3": CoreFilter.GetStringByFloat64(args.Price),
		"keyword4": args.OrderDes,
	}
	return SendMessageTemplate(&ArgsSendMessageTemplate{
		MerchantID:      args.MerchantID,
		OrgID:           args.OrgID,
		UserID:          args.UserID,
		UserOpenID:      args.UserOpenID,
		TemplateID:      WeixinXiaochengxuMessagePaySuccessTemplateID,
		Page:            "",
		FromID:          "",
		Data:            data,
		EmphasisKeyword: "keyword2",
	})
}
