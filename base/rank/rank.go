package BaseRank

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

// 获取列队列表
type ArgsGetRankList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//服务标识码
	ServiceMark string `json:"serviceMark"`
	//任务标识码
	MissionMark string `json:"missionMark"`
}

func GetRankList(args *ArgsGetRankList) (dataList []FieldsRank, dataCount int64, err error) {
	where := "service_mark=:service_mark AND mission_mark=:mission_mark"
	maps := map[string]interface{}{
		"service_mark": args.ServiceMark,
		"mission_mark": args.MissionMark,
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"core_rank",
		"id",
		fmt.Sprint(
			"SELECT id, create_at, expire_at, pick_min, pick_at, service_mark, mission_mark, mission_data FROM core_rank WHERE ",
			where,
		),
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "expire_at", "pick_at"},
	)
	return
}

// 获取完成列队
type ArgsGetRankOverList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//服务标识码
	ServiceMark string `json:"serviceMark"`
	//任务标识码
	MissionMark string `json:"missionMark"`
	//任务参数
	MissionData []byte `json:"missionData"`
}

func GetRankOverList(args *ArgsGetRankOverList) (dataList []FieldsRankOver, dataCount int64, err error) {
	where := "service_mark=:service_mark AND mission_mark=:mission_mark AND mission_data=:mission_data"
	maps := map[string]interface{}{
		"service_mark": args.ServiceMark,
		"mission_mark": args.MissionMark,
		"mission_data": args.MissionData,
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"core_rank_over",
		"id",
		fmt.Sprint(
			"SELECT id, create_at, expire_at, service_mark, mission_mark, mission_data FROM core_rank_over WHERE ",
			where,
		),
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "expire_at"},
	)
	return
}

// 写入列队
type ArgsAppendRank struct {
	//服务标识码
	ServiceMark string `json:"serviceMark"`
	//过期时间
	ExpireAt time.Time `json:"expireAt"`
	//提取最短间隔 s
	PickMin int64 `json:"pickMin"`
	//任务标识码
	MissionMark string `json:"missionMark"`
	//任务内容
	MissionData []byte `json:"missionData"`
}

func AppendRank(args *ArgsAppendRank) (data FieldsRank, err error) {
	var lastID int64
	lastID, err = CoreSQL.CreateOneAndID(Router2SystemConfig.MainDB.DB, "INSERT INTO core_rank(expire_at, pick_min, pick_at, service_mark, mission_mark, mission_data) VALUES(:expire_at, :pick_min, :pick_at, :service_mark, :mission_mark, :mission_data)", map[string]interface{}{
		"expire_at":    args.ExpireAt,
		"pick_min":     args.PickMin,
		"pick_at":      time.Time{},
		"service_mark": args.ServiceMark,
		"mission_mark": args.MissionMark,
		"mission_data": args.MissionData,
	})
	if err == nil {
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, expire_at, pick_min, pick_at, service_mark, mission_mark, mission_data FROM core_rank WHERE id=$1", lastID)
	}
	return
}

// 抽取一组数据用于计算
type ArgsPickRank struct {
	//服务标识码
	ServiceMark string `json:"serviceMark"`
	//任务标识码
	// 可以指定空任务标识码，可以获得该服务下所有任务
	MissionMark string `json:"missionMark"`
	//要提取几个
	Max int64 `json:"max"`
}

func PickRank(args *ArgsPickRank) (dataList []FieldsRank, err error) {
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT * FROM core_rank WHERE expire_at >= NOW() AND pick_at <= NOW() AND service_mark=$1 ORDER BY id LIMIT $2", args.ServiceMark, args.Max)
	if err != nil {
		return
	}
	for _, v := range dataList {
		var appendTime time.Time
		appendTime, err = CoreFilter.GetTimeByAdd(fmt.Sprint(v.PickMin, "s"))
		if err != nil {
			continue
		}
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE core_rank SET pick_at=:pick_at WHERE id=:id", map[string]interface{}{
			"pick_at": appendTime,
			"id":      v.ID,
		})
		if err != nil {
			return
		}
	}
	if err != nil {
		err = nil
	}
	return
}

// 标记完成
type ArgsOverRank struct {
	//ID
	ID int64 `json:"id"`
	//完成结果
	Result []byte `json:"result"`
}

func OverRank(args *ArgsOverRank) (err error) {
	var data FieldsRank
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, expire_at, pick_min, pick_at, service_mark, mission_mark, mission_data FROM core_rank WHERE id=$1", args.ID)
	if err != nil {
		return
	}
	err = overRankByData(&data, args.Result)
	return
}

// 直接终结数据
func overRankByData(data *FieldsRank, result []byte) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint("recover, ", e))
			return
		}
	}()
	tx := Router2SystemConfig.MainDB.MustBegin()
	_, err = tx.NamedExec("DELETE FROM core_rank WHERE id=:id", map[string]interface{}{
		"id": data.ID,
	})
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			err = errors.New("rollback failed: " + err.Error() + ", err: " + err2.Error())
			return
		}
		return
	}
	_, err = tx.NamedExec("INSERT INTO core_rank_over(create_at, expire_at, over_at, result, service_mark, mission_mark, mission_data) VALUES (:create_at, :expire_at, NOW(), :result, :service_mark, :mission_mark, :mission_data)", map[string]interface{}{
		"create_at":    data.CreateAt,
		"expire_at":    data.ExpireAt,
		"result":       result,
		"service_mark": data.ServiceMark,
		"mission_mark": data.MissionMark,
		"mission_data": data.MissionData,
	})
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			err = errors.New("rollback failed: " + err.Error() + ", err: " + err2.Error())
			return
		}
		return
	}
	err = tx.Commit()
	return
}
