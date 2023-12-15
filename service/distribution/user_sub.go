package ServiceDistribution

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"

	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	BaseQRCode "github.com/fotomxq/weeekj_core/v5/base/qrcode"
	BaseWeixinWXXQRCodeCore "github.com/fotomxq/weeekj_core/v5/base/weixin/wxx/qrcode"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLAddress "github.com/fotomxq/weeekj_core/v5/core/sql/address"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	MarketGivingUserSub "github.com/fotomxq/weeekj_core/v5/market/giving_user_sub"
	ServiceOrderWaitFields "github.com/fotomxq/weeekj_core/v5/service/order/wait_fields"
	UserSubscription "github.com/fotomxq/weeekj_core/v5/user/subscription"
	UserSubscriptionMod "github.com/fotomxq/weeekj_core/v5/user/subscription/mod"
)

// ArgsGetUserSubList 获取关联的用户订阅列表参数
type ArgsGetUserSubList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//对应分销商
	DistributionID int64 `db:"distribution_id" json:"distributionID" check:"id" empty:"true"`
	//是否被删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool" empty:"true"`
}

// GetUserSubList 获取关联的用户订阅列表
func GetUserSubList(args *ArgsGetUserSubList) (dataList []FieldsUserSub, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.DistributionID > -1 {
		where = where + " AND distribution_id = :distribution_id"
		maps["distribution_id"] = args.DistributionID
	}
	tableName := "service_distribution_user_sub"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, update_at, delete_at, org_id, distribution_id, sub_config_id, unit_price, market_giving_sub_id, cover_file_id, des, in_count, order_price, order_count FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	return
}

// ArgsInUserSub 扫码进入操作参数
type ArgsInUserSub struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// InUserSub 扫码进入操作
func InUserSub(args *ArgsInUserSub) (data FieldsUserSub, err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_distribution_user_sub SET in_count = in_count + 1 WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, distribution_id, sub_config_id, unit_price, market_giving_sub_id, cover_file_id, des, in_count, order_price, order_count FROM service_distribution_user_sub WHERE id = $1", args.ID)
	if err == nil && data.ID < 1 {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsBuyUserSub 购买操作参数
type ArgsBuyUserSub struct {
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
	//对应分销商
	DistributionUserSubID int64 `db:"distribution_user_sub_id" json:"distributionUserSubID" check:"id"`
}

// BuyUserSub 购买操作
func BuyUserSub(args *ArgsBuyUserSub) (data ServiceOrderWaitFields.FieldsWait, errCode string, err error) {
	//获取加盟商数据
	var distributionUserSubData FieldsUserSub
	err = Router2SystemConfig.MainDB.Get(&distributionUserSubData, "SELECT id, unit_price, market_giving_sub_id FROM service_distribution_user_sub WHERE id = $1 AND delete_at < to_timestamp(1000000)", args.DistributionUserSubID)
	if err != nil || distributionUserSubData.ID < 1 {
		errCode = "distribution_not_exist"
		err = errors.New("distribution not exist")
		return
	}
	args.Params = CoreSQLConfig.Set(args.Params, "force_market_giving_sub_id", distributionUserSubData.MarketGivingSubID)
	//创建订单
	data, errCode, err = UserSubscription.CreateSubOrder(&UserSubscription.ArgsCreateSubOrder{
		OrgID:              args.OrgID,
		UserID:             args.UserID,
		CreateFrom:         args.CreateFrom,
		AddressFrom:        args.AddressFrom,
		AddressTo:          args.AddressTo,
		Des:                args.Des,
		SubConfigID:        args.SubConfigID,
		Unit:               args.Unit,
		ReferrerNationCode: args.ReferrerNationCode,
		ReferrerPhone:      args.ReferrerPhone,
		Params:             args.Params,
		UnitPrice:          distributionUserSubData.UnitPrice,
	})
	return
}

// ArgsCreateUserSub 创建用户订阅关联参数
type ArgsCreateUserSub struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//对应分销商
	DistributionID int64 `db:"distribution_id" json:"distributionID" check:"id"`
	//会员配置
	SubConfigID int64 `db:"sub_config_id" json:"subConfigID" check:"id"`
	//强制约定价格
	UnitPrice int64 `db:"unit_price" json:"unitPrice" check:"price"`
	//指定奖励
	MarketGivingSubID int64 `db:"market_giving_sub_id" json:"marketGivingSubID" check:"id" empty:"true"`
	//宣传海报
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"1000" empty:"true"`
}

// CreateUserSub 创建用户订阅关联
func CreateUserSub(args *ArgsCreateUserSub) (err error) {
	//检查订阅
	err = UserSubscriptionMod.CheckConfigAndOrg(&UserSubscriptionMod.ArgsCheckConfigAndOrg{
		ID:    args.SubConfigID,
		OrgID: args.OrgID,
	})
	if err != nil {
		err = errors.New("user sub config not exist")
		return
	}
	//检查赠送
	if args.MarketGivingSubID > 0 {
		err = MarketGivingUserSub.CheckConfigAndOrg(&MarketGivingUserSub.ArgsCheckConfigAndOrg{
			ID:    args.MarketGivingSubID,
			OrgID: args.OrgID,
		})
		if err != nil {
			err = errors.New("market giving user sub config not exist")
			return
		}
	}
	//创建数据
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO service_distribution_user_sub (org_id, distribution_id, sub_config_id, unit_price, market_giving_sub_id, cover_file_id, des, in_count, order_price, order_count) VALUES (:org_id,:distribution_id,:sub_config_id,:unit_price,:market_giving_sub_id,:cover_file_id,:des,0,0,0)", args)
	return
}

// ArgsUpdateUserSub 修改用户订阅关联参数
type ArgsUpdateUserSub struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//对应分销商
	DistributionID int64 `db:"distribution_id" json:"distributionID" check:"id"`
	//会员配置
	SubConfigID int64 `db:"sub_config_id" json:"subConfigID" check:"id"`
	//强制约定价格
	UnitPrice int64 `db:"unit_price" json:"unitPrice" check:"price"`
	//指定奖励
	MarketGivingSubID int64 `db:"market_giving_sub_id" json:"marketGivingSubID" check:"id" empty:"true"`
	//宣传海报
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"1000" empty:"true"`
}

// UpdateUserSub 修改用户订阅关联
func UpdateUserSub(args *ArgsUpdateUserSub) (err error) {
	//检查订阅
	err = UserSubscriptionMod.CheckConfigAndOrg(&UserSubscriptionMod.ArgsCheckConfigAndOrg{
		ID:    args.SubConfigID,
		OrgID: args.OrgID,
	})
	if err != nil {
		err = errors.New("user sub config not exist")
		return
	}
	//检查赠送
	if args.MarketGivingSubID > 0 {
		err = MarketGivingUserSub.CheckConfigAndOrg(&MarketGivingUserSub.ArgsCheckConfigAndOrg{
			ID:    args.MarketGivingSubID,
			OrgID: args.OrgID,
		})
		if err != nil {
			err = errors.New("market giving user sub config not exist")
			return
		}
	}
	//更新数据
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_distribution_user_sub SET update_at = NOW(), distribution_id = :distribution_id, sub_config_id = :sub_config_id, unit_price = :unit_price, market_giving_sub_id = :market_giving_sub_id, cover_file_id = :cover_file_id, des = :des WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}

// ArgsDeleteUserSub 删除用户订阅关联参数
type ArgsDeleteUserSub struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteUserSub 删除用户订阅关联
func DeleteUserSub(args *ArgsDeleteUserSub) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "service_distribution_user_sub", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}

// ArgsGetUserSubQrcode 获取二维码参数
type ArgsGetUserSubQrcode struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//二维码类型
	QrcodeType string `json:"qrcodeType"`
	//尺寸
	// eg: 430
	Size int `db:"size" json:"size" check:"intThan0"`
	//是否需要透明底色
	IsHyaline bool `json:"isHyaline"`
	//自动配置线条颜色
	// 为 false 时生效, 使用 rgb 设置颜色 十进制表示
	AutoColor bool `json:"autoColor"`
	//色调
	// 50
	R string `json:"r"`
	G string `json:"g"`
	B string `json:"b"`
}

// GetUserSubQrcode 获取二维码
func GetUserSubQrcode(args *ArgsGetUserSubQrcode) (qrcodeData string, err error) {
	//尺寸不能太小
	if args.Size < 10 {
		err = errors.New("size too small")
		return
	}
	if args.Size > 1024 {
		err = errors.New("size too big")
		return
	}
	//获取会员数据包
	var data FieldsUserSub
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, org_id FROM service_distribution_user_sub WHERE id = $1 AND ($2 < 1 OR org_id = $2)", args.ID, args.OrgID)
	if err != nil || data.ID < 1 {
		err = errors.New("no data")
		return
	}
	//根据类型生成二维码
	switch args.QrcodeType {
	case "weixin_wxx":
		var dataByte []byte
		dataByte, err = BaseWeixinWXXQRCodeCore.GetQRByParam(&BaseWeixinWXXQRCodeCore.ArgsGetQRByParam{
			MerchantID: data.OrgID,
			Page:       "/pages/distribution/user_sub",
			Param:      fmt.Sprint("id=", args.ID),
			Width:      args.Size,
			IsHyaline:  args.IsHyaline,
			AutoColor:  args.AutoColor,
			R:          args.R,
			G:          args.G,
			B:          args.B,
		})
		if err != nil {
			return
		}
		qrcodeData = base64.StdEncoding.EncodeToString(dataByte)
		return
	case "app":
		type dataType struct {
			Action string `json:"action"`
			ID     int64  `json:"id"`
		}
		d := dataType{
			Action: "distribution_user_sub",
			ID:     data.ID,
		}
		var dJson []byte
		dJson, err = json.Marshal(d)
		if err != nil {
			return
		}
		return BaseQRCode.GetQRCode(&BaseQRCode.ArgsGetQRCode{
			Param: string(dJson),
			Size:  args.Size,
		})
	default:
		var appURL string
		appURL, err = BaseConfig.GetDataString("AppURL")
		if err != nil {
			return
		}
		return BaseQRCode.GetQRCode(&BaseQRCode.ArgsGetQRCode{
			Param: fmt.Sprint(appURL, "/distribution/user_sub?id=", data.ID),
			Size:  args.Size,
		})
	}
}
