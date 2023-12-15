package ServiceOrderAnalysis

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsCreateAnalysis 新增统计数据
type ArgsCreateAnalysis struct {
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//来源系统
	// 该订单创建来源的系统
	// eg: user_sub / org_sub / mall
	SystemMark string `db:"system_mark" json:"systemMark"`
	//创建来源
	// 0 商户平台下单、1 用户APP下单、2 微信小程序下单、3 支付宝小程序下单、4 电话系统下单、5 线下柜台下单、6 辅助设备下单、7 其他渠道
	CreateFrom int `db:"create_from" json:"createFrom"`
	// 货币
	Currency int `db:"currency" json:"currency"`
	// 总费用金额
	Price int64 `db:"price" json:"price"`
	//商品列
	Goods []ArgsCreateAnalysisGood `json:"goods"`
	//抵扣
	Exemptions []ArgsCreateAnalysisExemption `json:"exemptions"`
}

type ArgsCreateAnalysisGood struct {
	//获取来源
	// 如果商品mark带有virtual标记，且订单商品全部带有该标记，订单将在付款后直接完成
	From CoreSQLFrom.FieldsFrom `db:"from" json:"from"`
	//货物个数
	Count int64 `db:"count" json:"count"`
}

type ArgsCreateAnalysisExemption struct {
	//抵扣系统来源
	// integral 积分; ticket 票据; sub 订阅
	System string `db:"system" json:"system"`
	//抵扣配置ID
	// 可能不存在，如积分没有配置ID
	ConfigID int64 `db:"config_id" json:"configID"`
	//使用数量
	// 使用的张数、或使用积分的个数
	Count int64 `db:"count" json:"count"`
	//抵扣费用
	Price int64 `db:"price" json:"price"`
}

func CreateAnalysis(args *ArgsCreateAnalysis) (err error) {
	//构建小时数据
	minTime := CoreFilter.GetNowTimeCarbon().StartOfHour().Time
	maxTime := CoreFilter.GetNowTimeCarbon().EndOfHour().Time
	//组织基础统计
	var dataOrgOrder FieldsOrg
	if err = Router2SystemConfig.MainDB.Get(&dataOrgOrder, "SELECT id FROM service_order_analysis_org WHERE org_id = $1 AND day_time >= $2 AND day_time <= $3 AND create_from = $4 AND currency = $5 LIMIT 1", args.OrgID, minTime, maxTime, args.CreateFrom, args.Currency); err != nil {
		err = nil
		_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO service_order_analysis_org (org_id, create_from, count, currency, price) VALUES (:org_id,:create_from,1,:currency,:price)", map[string]interface{}{
			"org_id":      args.OrgID,
			"create_from": args.CreateFrom,
			"currency":    args.Currency,
			"price":       args.Price,
		})
	} else {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE service_order_analysis_org SET count = count + :count, price = price + :price WHERE id = :id", map[string]interface{}{
			"id":    dataOrgOrder.ID,
			"count": 1,
			"price": args.Price,
		})
	}
	//用户基础统计
	var dataOrgUser FieldsUser
	if err = Router2SystemConfig.MainDB.Get(&dataOrgUser, "SELECT id FROM service_order_analysis_user WHERE user_id = $1 AND day_time >= $2 AND day_time <= $3 AND currency = $4 LIMIT 1", args.UserID, minTime, maxTime, args.Currency); err != nil {
		err = nil
		_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO service_order_analysis_user (user_id, count, currency, price) VALUES (:user_id,1,:currency,:price)", map[string]interface{}{
			"org_id":   args.OrgID,
			"user_id":  args.UserID,
			"currency": args.Currency,
			"price":    args.Price,
		})
	} else {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE service_order_analysis_user SET count = count + :count, price = price + :price WHERE id = :id", map[string]interface{}{
			"id":    dataOrgUser.ID,
			"count": 1,
			"price": args.Price,
		})
	}
	//商品数量统计
	for _, v := range args.Goods {
		var dataFromCount FieldsFromCount
		if err = Router2SystemConfig.MainDB.Get(&dataFromCount, "SELECT id FROM service_order_analysis_from_count WHERE org_id = $1 AND day_time >= $2 AND day_time <= $3 AND from_system = $4 AND from_id = $5 LIMIT 1", &dataFromCount, v.From.System, v.From.ID); err != nil {
			_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO service_order_analysis_from_count (org_id, from_system, from_id, buy_count) VALUES (:org_id,:from_system,:from_id,:buy_count)", map[string]interface{}{
				"org_id":      args.OrgID,
				"user_id":     args.UserID,
				"from_system": v.From.System,
				"from_id":     v.From.ID,
				"buy_count":   v.Count,
			})
		} else {
			_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE service_order_analysis_from_count SET buy_count = buy_count + :buy_count WHERE id = :id", map[string]interface{}{
				"id":        dataOrgUser.ID,
				"buy_count": v.Count,
			})
		}
	}
	//优惠统计
	for _, v := range args.Exemptions {
		var dataOrgExemption FieldsOrgExemption
		if err = Router2SystemConfig.MainDB.Get(&dataOrgExemption, "SELECT id FROM service_order_analysis_org_exemption WHERE org_id = $1 AND day_time >= $2 AND day_time <= $3 AND from_system = $4 AND config_id = $5 LIMIT 1", &dataOrgExemption, v.System, v.ConfigID); err != nil {
			_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO service_order_analysis_org_exemption (org_id, from_system, config_id, count, price) VALUES (:org_id,:from_system,:config_id,:count,:price)", map[string]interface{}{
				"org_id":      args.OrgID,
				"user_id":     args.UserID,
				"from_system": v.System,
				"config_id":   v.ConfigID,
				"count":       v.Count,
				"price":       v.Price,
			})
		} else {
			_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE service_order_analysis_org_exemption SET count = count + :buy_count, price = price + :price WHERE id = :id", map[string]interface{}{
				"id":    dataOrgExemption.ID,
				"count": v.Count,
				"price": v.Price,
			})
		}
	}
	//反馈
	return
}
