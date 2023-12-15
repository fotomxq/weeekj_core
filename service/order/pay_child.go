package ServiceOrder

import (
	"errors"
	"fmt"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	FinancePay "github.com/fotomxq/weeekj_core/v5/finance/pay"
	FinancePayCreate "github.com/fotomxq/weeekj_core/v5/finance/pay_create"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsCreateChildPay 支付订单子项目参数
type ArgsCreateChildPay struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 可选，作为验证
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	// 可选，作为验证
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//日志
	//操作组织人员ID
	// 如果留空则说明为系统自动调整或创建人产生
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//支付方式
	// system: cash 现金 ; deposit 存储模块 ; weixin 微信支付 ; alipay 支付宝
	// mark: 子渠道信息，例如 weixin 的wxx/merchant
	PaymentChannel CoreSQLFrom.FieldsFrom `db:"payment_channel" json:"paymentChannel"`
	//支付备注
	// 用户环节可根据实际业务需求开放此项
	Des string `db:"des" json:"des" check:"des" min:"1" max:"600" empty:"true"`
	//要支付的priceList.priceType
	PriceType int `db:"price_type" json:"priceType"`
}

// CreateChildPay 支付订单子项目参数
func CreateChildPay(args *ArgsCreateChildPay) (data FinancePay.FieldsPayType, errCode string, err error) {
	//获取订单
	var orderData FieldsOrder
	orderData, err = GetByID(&ArgsGetByID{
		ID:     args.ID,
		OrgID:  args.OrgID,
		UserID: -1,
	})
	if err != nil {
		errCode = "order_not_exist"
		return
	}
	//检查是否已经支付?
	// 同时抽取出费用
	var price int64
	for _, v := range orderData.PriceList {
		if v.PriceType != args.PriceType {
			continue
		}
		if v.IsPay {
			errCode = "is_pay"
			err = errors.New("price type is pay")
			return
		}
		price = v.Price
	}
	//构建支付请求
	data, errCode, err = FinancePayCreate.CreateUserToOrg(&FinancePayCreate.ArgsCreateUserToOrg{
		UserID:         orderData.UserID,
		OrgID:          orderData.OrgID,
		IsRefund:       false,
		Currency:       orderData.Currency,
		Price:          price,
		PaymentChannel: args.PaymentChannel,
		ExpireAt:       orderData.ExpireAt,
		Des:            args.Des,
	})
	if err != nil {
		return
	}
	//修改数据
	for k, v := range orderData.PriceList {
		if v.PriceType != args.PriceType {
			continue
		}
		orderData.PriceList[k].PayID = data.ID
		orderData.PriceList[k].PayFailed = ""
	}
	//修改上述所有订单的支付ID
	var newLog string
	newLog, err = getLogData(args.UserID, args.OrgBindID, "pay_price_list_create", args.Des)
	if err != nil {
		errCode = "order_update_price_list_pay_id"
		return
	}
	if _, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_order SET price_list = :price_list, logs = logs || :log WHERE id = :id", map[string]interface{}{
		"id":       orderData.ID,
		"pay_list": orderData.PriceList,
		"log":      newLog,
	}); err != nil {
		errCode = "order_update_pay_id"
		err = errors.New(fmt.Sprint("order update pay id, ", err))
		return
	}
	//清理缓冲
	deleteOrderCache(orderData.ID)
	//反馈
	return
}
