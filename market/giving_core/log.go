package MarketGivingCore

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	FinanceDeposit "github.com/fotomxq/weeekj_core/v5/finance/deposit"
	MarketCore "github.com/fotomxq/weeekj_core/v5/market/core"
	OrgCore "github.com/fotomxq/weeekj_core/v5/org/core"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
	UserIntegral "github.com/fotomxq/weeekj_core/v5/user/integral"
	UserSubscriptionMod "github.com/fotomxq/weeekj_core/v5/user/subscription/mod"
	UserTicket "github.com/fotomxq/weeekj_core/v5/user/ticket"
	"github.com/golang-module/carbon"
	"math"
)

// ArgsGetLogList 获取日志列表参数
type ArgsGetLogList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//奖励的来源
	FromInfo CoreSQLFrom.FieldsFrom `db:"from_info" json:"fromInfo"`
	//对表成员的用户ID
	// 和成员对等，可用于一次性推荐的记录处理
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//推荐人
	ReferrerUserID int64 `db:"referrer_user_id" json:"referrerUserID" check:"id" empty:"true"`
	//推荐成员ID
	ReferrerBindID int64 `db:"referrer_bind_id" json:"referrerBindID" check:"id" empty:"true"`
	//奖励依据配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id" empty:"true"`
	//是否被删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetLogList 获取日志列表
func GetLogList(args *ArgsGetLogList) (dataList []FieldsLog, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.FromInfo.System != "" {
		where, maps, err = args.FromInfo.GetListAnd("from_info", "from_info", where, maps)
		if err != nil {
			return
		}
	}
	if args.UserID > -1 {
		where = where + " AND user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.ReferrerUserID > -1 {
		where = where + " AND referrer_user_id = :referrer_user_id"
		maps["referrer_user_id"] = args.ReferrerUserID
	}
	if args.ReferrerBindID > -1 {
		where = where + " AND referrer_bind_id = :referrer_bind_id"
		maps["referrer_bind_id"] = args.ReferrerBindID
	}
	if args.ConfigID > -1 {
		where = where + " AND config_id = :config_id"
		maps["config_id"] = args.ConfigID
	}
	if args.Search != "" {
		where = where + " AND (des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"market_giving_core_log",
		"id",
		"SELECT id, create_at, delete_at, org_id, from_info, user_id, referrer_user_id, config_id, user_integral, user_subs, user_tickets, price_total, deposit_config_mark, price, count, des, params FROM market_giving_core_log WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "delete_at"},
	)
	return
}

// ArgsCreateLog 新的营销记录参数
type ArgsCreateLog struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//奖励的来源
	FromInfo CoreSQLFrom.FieldsFrom `db:"from_info" json:"fromInfo"`
	//奖励的用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//推荐人用户ID
	ReferrerUserID int64 `db:"referrer_user_id" json:"referrerUserID" check:"id" empty:"true"`
	//推荐成员ID
	ReferrerBindID int64 `db:"referrer_bind_id" json:"referrerBindID" check:"id" empty:"true"`
	//奖励机制配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id"`
	//交易的金额
	// 用户发生交易的总金额
	PriceTotal int64 `db:"price_total" json:"priceTotal" check:"price" empty:"true"`
	//客户备注
	Des string `db:"des" json:"des" check:"des" min:"1" max:"1000" empty:"true"`
}

// CreateLog 新的营销记录
func CreateLog(args *ArgsCreateLog) (data FieldsLog, errCode string, err error) {
	//获取配置数据包
	var configData FieldsConfig
	configData, err = GetConfigByID(&ArgsGetConfigByID{
		ID:    args.ConfigID,
		OrgID: args.OrgID,
	})
	if err != nil || configData.ID < 1 {
		errCode = "config_not_exist"
		err = errors.New(fmt.Sprint("get config data, ", err))
		return
	}
	//检查配置的限制
	var count int64
	switch configData.LimitTimeType {
	case 0:
	case 1:
		count, err = CoreSQL.GetAllCount(Router2SystemConfig.MainDB.DB, "market_giving_core_log", "id", "org_id = $1 AND user_id = $2 AND delete_at < to_timestamp(1000000)", args.OrgID, args.UserID)
		if err == nil && count > int64(configData.LimitCount) {
			errCode = "config_limit"
			err = errors.New("config limit")
			return
		}
	case 2:
		count, err = CoreSQL.GetAllCount(Router2SystemConfig.MainDB.DB, "market_giving_core_log", "id", "org_id = $1 AND user_id = $2 AND delete_at < to_timestamp(1000000) AND create_at >= $3", args.OrgID, args.UserID, CoreFilter.GetNowTimeCarbon().SubDay().Time)
		if err == nil && count > int64(configData.LimitCount) {
			errCode = "config_limit"
			err = errors.New("config limit")
			return
		}
	case 3:
		count, err = CoreSQL.GetAllCount(Router2SystemConfig.MainDB.DB, "market_giving_core_log", "id", "org_id = $1 AND user_id = $2 AND delete_at < to_timestamp(1000000) AND create_at >= $3", args.OrgID, args.UserID, CoreFilter.GetNowTimeCarbon().SubWeek().Time)
		if err == nil && count > int64(configData.LimitCount) {
			errCode = "config_limit"
			err = errors.New("config limit")
			return
		}
	case 4:
		count, err = CoreSQL.GetAllCount(Router2SystemConfig.MainDB.DB, "market_giving_core_log", "id", "org_id = $1 AND user_id = $2 AND delete_at < to_timestamp(1000000) AND create_at >= $3", args.OrgID, args.UserID, CoreFilter.GetNowTimeCarbon().SubMonth().Time)
		if err == nil && count > int64(configData.LimitCount) {
			errCode = "config_limit"
			err = errors.New("config limit")
			return
		}
	case 5:
		count, err = CoreSQL.GetAllCount(Router2SystemConfig.MainDB.DB, "market_giving_core_log", "id", "org_id = $1 AND user_id = $2 AND delete_at < to_timestamp(1000000) AND create_at >= $3", args.OrgID, args.UserID, CoreFilter.GetNowTimeCarbon().SubMonths(3).Time)
		if err == nil && count > int64(configData.LimitCount) {
			errCode = "config_limit"
			err = errors.New("config limit")
			return
		}
	case 6:
		count, err = CoreSQL.GetAllCount(Router2SystemConfig.MainDB.DB, "market_giving_core_log", "id", "org_id = $1 AND user_id = $2 AND delete_at < to_timestamp(1000000) AND create_at >= $3", args.OrgID, args.UserID, CoreFilter.GetNowTimeCarbon().SubYear().Time)
		if err == nil && count > int64(configData.LimitCount) {
			errCode = "config_limit"
			err = errors.New("config limit")
			return
		}
	}
	//找到该绑定人的用户
	if args.ReferrerBindID > 0 && args.ReferrerUserID < 1 {
		var bindData OrgCore.FieldsBind
		bindData, err = OrgCore.GetBind(&OrgCore.ArgsGetBind{
			ID:     args.ReferrerBindID,
			OrgID:  args.OrgID,
			UserID: -1,
		})
		if err != nil {
			err = errors.New(fmt.Sprint("bind not exist, ", err))
			errCode = "bind_not_exist"
			return
		}
		var userData UserCore.FieldsUserType
		userData, err = UserCore.GetUserByID(&UserCore.ArgsGetUserByID{
			ID:    bindData.UserID,
			OrgID: -1,
		})
		if err != nil {
			errCode = "user_not_exist"
			err = errors.New(fmt.Sprint("get user data, ", err))
			return
		}
		args.ReferrerUserID = userData.ID
	}
	//根据配置数据，赠送给该用户相关内容
	haveGiving := false
	if configData.UserIntegral > 0 {
		err = UserIntegral.AddCount(&UserIntegral.ArgsAddCount{
			OrgID:    args.OrgID,
			UserID:   args.UserID,
			AddCount: configData.UserIntegral,
			Des:      "推荐成功奖励积分",
		})
		if err != nil {
			errCode = "user_integral"
			return
		}
		args.Des = fmt.Sprint(args.Des, "，赠送", configData.UserIntegral, "个积分")
		haveGiving = true
	}
	if len(configData.UserSubs) > 0 {
		for _, v := range configData.UserSubs {
			if v.ConfigID < 1 || (v.Count < 1 && v.CountTime < 1) {
				continue
			}
			var newExpireAt carbon.Carbon
			var subConfig UserSubscriptionMod.FieldsConfig
			var appendDes string
			if v.Count > 0 {
				subConfig, err = UserSubscriptionMod.GetConfigByID(&UserSubscriptionMod.ArgsGetConfigByID{
					ID:    v.ConfigID,
					OrgID: configData.OrgID,
				})
				if err != nil {
					errCode = "user_sub_config_not_exist"
					return
				}
				switch subConfig.TimeType {
				case 0:
					newExpireAt = CoreFilter.GetNowTimeCarbon().AddHours(subConfig.TimeN)
					appendDes = fmt.Sprint(subConfig.TimeN, "小时")
				case 1:
					newExpireAt = CoreFilter.GetNowTimeCarbon().AddDays(subConfig.TimeN)
					appendDes = fmt.Sprint(subConfig.TimeN, "天")
				case 2:
					newExpireAt = CoreFilter.GetNowTimeCarbon().AddWeeks(subConfig.TimeN)
					appendDes = fmt.Sprint(subConfig.TimeN, "周")
				case 3:
					newExpireAt = CoreFilter.GetNowTimeCarbon().AddMonths(subConfig.TimeN)
					appendDes = fmt.Sprint(subConfig.TimeN, "月")
				case 4:
					newExpireAt = CoreFilter.GetNowTimeCarbon().AddYears(subConfig.TimeN)
					appendDes = fmt.Sprint(subConfig.TimeN, "年")
				default:
					newExpireAt = CoreFilter.GetNowTimeCarbon()
				}
			} else {
				newExpireAt = CoreFilter.GetNowTimeCarbon().AddSeconds(int(v.CountTime))
			}
			err = UserSubscriptionMod.SetSub(UserSubscriptionMod.ArgsSetSub{
				OrgID:       configData.OrgID,
				ConfigID:    v.ConfigID,
				UserID:      args.UserID,
				ExpireAt:    newExpireAt.Time,
				HaveExpire:  true,
				UseFrom:     "market",
				UseFromName: fmt.Sprint("推荐人赠礼(", configData.Name, ")"),
			})
			if err != nil {
				errCode = "set_user_sub"
				return
			}
			args.Des = fmt.Sprint(args.Des, "，赠送", subConfig.Title, "订阅", appendDes)
			haveGiving = true
		}
	}
	if len(configData.UserTickets) > 0 {
		for _, v := range configData.UserTickets {
			if v.ConfigID < 1 || v.Count < 1 {
				continue
			}
			err = UserTicket.AddTicket(&UserTicket.ArgsAddTicket{
				OrgID:       configData.OrgID,
				ConfigID:    v.ConfigID,
				UserID:      args.UserID,
				Count:       v.Count,
				UseFromName: "推荐人赠礼",
			})
			if err != nil {
				errCode = "add_user_ticket"
				return
			}
			args.Des = fmt.Sprint(args.Des, "，赠送", configData.Name, "票据", v.Count, "张")
			haveGiving = true
		}
	}
	if configData.Price > 0 && configData.DepositConfigMark != "" {
		_, _, err = FinanceDeposit.SetByFrom(&FinanceDeposit.ArgsSetByFrom{
			UpdateHash: "",
			CreateInfo: CoreSQLFrom.FieldsFrom{
				System: "user",
				ID:     args.UserID,
				Mark:   "",
				Name:   "",
			},
			FromInfo: CoreSQLFrom.FieldsFrom{
				System: "org",
				ID:     configData.OrgID,
				Mark:   "",
				Name:   "",
			},
			ConfigMark:      configData.DepositConfigMark,
			AppendSavePrice: configData.Price,
		})
		if err != nil {
			errCode = "set_finance_deposit"
			return
		}
		args.Des = fmt.Sprint(args.Des, "，赠送储蓄", math.Round(float64(configData.Price/100)), "元")
		haveGiving = true
	}
	//如果存在赠礼内容，则记录
	if haveGiving {
		//创建新的记录
		err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "market_giving_core_log", "INSERT INTO market_giving_core_log (org_id, from_info, user_id, referrer_user_id, referrer_bind_id, config_id, user_integral, user_subs, user_tickets, price_total, deposit_config_mark, price, count, des, params) VALUES (:org_id,:from_info,:user_id,:referrer_user_id,:referrer_bind_id,:config_id,:user_integral,:user_subs,:user_tickets,:price_total,:deposit_config_mark,:price,:count,:des,:params)", map[string]interface{}{
			"org_id":              args.OrgID,
			"from_info":           args.FromInfo,
			"user_id":             args.UserID,
			"referrer_user_id":    args.ReferrerUserID,
			"referrer_bind_id":    args.ReferrerBindID,
			"config_id":           configData.ID,
			"user_integral":       configData.UserIntegral,
			"user_subs":           configData.UserSubs,
			"user_tickets":        configData.UserTickets,
			"price_total":         args.PriceTotal,
			"deposit_config_mark": configData.DepositConfigMark,
			"price":               configData.Price,
			"count":               configData.Count,
			"des":                 args.Des,
			"params":              CoreSQLConfig.FieldsConfigsType{},
		}, &data)
		if err != nil {
			errCode = "insert"
			return
		}
	}
	//推荐人奖励机制
	if configData.MarketConfigID > 0 {
		if (args.ReferrerUserID > 0 || args.ReferrerBindID > 0) && args.ReferrerUserID != args.UserID {
			//注意，奖励人和推荐人是反过来写的
			_, errCode, err = MarketCore.CreateLog(&MarketCore.ArgsCreateLog{
				OrgID:      args.OrgID,
				UserID:     args.ReferrerUserID,
				BindID:     args.ReferrerBindID,
				BindUserID: args.UserID,
				BindInfoID: 0,
				ConfigID:   configData.MarketConfigID,
				FromInfo:   args.FromInfo,
				PriceTotal: args.PriceTotal,
				Des:        args.Des,
			})
			if err != nil {
				return
			}
		}
	}
	return
}

// ArgsDeleteLog 销毁营销记录参数
type ArgsDeleteLog struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 可选
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteLog 销毁营销记录
func DeleteLog(args *ArgsDeleteLog) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "market_giving_core_log", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}
