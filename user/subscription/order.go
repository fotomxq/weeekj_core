package UserSubscription

import (
	"errors"
	"fmt"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	ClassConfig "github.com/fotomxq/weeekj_core/v5/class/config"
	CoreSQLAddress "github.com/fotomxq/weeekj_core/v5/core/sql/address"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	OrgCore "github.com/fotomxq/weeekj_core/v5/org/core"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	ServiceOrderWait "github.com/fotomxq/weeekj_core/v5/service/order/wait"
	ServiceOrderWaitFields "github.com/fotomxq/weeekj_core/v5/service/order/wait_fields"
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
	//推荐人手机号
	ReferrerNationCode string `db:"referrer_nation_code" json:"referrerNationCode" check:"nationCode" empty:"true"`
	ReferrerPhone      string `json:"referrerPhone" check:"phone" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
	//强制约定价格
	UnitPrice int64 `db:"unit_price" json:"unitPrice"`
}

// CreateSubOrder 创建新订阅请求
// 用户专用请求
func CreateSubOrder(args *ArgsCreateSubOrder) (data ServiceOrderWaitFields.FieldsWait, errCode string, err error) {
	//限制购买数量
	if args.Unit < 1 && args.Unit > 9999 {
		errCode = "unit_limit"
		err = errors.New("unit too many or less")
		return
	}
	//获取配置
	var configData FieldsConfig
	configData, err = GetConfigByID(&ArgsGetConfigByID{
		ID: args.SubConfigID,
	})
	if err != nil {
		errCode = "sub_config_not_exist"
		return
	}
	if args.UnitPrice > 0 {
		configData.Price = args.UnitPrice
	}
	//检查平台是否启用了锁定机制
	var userSubLevelLockGlob bool
	userSubLevelLockGlob, err = BaseConfig.GetDataBool("UserSubscriptionLevelLock")
	if err != nil {
		userSubLevelLockGlob = true
	}
	//获取商户配置，订阅锁定机制
	var userSubLevelLock bool
	if userSubLevelLockGlob {
		userSubLevelLock, err = OrgCore.Config.GetConfigValBool(&ClassConfig.ArgsGetConfig{
			BindID:    args.OrgID,
			Mark:      "UserSubLevelLock",
			VisitType: "admin",
		})
		if err != nil {
			userSubLevelLock = true
			err = nil
		}
	}
	if userSubLevelLock {
		//获取用户当前在续订阅情况
		var subList []FieldsSub
		err = Router2SystemConfig.MainDB.Select(&subList, "SELECT id, config_id FROM user_sub WHERE expire_at > NOW() AND delete_at < to_timestamp(1000000) AND org_id = $1 AND user_id = $2", args.OrgID, args.UserID)
		//检查如果存在价格更高的，则拒绝生成订单
		for _, v := range subList {
			var vConfig FieldsConfig
			vConfig, err = GetConfigByID(&ArgsGetConfigByID{
				ID:    v.ConfigID,
				OrgID: args.OrgID,
			})
			if err != nil {
				err = nil
				continue
			}
			if vConfig.Price > configData.Price {
				errCode = "level_lock"
				err = errors.New("user have high level user sub")
				return
			}
		}
		err = nil
	}
	//费用计算
	price := configData.Price
	subPrice := args.Unit * price
	for _, v := range configData.ExemptionTime {
		if int64(v.TimeN) == args.Unit {
			price = v.Price / args.Unit
			subPrice = v.Price
		}
	}
	//创建订单
	data, errCode, err = ServiceOrderWait.CreateOrder(&ServiceOrderWait.ArgsCreateOrder{
		SystemMark:  "user_sub",
		OrgID:       args.OrgID,
		UserID:      args.UserID,
		CreateFrom:  args.CreateFrom,
		AddressFrom: args.AddressFrom,
		AddressTo:   args.AddressTo,
		Goods: []ServiceOrderWaitFields.FieldsGood{
			{
				From: CoreSQLFrom.FieldsFrom{
					System: "user_sub",
					ID:     configData.ID,
					Mark:   "virtual",
					Name:   configData.Title,
				},
				OptionKey:  "",
				Count:      args.Unit,
				Price:      price,
				Exemptions: []ServiceOrderWaitFields.FieldsExemption{},
			},
		},
		Exemptions:         []ServiceOrderWaitFields.FieldsExemption{},
		NeedAllowAutoAudit: true,
		AllowAutoAudit:     true,
		TransportAllowAuto: false,
		TransportTaskAt:    time.Time{},
		TransportPayAfter:  false,
		TransportSystem:    "",
		PriceList: []ServiceOrderWaitFields.FieldsPrice{
			{
				PriceType: 0,
				IsPay:     false,
				Price:     subPrice,
			},
		},
		PricePay:           false,
		NeedExPrice:        false,
		Currency:           configData.Currency,
		Des:                args.Des,
		Logs:               []ServiceOrderWaitFields.FieldsLog{},
		ReferrerNationCode: args.ReferrerNationCode,
		ReferrerPhone:      args.ReferrerPhone,
		Params:             args.Params,
	})
	if err != nil {
		err = errors.New(fmt.Sprint("create_order_fail:", err))
		return
	}
	return
}
