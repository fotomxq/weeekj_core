package MarketGivingBuyMall

import (
	"fmt"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	MarketGivingCore "gitee.com/weeekj/weeekj_core/v5/market/giving_core"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsCreateLog 创建新的请求参数
type ArgsCreateLog struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//奖励的用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//推荐人用户ID
	ReferrerUserID int64 `db:"referrer_user_id" json:"referrerUserID" check:"id" empty:"true"`
	//推荐成员ID
	ReferrerBindID int64 `db:"referrer_bind_id" json:"referrerBindID" check:"id" empty:"true"`
	//交易的金额
	// 用户发生交易的总金额
	PriceTotal int64 `db:"price_total" json:"priceTotal" check:"price" empty:"true"`
	//商品ID
	MallProductID pq.Int64Array `db:"mall_product_id" json:"mallProductID" check:"ids" empty:"true"`
	//商品分类
	SortID pq.Int64Array `db:"sort_id" json:"sortID" check:"ids" empty:"true"`
	//商品标签
	Tag pq.Int64Array `db:"tag" json:"tag" check:"ids" empty:"true"`
}

// CreateLog 创建新的请求
func CreateLog(args *ArgsCreateLog) (errCode string, err error) {
	//检查符合的条件
	var conditionsList []FieldsConditions
	err = Router2SystemConfig.MainDB.Select(&conditionsList, "SELECT id, config_id, params FROM market_giving_buy_mall WHERE delete_at < to_timestamp(1000000) AND org_id = $1 AND min_price <= $2 AND (mall_product_id = ANY($3) OR sort_id = ANY($4) OR tag = ANY($5))", args.OrgID, args.PriceTotal, args.MallProductID, args.SortID, args.Tag)
	if err != nil || len(conditionsList) < 1 {
		//没有符合条件则退出
		err = nil
		return
	}
	//遍历条件并触发
	for _, v := range conditionsList {
		_, errCode, err = MarketGivingCore.CreateLog(&MarketGivingCore.ArgsCreateLog{
			OrgID: args.OrgID,
			FromInfo: CoreSQLFrom.FieldsFrom{
				System: "buy_mall",
				ID:     0,
				Mark:   "",
				Name:   "",
			},
			UserID:         args.UserID,
			ReferrerUserID: args.ReferrerUserID,
			ReferrerBindID: args.ReferrerBindID,
			ConfigID:       v.ConfigID,
			PriceTotal:     args.PriceTotal,
			Des:            fmt.Sprint("用户购买商品奖励"),
		})
		if err != nil {
			if errCode == "config_limit" {
				err = nil
				continue
			}
			return
		}
	}
	//反馈
	return
}
