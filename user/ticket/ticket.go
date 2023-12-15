package UserTicket

import (
	"database/sql"
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	OrgUserMod "github.com/fotomxq/weeekj_core/v5/org/user/mod"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"time"
)

// ArgsGetTicketList 获取票据列表参数
type ArgsGetTicketList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//是否到期
	NeedIsExpire bool `json:"needIsExpire" check:"bool" empty:"true"`
	IsExpire     bool `json:"isExpire" check:"bool" empty:"true"`
	//是否被删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool" empty:"true"`
	//是否聚合配置
	NeedAgg bool `json:"needAgg"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetTicketList 获取票据列表
func GetTicketList(args *ArgsGetTicketList) (dataList []FieldsTicket, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.ConfigID > -1 {
		where = where + " AND config_id = :config_id"
		maps["config_id"] = args.ConfigID
	}
	if args.UserID > -1 {
		where = where + " AND user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.NeedIsExpire {
		if args.IsExpire {
			where = where + " AND (expire_at > to_timestamp(1000000) AND expire_at < NOW())"
		} else {
			where = where + " AND (expire_at < to_timestamp(1000000) OR expire_at >= NOW())"
		}
	}
	tableName := "user_ticket"
	if args.NeedAgg {
		args.Pages.Sort = "config_id"
		dataCount, err = CoreSQL.GetListPageAndCount(
			Router2SystemConfig.MainDB.DB,
			&dataList,
			tableName,
			"config_id",
			"SELECT config_id, SUM(count) as count, SUM(res_count) as res_count FROM user_ticket WHERE "+where+" GROUP BY config_id, org_id, user_id",
			where,
			maps,
			&args.Pages,
			[]string{"config_id"},
		)
		return
	} else {
		dataCount, err = CoreSQL.GetListPageAndCount(
			Router2SystemConfig.MainDB.DB,
			&dataList,
			tableName,
			"id",
			"SELECT id, create_at, update_at, delete_at, expire_at, org_id, user_id, config_id, count, res_count FROM "+tableName+" WHERE "+where,
			where,
			maps,
			&args.Pages,
			[]string{"id", "create_at", "update_at", "delete_at", "expire_at"},
		)
		return
	}
}

// ArgsGetTicketCount 查询用户可用票数参数
type ArgsGetTicketCount struct {
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
}

// GetTicketCount 查询用户可用票数
func GetTicketCount(args *ArgsGetTicketCount) (count int64, err error) {
	type newData struct {
		Count int64 `db:"count"`
	}
	var data newData
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT SUM(count) as count FROM user_ticket WHERE config_id = $1 AND user_id = $2 AND (expire_at < to_timestamp(1000000) OR expire_at >= NOW()) AND delete_at < to_timestamp(1000000)", args.ConfigID, args.UserID)
	if err != nil {
		return
	}
	if err == nil && data.Count > 0 {
		count = data.Count
	}
	//更新组织用户数据
	OrgUserMod.PushUpdateUserData(0, args.UserID)
	//反馈
	return
}

// GetTicketCountByOrgID 查询组织持有的票数
func GetTicketCountByOrgID(orgID int64) (count int64) {
	_ = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM user_ticket WHERE org_id = $1 AND (expire_at < to_timestamp(1000000) OR expire_at >= NOW()) AND delete_at < to_timestamp(1000000)", orgID)
	return
}

// ArgsAddTicket 添加指定用户和票据配置的张数参数
type ArgsAddTicket struct {
	//组织ID
	// 可选
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//张数
	Count int64 `db:"count" json:"count" check:"int64Than0"`
	//使用来源
	UseFromName string `db:"use_from_name" json:"useFromName"`
}

// AddTicket 添加指定用户和票据配置的张数
func AddTicket(args *ArgsAddTicket) (err error) {
	//获取配置
	var configData FieldsConfig
	configData, err = GetConfigByID(&ArgsGetConfigByID{
		ID:    args.ConfigID,
		OrgID: args.OrgID,
	})
	if err != nil {
		return
	}
	if configData.DeleteAt.Unix() > 0 {
		err = errors.New("config not exist")
		return
	}
	//生成过期时间
	var expireAt time.Time
	if configData.DefaultExpireTime > 0 {
		expireAt = CoreFilter.GetNowTimeCarbon().AddSeconds(int(configData.DefaultExpireTime)).Time
	}
	//构建新的票据
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO user_ticket (expire_at, org_id, user_id, config_id, count, res_count) VALUES (:expire_at, :org_id, :user_id, :config_id, :count, :count)", map[string]interface{}{
		"expire_at": expireAt,
		"org_id":    configData.OrgID,
		"config_id": configData.ID,
		"user_id":   args.UserID,
		"count":     args.Count,
	})
	if err != nil {
		return
	}
	//记录日志
	_ = appendLog(configData.OrgID, configData.ID, args.UserID, 1, args.Count, fmt.Sprint(args.UseFromName, "添加[", args.Count, "]张[", configData.Title, "]票据"))
	//更新组织用户数据
	OrgUserMod.PushUpdateUserData(args.OrgID, args.UserID)
	//反馈
	return
}

// ArgsAddTickets 批量给予用户票据参数
type ArgsAddTickets struct {
	//组织ID
	// 可选
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//赠送列表
	Data []ArgsAddTicketsChild `json:"data"`
	//可退票据配置列
	CanRefundConfigIDs []int64 `json:"canRefundConfigIDs"`
}

type ArgsAddTicketsChild struct {
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id"`
	//张数
	Count int64 `db:"count" json:"count" check:"int64Than0"`
	//使用来源
	UseFromName string `db:"use_from_name" json:"useFromName"`
}

// AddTickets 批量给予用户票据
func AddTickets(args *ArgsAddTickets) (newTicketIDs []int64, newTicketRefundIDs []int64, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			return
		}
	}()
	tx := Router2SystemConfig.MainDB.MustBegin()
	for _, v := range args.Data {
		//获取票据配置
		var configData FieldsConfig
		err = Router2SystemConfig.MainDB.Get(&configData, "SELECT id, default_expire_time FROM user_ticket_config WHERE id = $1 AND org_id = $2 AND delete_at < to_timestamp(1000000)", v.ConfigID, args.OrgID)
		if err != nil && configData.ID < 1 {
			return
		}
		//生成过期时间
		var expireAt time.Time
		if configData.DefaultExpireTime > 0 {
			expireAt = CoreFilter.GetNowTimeCarbon().AddSeconds(int(configData.DefaultExpireTime)).Time
		}
		//新增票据
		var stmt *sqlx.NamedStmt
		if stmt, err = tx.PrepareNamed("INSERT INTO user_ticket (expire_at, org_id, user_id, config_id, count, res_count) VALUES (:expire_at, :org_id, :user_id, :config_id, :count, :count) RETURNING id;"); err != nil {
			if err2 := tx.Rollback(); err2 != nil {
				err = errors.New(fmt.Sprint("insert user ticket, ", err, ", rollback failed, ", err2))
				return
			}
			return
		}
		//记录添加的票据ID
		var newID int64
		newID, err = CoreSQL.LastRowsAffectedCreate(tx, stmt, map[string]interface{}{
			"expire_at": expireAt,
			"org_id":    args.OrgID,
			"config_id": configData.ID,
			"user_id":   args.UserID,
			"count":     v.Count,
		}, err)
		if err != nil {
			if err2 := tx.Rollback(); err2 != nil {
				err = errors.New(fmt.Sprint("insert user ticket and get last id, ", err, ", rollback failed, ", err2))
				return
			}
			return
		}
		newTicketIDs = append(newTicketIDs, newID)
		//如果配置符合，则添加
		for _, v := range args.CanRefundConfigIDs {
			if v == configData.ID {
				newTicketRefundIDs = append(newTicketRefundIDs, newID)
			}
		}
		//记录日志
		var stmtLog *sqlx.NamedStmt
		stmtLog, err = tx.PrepareNamed("INSERT INTO user_ticket_log (org_id, config_id, user_id, mode, count, des) VALUES (:org_id,:config_id,:user_id,:mode,:count,:des) RETURNING id;")
		if err != nil {
			if err2 := tx.Rollback(); err2 != nil {
				err = errors.New(fmt.Sprint("insert user ticket log, ", err, ", rollback failed, ", err2))
				return
			}
			return
		}
		if _, err = CoreSQL.LastRowsAffectedCreate(tx, stmtLog, map[string]interface{}{
			"org_id":    args.OrgID,
			"config_id": configData.ID,
			"user_id":   args.UserID,
			"mode":      1,
			"count":     v.Count,
			"des":       fmt.Sprint(v.UseFromName, "添加[", v.Count, "]张[", configData.Title, "]票据"),
		}, err); err != nil {
			err = errors.New(fmt.Sprint("get user ticket log last id: ", err))
			if err2 := tx.Rollback(); err2 != nil {
				err = errors.New(fmt.Sprint("insert user ticket log and get last id, ", err, ", rollback failed, ", err2))
				return
			}
			return
		}
	}
	err = tx.Commit()
	//更新组织用户数据
	OrgUserMod.PushUpdateUserData(args.OrgID, args.UserID)
	//反馈
	return
}

// ArgsUseTicket 使用N张票据参数
type ArgsUseTicket struct {
	//ID
	// 如果给予，则查询ID，否则根据configID和userID检索数据
	ID int64 `db:"id" json:"id" check:"id" empty:"true"`
	//组织ID
	// 可选
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//张数
	// 如果是ID，则不能超出ID总数
	Count int64 `db:"count" json:"count" check:"int64Than0"`
	//使用来源
	UseFromName string `db:"use_from_name" json:"useFromName"`
}

// UseTicket 使用N张票据
func UseTicket(args *ArgsUseTicket) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint("recover, ", e))
			return
		}
	}()
	//获取票据配置数据
	var configData FieldsConfig
	configData, err = GetConfigByID(&ArgsGetConfigByID{
		ID:    args.ConfigID,
		OrgID: args.OrgID,
	})
	if err != nil {
		err = errors.New(fmt.Sprint("config not exist, ", err))
		return
	}
	//检查使用限制
	switch configData.LimitTimeType {
	case 1:
		//一次性
		var count int64 = 0
		count, _ = CoreSQL.GetAllCountMap(Router2SystemConfig.MainDB.DB, "user_ticket_log", "id", "config_id = :config_id", map[string]interface{}{
			"config_id": args.ConfigID,
		})
		if count > 0 {
			err = errors.New("too many")
			return
		}
	case 2:
		//日限制
		var count int64 = 0
		count, _ = CoreSQL.GetAllCountMap(Router2SystemConfig.MainDB.DB, "user_ticket_log", "id", "config_id = :config_id AND create_at >= :create_at", map[string]interface{}{
			"config_id": args.ConfigID,
			"create_at": CoreFilter.GetNowTimeCarbon().SubDay().Time,
		})
		if count > 0 {
			err = errors.New("too many")
			return
		}
	case 3:
		//周限制
		var count int64 = 0
		count, _ = CoreSQL.GetAllCountMap(Router2SystemConfig.MainDB.DB, "user_ticket_log", "id", "config_id = :config_id AND create_at >= :create_at", map[string]interface{}{
			"config_id": args.ConfigID,
			"create_at": CoreFilter.GetNowTimeCarbon().SubWeek().Time,
		})
		if count > 0 {
			err = errors.New("too many")
			return
		}
	case 4:
		//月限制
		var count int64 = 0
		count, _ = CoreSQL.GetAllCountMap(Router2SystemConfig.MainDB.DB, "user_ticket_log", "id", "config_id = :config_id AND create_at >= :create_at", map[string]interface{}{
			"config_id": args.ConfigID,
			"create_at": CoreFilter.GetNowTimeCarbon().SubMonth().Time,
		})
		if count > 0 {
			err = errors.New("too many")
			return
		}
	case 5:
		//季度限制
		var count int64 = 0
		count, _ = CoreSQL.GetAllCountMap(Router2SystemConfig.MainDB.DB, "user_ticket_log", "id", "config_id = :config_id AND create_at >= :create_at", map[string]interface{}{
			"config_id": args.ConfigID,
			"create_at": CoreFilter.GetNowTimeCarbon().SubMonths(3).Time,
		})
		if count > 0 {
			err = errors.New("too many")
			return
		}
	case 6:
		//年限制
		var count int64 = 0
		count, _ = CoreSQL.GetAllCountMap(Router2SystemConfig.MainDB.DB, "user_ticket_log", "id", "config_id = :config_id AND create_at >= :create_at", map[string]interface{}{
			"config_id": args.ConfigID,
			"create_at": CoreFilter.GetNowTimeCarbon().SubYear().Time,
		})
		if count > 0 {
			err = errors.New("too many")
			return
		}
	default:
		//不限制
	}
	//使用票据
	if args.ID > 0 {
		var data FieldsTicket
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, config_id, count FROM user_ticket WHERE id = $1 AND ($2 < 1 OR org_id = $2) AND (expire_at < to_timestamp(1000000) OR expire_at >= NOW()) AND delete_at < to_timestamp(1000000)", args.ID, args.OrgID)
		if err != nil {
			err = errors.New(fmt.Sprint("ticket not exist, ticket id: ", args.ID, ", err: ", err))
			return
		}
		if data.Count < args.Count {
			err = errors.New("count not enough")
			return
		}
		if data.Count == args.Count {
			_, err = CoreSQL.DeleteOneSoft(Router2SystemConfig.MainDB.DB, "user_ticket", "id", map[string]interface{}{
				"id": data.ID,
			})
			if err != nil {
				err = errors.New(fmt.Sprint("delete user ticket, ", err))
			}
		} else {
			_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE user_ticket SET count = count - :count WHERE id = :id", map[string]interface{}{
				"id":    data.ID,
				"count": args.Count,
			})
			if err != nil {
				err = errors.New(fmt.Sprint("update user ticket count, ", err))
			}
		}
		//记录日志
		if err == nil {
			_ = appendLog(configData.OrgID, configData.ID, args.UserID, 2, args.Count, fmt.Sprint(args.UseFromName, "使用[", args.Count, "]张[", configData.Title, "]票据"))
		}
	} else {
		var count int64
		count, err = CoreSQL.GetAllSumMap(Router2SystemConfig.MainDB.DB, "user_ticket", "count", "config_id = :config_id AND user_id = :user_id AND (:org_id < 1 OR org_id = :org_id) AND (expire_at < to_timestamp(1000000) OR expire_at >= NOW()) AND delete_at < to_timestamp(1000000)", map[string]interface{}{
			"config_id": args.ConfigID,
			"user_id":   args.UserID,
			"org_id":    args.OrgID,
		})
		if err != nil || count < args.Count {
			err = errors.New(fmt.Sprint("no have ticket count, ", err))
			return
		}
		step := 0
		needCount := args.Count + 0
		for {
			if needCount < 1 {
				break
			}
			var dataList []FieldsTicket
			err = CoreSQL.GetList(Router2SystemConfig.MainDB.DB, &dataList, fmt.Sprint("SELECT id, count FROM user_ticket WHERE config_id = :config_id AND user_id = :user_id AND (:org_id < 1 OR org_id = :org_id) AND (expire_at < to_timestamp(1000000) OR expire_at >= NOW()) AND delete_at < to_timestamp(1000000) ORDER BY expire_at LIMIT 100 OFFSET ", step), map[string]interface{}{
				"config_id": args.ConfigID,
				"user_id":   args.UserID,
				"org_id":    args.OrgID,
			})
			if err != nil {
				break
			}
			if len(dataList) < 1 {
				break
			}
			tx := Router2SystemConfig.MainDB.MustBegin()
			for _, v := range dataList {
				if needCount < 1 {
					break
				}
				if v.Count <= needCount {
					var result sql.Result
					result, err = tx.NamedExec("UPDATE user_ticket SET delete_at = NOW(), count = 0 WHERE id = :id", map[string]interface{}{
						"id": v.ID,
					})
					if err != nil {
						if err2 := tx.Rollback(); err2 != nil {
							err = errors.New("tx update ticket, " + err.Error() + ", rollback, " + err2.Error())
							return
						}
						err = errors.New(fmt.Sprint("tx delete user ticket, ", err))
						return
					}
					err = CoreSQL.LastRowsAffected(tx, result, err)
					if err != nil {
						err = errors.New(fmt.Sprint("tx delete user ticket, get last rows affected, ", err))
						return
					}
					needCount = needCount - v.Count
				} else {
					var result sql.Result
					result, err = tx.NamedExec("UPDATE user_ticket SET count = count - :count WHERE id = :id", map[string]interface{}{
						"id":    v.ID,
						"count": needCount,
					})
					if err != nil {
						if err2 := tx.Rollback(); err2 != nil {
							err = errors.New("tx update ticket, " + err.Error() + ", rollback, " + err2.Error())
							return
						}
						err = errors.New(fmt.Sprint("tx update user ticket, ", err))
						return
					}
					err = CoreSQL.LastRowsAffected(tx, result, err)
					if err != nil {
						err = errors.New(fmt.Sprint("tx update user ticket, get last rows affected, ", err))
						return
					}
					needCount = 0
				}
			}
			err = tx.Commit()
			if err != nil {
				if err2 := tx.Rollback(); err2 != nil {
					err = errors.New("tx update ticket, " + err.Error() + ", rollback, " + err2.Error())
					return
				}
				err = errors.New(fmt.Sprint("tx update user ticket, commit, ", err))
				return
			}
			if needCount < 1 {
				break
			}
			step += 100
		}
		if needCount > 0 {
			err = errors.New(fmt.Sprint("have ticket not use, user id: ", args.UserID, ", config id: ", args.ConfigID, ", need count: ", needCount))
			return
		}
		//记录日志
		var configData FieldsConfig
		configData, err = GetConfigByID(&ArgsGetConfigByID{
			ID:    args.ConfigID,
			OrgID: 0,
		})
		if err != nil {
			err = errors.New(fmt.Sprint("config not exist, ", err))
			return
		}
		_ = appendLog(configData.OrgID, configData.ID, args.UserID, 2, args.Count, fmt.Sprint(args.UseFromName, "使用[", args.Count, "]张[", configData.Title, "]票据"))
	}
	if err != nil {
		return
	}
	//更新组织用户数据
	OrgUserMod.PushUpdateUserData(args.OrgID, args.UserID)
	//反馈
	return
}

// ArgsRefundUseTicket 强制用掉用户的票据参数
type ArgsRefundUseTicket struct {
	//组织ID
	// 可选
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//退还的票据ID列
	IDs pq.Int64Array `db:"ids" json:"ids"`
	//描述信息前缀
	Des string `json:"des" check:"des" min:"1" max:"1000"`
}

// RefundUseTicket 强制用掉用户的票据
// 解决退票问题
func RefundUseTicket(args *ArgsRefundUseTicket) (err error) {
	//强制删除用户的票据
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "user_ticket", "id = :ids AND org_id = :org_id AND user_id = :user_id", map[string]interface{}{
		"ids":     args.IDs,
		"org_id":  args.OrgID,
		"user_id": args.UserID,
	})
	if err != nil {
		return
	}
	_ = appendLog(args.OrgID, 0, args.UserID, 2, 0, fmt.Sprint(args.Des, ", 票据[", args.IDs, "]批量作废"))
	//更新组织用户数据
	OrgUserMod.PushUpdateUserData(args.OrgID, args.UserID)
	//反馈
	return
}

// ArgsClearTicket 清理指定配置的票据参数
type ArgsClearTicket struct {
	//配置ID
	ConfigD int64 `db:"config_id" json:"configID" check:"id"`
	//组织ID
	// 可选
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// ClearTicket 清理指定配置的票据
func ClearTicket(args *ArgsClearTicket) (err error) {
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "user_ticket", "config_id = :config_id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	return
}
