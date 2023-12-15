package FinancePhysicalPay

import (
	"database/sql"
	"errors"
	"fmt"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	CoreSQLTime "github.com/fotomxq/weeekj_core/v5/core/sql/time"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// ArgsGetLogList 获取请求列表参数
type ArgsGetLogList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//成员ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//服务用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//支付分发渠道
	// order 订单 / tms 配送 / housekeeping 家政服务
	System string `db:"system" json:"system" check:"mark" empty:"true"`
	//标的物
	PhysicalID int64 `db:"physical_id" json:"physicalID" check:"id" empty:"true"`
	//时间范围
	// 部分统计支持
	TimeBetween CoreSQLTime.DataCoreTime `json:"timeBetween"`
}

// GetLogList 获取请求列表
func GetLogList(args *ArgsGetLogList) (dataList []FieldsLog, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.OrgID > -1 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.BindID > -1 {
		where = where + " AND (bind_id = :bind_id OR :bind_id = ANY(other_binds))"
		maps["bind_id"] = args.BindID
	}
	if args.UserID > -1 {
		where = where + " AND user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.System != "" {
		where = where + " AND system = :system"
		maps["system"] = args.System
	}
	if args.PhysicalID > -1 {
		where = where + " AND physical_id = :physical_id"
		maps["physical_id"] = args.PhysicalID
	}
	if args.TimeBetween.MinTime != "" && args.TimeBetween.MaxTime != "" {
		var timeBetween CoreSQLTime.FieldsCoreTime
		timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
		if err != nil {
			return
		}
		where = where + " AND create_at >= :start_at AND create_at <= :end_at"
		maps["start_at"] = timeBetween.MinTime
		maps["end_at"] = timeBetween.MaxTime
	}
	if where == "" {
		where = "true"
	}
	tableName := "finance_physical_pay_log"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, org_id, bind_id, user_id, physical_id, physical_count, bind_from, bind_count, params FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	return
}

// ArgsCreateLog 创建新的请求参数
type ArgsCreateLog struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//配送人员
	BindID int64 `db:"bind_id" json:"bindID" check:"id"`
	//服务用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//支付分发渠道
	// order 订单 / tms 配送 / housekeeping 家政服务
	System string `db:"system" json:"system" check:"mark"`
	//置换数量
	Data []ArgsCreateLogData `json:"data"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

type ArgsCreateLogData struct {
	//给予标的物数量
	PhysicalCount int64 `db:"physical_count" json:"physicalCount" check:"int64Than0"`
	//货物来源标的物
	// 如果没有标的物ID，则通过此数据获取配置；否则会报错
	BindFrom CoreSQLFrom.FieldsFrom `db:"bind_from" json:"bindFrom"`
	//置换商品的数量
	BindCount int64 `db:"bind_count" json:"bindCount" check:"int64Than0"`
}

// CreateLog 创建新的请求
func CreateLog(args *ArgsCreateLog) (newLogIDs pq.Int64Array, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint("recover, ", e))
			return
		}
	}()
	//锁定机制
	logLock.Lock()
	defer logLock.Unlock()
	//获取置换物配置
	tx := Router2SystemConfig.MainDB.MustBegin()
	for _, v := range args.Data {
		var physicalData FieldsPhysical
		if v.BindFrom.System == "" {
			err = errors.New("no bind from")
			if err2 := tx.Rollback(); err2 != nil {
				err = errors.New(err.Error() + ", rollback failed: " + err2.Error())
				return
			}
			return
		}
		physicalData, err = GetPhysicalByFrom(&ArgsGetPhysicalByFrom{
			OrgID:    args.OrgID,
			BindFrom: v.BindFrom,
		})
		if err != nil {
			err = errors.New(fmt.Sprint("no physical data by bind from, ", err))
			if err2 := tx.Rollback(); err2 != nil {
				err = errors.New(err.Error() + ", rollback failed: " + err2.Error())
				return
			}
			return
		}
		if physicalData.ID < 1 {
			err = errors.New("no physical data")
			if err2 := tx.Rollback(); err2 != nil {
				err = errors.New(err.Error() + ", rollback failed: " + err2.Error())
				return
			}
			return
		}
		//检查本次置换物品是否超出限制
		if v.PhysicalCount+physicalData.TakeCount > physicalData.LimitCount {
			err = errors.New("more than limit")
			if err2 := tx.Rollback(); err2 != nil {
				err = errors.New(err.Error() + ", rollback failed: " + err2.Error())
				return
			}
			return
		}
		//检查置换商品所需的置换物，必须完全符合条件
		if v.PhysicalCount != physicalData.NeedCount*v.BindCount {
			err = errors.New("physical count not equal")
			if err2 := tx.Rollback(); err2 != nil {
				err = errors.New(err.Error() + ", rollback failed: " + err2.Error())
				return
			}
			return
		}
		//递减处理
		var result sql.Result
		result, err = tx.NamedExec("UPDATE finance_physical_pay_physical SET update_at = NOW(), take_count = take_count + :take_count WHERE id = :id AND limit_count > take_count + :take_count AND delete_at < to_timestamp(1000000)", map[string]interface{}{
			"id":         physicalData.ID,
			"take_count": v.PhysicalCount,
		})
		if err != nil {
			err = CoreSQL.LastRowsAffected(tx, result, errors.New(fmt.Sprint("update physical count limit or failed, ", err)))
			return
		}
		//创建记录
		var stmt *sqlx.NamedStmt
		stmt, err = tx.PrepareNamed("INSERT INTO finance_physical_pay_log(org_id, bind_id, user_id, system, physical_id, physical_count, bind_from, bind_count, params) VALUES(:org_id, :bind_id, :user_id, :system, :physical_id, :physical_count, :bind_from, :bind_count, :params) RETURNING id;")
		if err != nil {
			err = errors.New(fmt.Sprint("create log, ", err))
			return
		}
		var newID int64
		newID, err = CoreSQL.LastRowsAffectedCreate(tx, stmt, map[string]interface{}{
			"org_id":         args.OrgID,
			"bind_id":        args.BindID,
			"user_id":        args.UserID,
			"system":         args.System,
			"physical_id":    physicalData.ID,
			"physical_count": v.PhysicalCount,
			"bind_from":      v.BindFrom,
			"bind_count":     v.BindCount,
			"params":         args.Params,
		}, err)
		if err != nil {
			return
		}
		newLogIDs = append(newLogIDs, newID)
	}
	err = tx.Commit()
	if err != nil {
		err = CoreSQL.LastRows(tx, err)
		return
	}
	return
}
