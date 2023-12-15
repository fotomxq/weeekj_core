package MarketCore

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	FinanceDeposit "github.com/fotomxq/weeekj_core/v5/finance/deposit"
	OrgCore "github.com/fotomxq/weeekj_core/v5/org/core"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
	UserIntegral "github.com/fotomxq/weeekj_core/v5/user/integral"
	UserSubscriptionMod "github.com/fotomxq/weeekj_core/v5/user/subscription/mod"
	UserTicket "github.com/fotomxq/weeekj_core/v5/user/ticket"
	"github.com/golang-module/carbon"
	"github.com/lib/pq"
	"math"
)

// ArgsGetLogList 获取日志列表参数
type ArgsGetLogList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//对表成员的用户ID
	// 和成员对等，可用于一次性推荐的记录处理
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//成员ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//绑定的用户
	BindUserID int64 `db:"bind_user_id" json:"bindUserID" check:"id" empty:"true"`
	//绑定的档案
	BindInfoID int64 `db:"bind_info_id" json:"bindInfoID" check:"id" empty:"true"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签列
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//奖励的来源
	FromInfo CoreSQLFrom.FieldsFrom `db:"from_info" json:"fromInfo"`
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
	if args.UserID > -1 {
		where = where + " AND user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.BindID > -1 {
		where = where + " AND bind_id = :bind_id"
		maps["bind_id"] = args.BindID
	}
	if args.BindUserID > -1 {
		where = where + " AND bind_user_id = :bind_user_id"
		maps["bind_user_id"] = args.BindUserID
	}
	if args.BindInfoID > -1 {
		where = where + " AND bind_info_id = :bind_info_id"
		maps["bind_info_id"] = args.BindInfoID
	}
	if args.SortID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "sort_id = :sort_id"
		maps["sort_id"] = args.SortID
	}
	if len(args.Tags) > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "tags @> :tags"
		maps["tags"] = args.Tags
	}
	if args.FromInfo.System != "" {
		where, maps, err = args.FromInfo.GetListAnd("from_info", "from_info", where, maps)
		if err != nil {
			return
		}
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
		"market_core_log",
		"id",
		"SELECT id, create_at, delete_at, org_id, user_id, bind_id, bind_user_id, bind_info_id, sort_id, tags, user_integral, user_subs, user_tickets, from_info, price_total, deposit_config_mark, price, count, config_id, des, params FROM market_core_log WHERE "+where,
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
	//对表成员的用户ID
	// 和成员对等，可用于一次性推荐的记录处理
	// 可以不给予，则按照成员ID走
	// 推荐的人用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//成员ID
	// 可以不提供，会自动根据该用户关联营销人员走
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//绑定的用户
	// 被推荐的用户ID
	BindUserID int64 `db:"bind_user_id" json:"bindUserID" check:"id" empty:"true"`
	//绑定的档案
	BindInfoID int64 `db:"bind_info_id" json:"bindInfoID" check:"id" empty:"true"`
	//奖励机制配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id"`
	//奖励的来源
	FromInfo CoreSQLFrom.FieldsFrom `db:"from_info" json:"fromInfo"`
	//交易的金额
	// 用户发生交易的总金额
	PriceTotal int64 `db:"price_total" json:"priceTotal" check:"price" empty:"true"`
	//客户备注
	Des string `db:"des" json:"des" check:"des" min:"1" max:"1000" empty:"true"`
}

// CreateLog 新的营销记录
func CreateLog(args *ArgsCreateLog) (data FieldsLog, errCode string, err error) {
	//必须存在档案或用户
	if args.BindUserID < 1 && args.BindInfoID < 1 {
		errCode = "user_info_empty"
		err = errors.New("user id or info id less 1")
		return
	}
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
		count, err = CoreSQL.GetAllCount(Router2SystemConfig.MainDB.DB, "market_core_log", "id", "org_id = $1 AND user_id = $2 AND bind_id = $3 AND bind_user_id = $4 AND delete_at < to_timestamp(1000000)", args.OrgID, args.UserID, args.BindID, args.BindUserID)
		if err == nil && count > int64(configData.LimitCount) {
			errCode = "config_limit"
			err = errors.New("config limit")
			return
		}
	case 2:
		count, err = CoreSQL.GetAllCount(Router2SystemConfig.MainDB.DB, "market_core_log", "id", "org_id = $1 AND user_id = $2 AND bind_id = $3 AND bind_user_id = $4 AND delete_at < to_timestamp(1000000) AND create_at >= $5", args.OrgID, args.UserID, args.BindID, args.BindUserID, CoreFilter.GetNowTimeCarbon().SubDay().Time)
		if err == nil && count > int64(configData.LimitCount) {
			errCode = "config_limit"
			err = errors.New("config limit")
			return
		}
	case 3:
		count, err = CoreSQL.GetAllCount(Router2SystemConfig.MainDB.DB, "market_core_log", "id", "org_id = $1 AND user_id = $2 AND bind_id = $3 AND bind_user_id = $4 AND delete_at < to_timestamp(1000000) AND create_at >= $5", args.OrgID, args.UserID, args.BindID, args.BindUserID, CoreFilter.GetNowTimeCarbon().SubWeek().Time)
		if err == nil && count > int64(configData.LimitCount) {
			errCode = "config_limit"
			err = errors.New("config limit")
			return
		}
	case 4:
		count, err = CoreSQL.GetAllCount(Router2SystemConfig.MainDB.DB, "market_core_log", "id", "org_id = $1 AND user_id = $2 AND bind_id = $3 AND bind_user_id = $4 AND delete_at < to_timestamp(1000000) AND create_at >= $5", args.OrgID, args.UserID, args.BindID, args.BindUserID, CoreFilter.GetNowTimeCarbon().SubMonth().Time)
		if err == nil && count > int64(configData.LimitCount) {
			errCode = "config_limit"
			err = errors.New("config limit")
			return
		}
	case 5:
		count, err = CoreSQL.GetAllCount(Router2SystemConfig.MainDB.DB, "market_core_log", "id", "org_id = $1 AND user_id = $2 AND bind_id = $3 AND bind_user_id = $4 AND delete_at < to_timestamp(1000000) AND create_at >= $5", args.OrgID, args.UserID, args.BindID, args.BindUserID, CoreFilter.GetNowTimeCarbon().SubMonths(3).Time)
		if err == nil && count > int64(configData.LimitCount) {
			errCode = "config_limit"
			err = errors.New("config limit")
			return
		}
	case 6:
		count, err = CoreSQL.GetAllCount(Router2SystemConfig.MainDB.DB, "market_core_log", "id", "org_id = $1 AND user_id = $2 AND bind_id = $3 AND bind_user_id = $4 AND delete_at < to_timestamp(1000000) AND create_at >= $5", args.OrgID, args.UserID, args.BindID, args.BindUserID, CoreFilter.GetNowTimeCarbon().SubYear().Time)
		if err == nil && count > int64(configData.LimitCount) {
			errCode = "config_limit"
			err = errors.New("config limit")
			return
		}
	}
	//检查和获取用户关联的销售人员
	var sortID int64 = 0
	var tags = pq.Int64Array{}
	if args.UserID < 1 && args.BindID < 1 {
		var bindData FieldsBind
		err = Router2SystemConfig.MainDB.Get(&bindData, "SELECT id, sort_id, tags, bind_id FROM market_core_bind WHERE (bind_user_id = $1 OR bind_info_id = $2) AND org_id = $3 AND delete_at < to_timestamp(1000000)", args.BindUserID, args.BindInfoID, args.OrgID)
		if err != nil || bindData.ID < 1 {
			err = errors.New(fmt.Sprint("get bind data, ", err))
			errCode = "market_bind_not_exist"
			return
		}
		args.BindID = bindData.BindID
		sortID = bindData.SortID
		tags = bindData.Tags
	}
	//创建新的记录
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "market_core_log", "INSERT INTO market_core_log (org_id, user_id, bind_id, bind_user_id, bind_info_id, sort_id, tags, user_integral, user_subs, user_tickets, from_info, price_total, deposit_config_mark, price, count, config_id, des, params) VALUES (:org_id,:user_id,:bind_id,:bind_user_id,:bind_info_id,:sort_id,:tags,:user_integral,:user_subs,:user_tickets,:from_info,:price_total,:deposit_config_mark,:price,:count,:config_id,:des,:params)", map[string]interface{}{
		"org_id":              args.OrgID,
		"user_id":             args.UserID,
		"bind_id":             args.BindID,
		"bind_user_id":        args.BindUserID,
		"bind_info_id":        args.BindInfoID,
		"sort_id":             sortID,
		"tags":                tags,
		"user_integral":       configData.UserIntegral,
		"user_subs":           configData.UserSubs,
		"user_tickets":        configData.UserTickets,
		"from_info":           args.FromInfo,
		"price_total":         args.PriceTotal,
		"deposit_config_mark": configData.DepositConfigMark,
		"price":               configData.Price,
		"count":               configData.Count,
		"config_id":           configData.ID,
		"des":                 args.Des,
		"params":              CoreSQLConfig.FieldsConfigsType{},
	}, &data)
	if err != nil {
		errCode = "insert"
		return
	}
	//找到该绑定人的用户
	if args.UserID < 1 {
		var bindData OrgCore.FieldsBind
		bindData, err = OrgCore.GetBind(&OrgCore.ArgsGetBind{
			ID:     args.BindID,
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
		args.UserID = userData.ID
	}
	//根据配置数据，赠送给该用户相关内容
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
	}
	if len(configData.UserSubs) > 0 {
		for _, v := range configData.UserSubs {
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
				UseFromName: fmt.Sprint("营销赠礼(", configData.Name, ")"),
			})
			if err != nil {
				errCode = "set_user_sub"
				return
			}
			args.Des = fmt.Sprint(args.Des, "，赠送", subConfig.Title, "订阅", appendDes)
		}
	}
	if len(configData.UserTickets) > 0 {
		for _, v := range configData.UserTickets {
			err = UserTicket.AddTicket(&UserTicket.ArgsAddTicket{
				OrgID:       configData.OrgID,
				ConfigID:    v.ConfigID,
				UserID:      args.UserID,
				Count:       v.Count,
				UseFromName: "营销赠礼",
			})
			if err != nil {
				errCode = "add_user_ticket"
				return
			}
			args.Des = fmt.Sprint(args.Des, "，赠送", configData.Name, "票据", v.Count, "张")
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
	//成员ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
}

// DeleteLog 销毁营销记录
func DeleteLog(args *ArgsDeleteLog) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "market_core_log", "id = :id AND (:org_id < 1 OR org_id = :org_id) AND (:bind_id < 1 OR bind_id = :bind_id)", args)
	return
}
