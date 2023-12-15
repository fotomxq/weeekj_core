package AnalysisAny

import (
	"errors"
	"fmt"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/golang-module/carbon"
	"time"

	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLTime "github.com/fotomxq/weeekj_core/v5/core/sql/time"
)

// ArgsAppendAny 添加新的记录参数
type ArgsAppendAny struct {
	//创建时间
	// 如果给空，则默认当前时间
	CreateAt string `db:"create_at" json:"createAt"`
	//组织ID
	// 可留空
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	// 可留空
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//绑定ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//扩展参数1
	// 可选扩展参数，默认给0无视
	Param1 int64 `db:"params1" json:"params1" check:"id" empty:"true"`
	//扩展参数2
	// 可选扩展参数，默认给0无视
	Param2 int64 `db:"params2" json:"params2" check:"id" empty:"true"`
	//数据标识码
	Mark string `db:"mark" json:"mark" check:"mark"`
	//数据
	Data    int64  `db:"data" json:"data"`
	DataVal string `db:"data_val" json:"dataVal"`
}

// AppendAny 添加新的记录
// Deprecated
func AppendAny(args *ArgsAppendAny) (err error) {
	//分析当前时间
	var createAt time.Time
	if args.CreateAt == "" {
		createAt = CoreFilter.GetNowTime()
	} else {
		createAt, err = CoreFilter.GetTimeByISO(args.CreateAt)
		if err != nil {
			return
		}
	}
	createAtCarbon := carbon.CreateFromGoTime(createAt)
	//获取配置
	var configData FieldsConfig
	err = Router2SystemConfig.MainDB.Get(&configData, "SELECT id, mqtt_org, mqtt_user, mqtt_user FROM analysis_any_config WHERE mark = $1 AND delete_at < to_timestamp(1000000)", args.Mark)
	if err != nil || configData.ID < 1 {
		err = errors.New(fmt.Sprint("no data, ", err))
		return
	}
	//计算hash
	var newHash string
	newHash, err = CoreFilter.GetSha1ByString(fmt.Sprint(args.Data, args.DataVal))
	if err != nil {
		return
	}
	//获取最新一条数据
	var data FieldsAny
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, hash FROM analysis_any WHERE config_id = $1 AND org_id = $2 AND user_id = $3 AND bind_id = $4 AND params1 = $5 AND params2 = $6 AND create_at >= $7 AND create_at <= $8 ORDER BY id DESC LIMIT 1", configData.ID, args.OrgID, args.UserID, args.BindID, args.Param1, args.Param2, createAtCarbon.SubSecond().Time, createAtCarbon.AddSecond().Time)
	//不需要管是否存在数据
	if err == nil && data.ID > 0 {
		if newHash != data.Hash {
			_, _ = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE analysis_any SET data = :data, data_val = :data_val WHERE id = :id", map[string]interface{}{
				"id":       data.ID,
				"data":     args.Data,
				"data_val": args.DataVal,
			})
		}
	} else {
		//写入数据
		var newData FieldsAny
		err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "analysis_any", "INSERT INTO analysis_any (create_at, org_id, user_id, bind_id, params1, params2, config_id, hash, data, data_val) VALUES (:create_at,:org_id,:user_id,:bind_id,:params1,:params2,:config_id,:hash,:data,:data_val)", map[string]interface{}{
			"create_at": createAt,
			"org_id":    args.OrgID,
			"user_id":   args.UserID,
			"bind_id":   args.BindID,
			"params1":   args.Param1,
			"params2":   args.Param2,
			"config_id": configData.ID,
			"hash":      newHash,
			"data":      args.Data,
			"data_val":  args.DataVal,
		}, &newData)
		if err != nil {
			return
		}
	}
	//如果数据不同，则推送数据
	if !NoMqtt {
		if newHash != data.Hash {
			//更新推送时间
			_, _ = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE analysis_any_config SET last_mqtt = NOW(), last_hash = :last_hash WHERE id = :id", map[string]interface{}{
				"id":        configData.ID,
				"last_hash": newHash,
			})
		}
	}
	//删除缓冲
	cacheMark := fmt.Sprint("analysis:any:get:", args.OrgID, ".", args.Mark, ".", args.UserID, ".", args.BindID, ".", args.Param1, ".", args.Param2)
	Router2SystemConfig.MainCache.DeleteMark(cacheMark)
	Router2SystemConfig.MainCache.DeleteSearchMark(fmt.Sprint("analysis:any:get:", args.OrgID, ".", args.Mark))
	cacheMark2 := fmt.Sprint("analysis:any:sum:", args.OrgID, ".", args.Mark, ".", args.UserID, ".", args.BindID, ".", args.Param1, ".", args.Param2)
	Router2SystemConfig.MainCache.DeleteMark(cacheMark2)
	Router2SystemConfig.MainCache.DeleteSearchMark(fmt.Sprint("analysis:any:sum:", args.OrgID, ".", args.Mark))
	//反馈
	return
}

// ArgsGetAnyByMark 获取指定的记录参数
type ArgsGetAnyByMark struct {
	//组织ID
	// 可留空，用于验证
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	// 可留空
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//数据标识码
	Mark string `db:"mark" json:"mark" check:"mark"`
	//绑定ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//扩展参数1
	// 可选扩展参数，默认给0无视
	Param1 int64 `db:"params1" json:"params1" check:"id" empty:"true"`
	//扩展参数2
	// 可选扩展参数，默认给0无视
	Param2 int64 `db:"params2" json:"params2" check:"id" empty:"true"`
	//时间范围
	BetweenTime CoreSQLTime.DataCoreTime `json:"betweenTime"`
}

type DataGetAnyByMark struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	// 可留空
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	// 可留空
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//绑定ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//扩展参数1
	// 可选扩展参数，默认给0无视
	Param1 int64 `db:"params1" json:"params1" check:"id" empty:"true"`
	//扩展参数2
	// 可选扩展参数，默认给0无视
	Param2 int64 `db:"params2" json:"params2" check:"id" empty:"true"`
	//数据类型
	Mark string `db:"mark" json:"mark"`
	//数据Hash
	Hash string `db:"hash" json:"hash"`
	//数据
	Data    int64  `db:"data" json:"data"`
	DataVal string `db:"data_val" json:"dataVal"`
}

// GetAnyByMark 获取指定的记录
// Deprecated
func GetAnyByMark(args *ArgsGetAnyByMark) (data DataGetAnyByMark, err error) {
	//分析时间范围
	var betweenTime CoreSQLTime.FieldsCoreTime
	betweenTime, err = CoreSQLTime.GetBetweenByISO(args.BetweenTime)
	if err != nil {
		return
	}
	//获取缓存
	cacheMark := fmt.Sprint("analysis:any:sum:", args.OrgID, ".", args.Mark, ".", args.UserID, ".", args.BindID, ".", args.Param1, ".", args.Param2, ".", betweenTime.MinTime.Format("20060102150405"), ".", betweenTime.MaxTime.Format("20060102150405"))
	if err = Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil {
		if data.Mark != "" {
			return
		}
		//CoreLog.Info("debug: mark: ", data.Mark, ", val: ", data.Data)
	}
	//获取配置
	var configID int64
	err = Router2SystemConfig.MainDB.Get(&configID, "SELECT id FROM analysis_any_config WHERE mark = $1 AND delete_at < to_timestamp(1000000)", args.Mark)
	if err != nil || configID < 1 {
		err = errors.New(fmt.Sprint("no data, ", err))
		return
	}
	//获取数据
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, org_id, user_id, bind_id, params1, params2, hash, data, data_val FROM analysis_any WHERE ($1 < 1 OR org_id = $1) AND ($2 < 1 OR user_id = $2) AND config_id = $3 AND create_at >= $4 AND create_at < $5 AND ($6 < 0 OR bind_id = $6) AND ($7 < 0 OR params1 = $7) AND ($8 < 0 OR params2 = $8) ORDER BY id DESC LIMIT 1", args.OrgID, args.UserID, configID, betweenTime.MinTime, betweenTime.MaxTime, args.BindID, args.Param1, args.Param2)
	if err != nil || data.ID < 1 {
		//从归档数据抽取
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, org_id, user_id, bind_id, params1, params2, hash, data, data_val FROM analysis_any_file WHERE ($1 < 1 OR org_id = $1) AND ($2 < 1 OR user_id = $2) AND config_id = $3 AND create_at >= $4 AND create_at < $5  AND ($6 < 0 OR bind_id = $6) AND ($7 < 0 OR params1 = $7) AND ($8 < 0 OR params2 = $8) ORDER BY id DESC LIMIT 1", args.OrgID, args.UserID, configID, betweenTime.MinTime, betweenTime.MaxTime, args.BindID, args.Param1, args.Param2)
		if err == nil && data.ID < 1 {
			err = errors.New(fmt.Sprint("no data, ", err))
			return
		}
	}
	//反馈数据
	data.Mark = args.Mark
	//写入缓冲
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 3600)
	//反馈
	return
}

// GetAnyInt64ByMark 单独获取数据的统计数字
// Deprecated
func GetAnyInt64ByMark(args *ArgsGetAnyByMark) (count int64, err error) {
	var data DataGetAnyByMark
	data, err = GetAnyByMark(args)
	if err != nil {
		return
	}
	count = data.Data
	return
}

// Deprecated
func GetAnyInt64ByMarkNoErr(args *ArgsGetAnyByMark) (count int64) {
	count, _ = GetAnyInt64ByMark(args)
	return
}

type DataGetAnySumByMark struct {
	//影响行数
	IDCount int64 `db:"id_count" json:"idCount"`
	//创建时间
	CreateMinAt time.Time `db:"create_min_at" json:"createMinAt"`
	CreateMaxAt time.Time `db:"create_max_at" json:"createMaxAt"`
	//组织ID
	// 可留空
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	// 可留空
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//绑定ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//扩展参数1
	// 可选扩展参数，默认给0无视
	Param1 int64 `db:"params1" json:"params1" check:"id" empty:"true"`
	//扩展参数2
	// 可选扩展参数，默认给0无视
	Param2 int64 `db:"params2" json:"params2" check:"id" empty:"true"`
	//数据类型
	Mark string `db:"mark" json:"mark"`
	//数据
	Data    int64  `db:"data" json:"data"`
	DataVal string `db:"data_val" json:"dataVal"`
}

// GetAnySumByMark 获取指定的记录(同一个阶段的合计数)
// Deprecated
func GetAnySumByMark(args *ArgsGetAnyByMark) (data DataGetAnySumByMark, err error) {
	//分析时间范围
	var betweenTime CoreSQLTime.FieldsCoreTime
	betweenTime, err = CoreSQLTime.GetBetweenByISO(args.BetweenTime)
	if err != nil {
		return
	}
	//获取缓存
	cacheMark := fmt.Sprint("analysis:any:sum:", args.OrgID, ".", args.Mark, ".", args.UserID, ".", args.BindID, ".", args.Param1, ".", args.Param2, ".", betweenTime.MinTime.Format("20060102150405"), ".", betweenTime.MaxTime.Format("20060102150405"))
	if err = Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil {
		if data.Mark != "" {
			return
		}
		//CoreLog.Info("data: ", data.Mark, ":", data.Data)
	}
	//获取配置
	var configID int64
	err = Router2SystemConfig.MainDB.Get(&configID, "SELECT id FROM analysis_any_config WHERE mark = $1 AND delete_at < to_timestamp(1000000)", args.Mark)
	if err != nil || configID < 1 {
		err = errors.New(fmt.Sprint("no data, ", err))
		return
	}
	//获取数据
	_ = Router2SystemConfig.MainDB.Get(&data, "select COUNT(id) as id_count, MIN(create_at) as create_min_at, MAX(create_at) as create_max_at, SUM(data) as data, MAX(data_val) as data_val from analysis_any WHERE ($1 < 1 OR org_id = $1) AND ($2 < 1 OR user_id = $2) AND config_id = $3 AND create_at >= $4 AND create_at < $5 AND ($6 < 0 OR bind_id = $6) AND ($7 < 0 OR params1 = $7) AND ($8 < 0 OR params2 = $8) ORDER BY id DESC LIMIT 1", args.OrgID, args.UserID, configID, betweenTime.MinTime, betweenTime.MaxTime, args.BindID, args.Param1, args.Param2)
	//抽取归档数据
	var data2 DataGetAnySumByMark
	_ = Router2SystemConfig.MainDB.Get(&data2, "select COUNT(id) as id_count, MIN(create_at) as create_min_at, MAX(create_at) as create_max_at, SUM(data) as data, MAX(data_val) as data_val from analysis_any_file WHERE ($1 < 1 OR org_id = $1) AND ($2 < 1 OR user_id = $2) AND config_id = $3 AND create_at >= $4 AND create_at < $5  AND ($6 < 0 OR bind_id = $6) AND ($7 < 0 OR params1 = $7) AND ($8 < 0 OR params2 = $8) ORDER BY id DESC LIMIT 1", args.OrgID, args.UserID, configID, betweenTime.MinTime, betweenTime.MaxTime, args.BindID, args.Param1, args.Param2)
	//合并数据
	data.Data = data.Data + data2.Data
	if data.DataVal == "" && data2.DataVal != "" {
		data.DataVal = data2.DataVal
	}
	//补全参数
	data.OrgID = args.OrgID
	data.UserID = args.UserID
	data.BindID = args.BindID
	data.Param1 = args.Param1
	data.Param2 = args.Param2
	//反馈数据
	data.Mark = args.Mark
	//写入缓冲
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 3600)
	//反馈
	return
}

// CheckAnyHaveData 检查条件下是否存在数据？
// Deprecated
func CheckAnyHaveData(args *ArgsGetAnyByMark) (haveData bool) {
	//获取配置
	var configID int64
	err := Router2SystemConfig.MainDB.Get(&configID, "SELECT id FROM analysis_any_config WHERE mark = $1 AND delete_at < to_timestamp(1000000)", args.Mark)
	if err != nil || configID < 1 {
		err = errors.New(fmt.Sprint("no data, ", err))
		return
	}
	//分析时间范围
	var betweenTime CoreSQLTime.FieldsCoreTime
	betweenTime, err = CoreSQLTime.GetBetweenByISO(args.BetweenTime)
	if err != nil {
		return
	}
	//获取数据
	var data int64
	err = Router2SystemConfig.MainDB.Get(&data, "select COUNT(id) as id_count from analysis_any WHERE ($1 < 1 OR org_id = $1) AND ($2 < 1 OR user_id = $2) AND config_id = $3 AND create_at >= $4 AND create_at < $5 AND ($6 < 0 OR bind_id = $6) AND ($7 < 0 OR params1 = $7) AND ($8 < 0 OR params2 = $8) ORDER BY id DESC LIMIT 1", args.OrgID, args.UserID, configID, betweenTime.MinTime, betweenTime.MaxTime, args.BindID, args.Param1, args.Param2)
	if err == nil && data < 1 {
		//从归档数据抽取
		err = Router2SystemConfig.MainDB.Get(&data, "select COUNT(id) as id_count from analysis_any_file WHERE ($1 < 1 OR org_id = $1) AND ($2 < 1 OR user_id = $2) AND config_id = $3 AND create_at >= $4 AND create_at < $5  AND ($6 < 0 OR bind_id = $6) AND ($7 < 0 OR params1 = $7) AND ($8 < 0 OR params2 = $8) ORDER BY id DESC LIMIT 1", args.OrgID, args.UserID, configID, betweenTime.MinTime, betweenTime.MaxTime, args.BindID, args.Param1, args.Param2)
		if err == nil && data < 1 {
			return
		}
	}
	haveData = data > 0
	return
}
