package FinanceAnalysis

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	CoreSQLTime "gitee.com/weeekj/weeekj_core/v5/core/sql/time"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsGetAnalysisTakePrice 储蓄账户金额参数
type ArgsGetAnalysisTakePrice struct {
	//收款渠道
	// system: cash 现金 ; deposit 存储模块 ; weixin 微信支付 ; alipay 支付宝
	TakeChannel CoreSQLFrom.FieldsFrom `db:"take_channel" json:"takeChannel"`
	//收款方来源
	// 留空则代表平台方，否则为商户或加盟商
	TakeFrom CoreSQLFrom.FieldsFrom `db:"take_from" json:"takeFrom"`
	//时间范围
	// 部分统计支持
	TimeBetween CoreSQLTime.DataCoreTime `json:"timeBetween"`
}

// GetAnalysisTakePrice 检查目标转入的资金总量
func GetAnalysisTakePrice(args *ArgsGetAnalysisTakePrice) (count int64, err error) {
	var timeBetween CoreSQLTime.FieldsCoreTime
	timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
	if err != nil {
		return
	}
	where := "day_time >= :start_at AND day_time <= :end_at"
	maps := map[string]interface{}{
		"start_at": timeBetween.MinTime,
		"end_at":   timeBetween.MaxTime,
	}
	var newWhere string
	newWhere, maps, err = args.TakeChannel.GetList("take_channel", "take_channel", maps)
	if err != nil {
		return
	} else {
		if newWhere != "" {
			where = where + " AND " + newWhere
		}
	}
	newWhere, maps, err = args.TakeChannel.GetList("take_from", "take_from", maps)
	if err != nil {
		return
	} else {
		if newWhere != "" {
			where = where + " AND " + newWhere
		}
	}
	count, err = CoreSQL.GetAllSumMap(Router2SystemConfig.MainDB.DB, "finance_analysis", "price", where, maps)
	return
}
