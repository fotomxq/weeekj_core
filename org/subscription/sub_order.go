package OrgSubscription

import (
	"fmt"
	CoreSQLAddress "gitee.com/weeekj/weeekj_core/v5/core/sql/address"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	ServiceOrderWait "gitee.com/weeekj/weeekj_core/v5/service/order/wait"
	ServiceOrderWaitFields "gitee.com/weeekj/weeekj_core/v5/service/order/wait_fields"
	"time"
)

// ArgsCreateSubOrder 创建新的订阅请求参数
type ArgsCreateSubOrder struct {
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//创建来源
	// 0 商户平台下单、1 用户APP下单、2 微信小程序下单、3 支付宝小程序下单、4 电话系统下单、5 线下柜台下单、6 辅助设备下单、7 其他渠道
	CreateFrom int `db:"create_from" json:"createFrom"`
	//收取货物地址
	AddressFrom CoreSQLAddress.FieldsAddress `db:"address_from" json:"addressFrom"`
	//送货地址
	AddressTo CoreSQLAddress.FieldsAddress `db:"address_to" json:"addressTo"`
	//备注信息
	Des string `db:"des" json:"des" check:"des" min:"1" max:"600" empty:"true"`
	//订阅配置ID
	SubConfigID int64 `db:"sub_config_id" json:"subConfigID" check:"id"`
	//购买的单位
	Unit int64 `db:"unit" json:"unit" check:"int64Than0"`
}

// CreateSubOrder 创建新的订阅请求
// 用户专用请求
func CreateSubOrder(args *ArgsCreateSubOrder) (data ServiceOrderWaitFields.FieldsWait, errCode string, err error) {
	var configData FieldsConfig
	configData, err = GetConfigByID(&ArgsGetConfigByID{
		ID: args.SubConfigID,
	})
	if err != nil {
		errCode = "sub_config_not_exist"
		return
	}
	data, errCode, err = ServiceOrderWait.CreateOrder(&ServiceOrderWait.ArgsCreateOrder{
		SystemMark:  "org_sub",
		OrgID:       args.OrgID,
		UserID:      args.UserID,
		CreateFrom:  args.CreateFrom,
		AddressFrom: args.AddressFrom,
		AddressTo:   args.AddressTo,
		Goods: []ServiceOrderWaitFields.FieldsGood{
			{
				From: CoreSQLFrom.FieldsFrom{
					System: "org_sub",
					ID:     configData.ID,
					Mark:   "virtual",
					Name:   configData.Title,
				},
				OptionKey:  "",
				Count:      args.Unit,
				Price:      configData.Price,
				Exemptions: []ServiceOrderWaitFields.FieldsExemption{},
			},
		},
		Exemptions:         []ServiceOrderWaitFields.FieldsExemption{},
		NeedAllowAutoAudit: true,
		AllowAutoAudit:     true,
		TransportAllowAuto: false,
		TransportTaskAt:    time.Time{},
		TransportPayAfter:  false,
		PriceList: []ServiceOrderWaitFields.FieldsPrice{
			{
				PriceType: 0,
				IsPay:     false,
				Price:     args.Unit * configData.Price,
			},
		},
		PricePay:    false,
		NeedExPrice: false,
		Currency:    configData.Currency,
		Des:         args.Des,
		Logs:        []ServiceOrderWaitFields.FieldsLog{},
		Params: []CoreSQLConfig.FieldsConfigType{
			{
				Mark: "org_id",
				Val:  fmt.Sprint(args.OrgID),
			},
		},
	})
	if err != nil {
		return
	}
	return
}
