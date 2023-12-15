package OrgMission

import (
	"errors"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNextTime "github.com/fotomxq/weeekj_core/v5/core/next_time"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/golang-module/carbon"
)

func runAuto() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("org mission auto run, ", r)
		}
	}()
	//批量获取待处理自动化
	limit := 100
	step := 0
	for {
		var dataList []FieldsAuto
		if err := Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, create_at, update_at, delete_at, time_type, time_n, skip_holiday, start_hour, start_minute, end_hour, end_minute, org_id, create_bind_id, bind_id, other_bind_ids, title, des, des_files, start_at, next_at, end_at, tip_id, level, sort_id, tags, params FROM org_mission_auto WHERE delete_at < to_timestamp(1000000) AND start_at < NOW() AND end_at >= NOW() AND next_at <= NOW() LIMIT $1 OFFSET $2", limit, step); err != nil {
			break
		}
		if len(dataList) < 1 {
			break
		}
		for _, v := range dataList {
			if err := runAutoChild(&v); err != nil {
				CoreLog.Error("org mission auto run, ", err)
			}
		}
		step += limit
	}
}

func runAutoChild(autoData *FieldsAuto) (err error) {
	//检查该计划是否存在关联任务？
	var missionID int64
	if err = Router2SystemConfig.MainDB.Get(&missionID, "SELECT id FROM org_mission WHERE auto_id = $1 AND delete_at < NOW() AND status = 0"); err == nil && missionID > 0 {
		//存在还未完成任务，则跳出
		return
	}
	//计算下一次执行的时间
	nextAt := carbon.CreateFromTimestamp(autoData.NextAt.Unix())
	var b, needDeleteConfig bool
	nextAt, needDeleteConfig, b = CoreNextTime.MakeNextAt(autoData.TimeType, autoData.TimeN, autoData.SkipHoliday, nextAt)
	if !b {
		return
	}
	if needDeleteConfig {
		if _, err = CoreSQL.DeleteOneSoft(Router2SystemConfig.MainDB.DB, "org_mission_auto", "id", map[string]interface{}{
			"id": autoData.ID,
		}); err != nil {
			CoreLog.Error("org mission auto, delete schedule data, " + err.Error())
			return
		}
	}
	startAt := nextAt.SetHour(autoData.StartHour).SetMinute(autoData.StartMinute).SetSecond(0)
	endAt := nextAt.SetHour(autoData.EndHour).SetMinute(autoData.EndMinute).SetSecond(0)
	//建立任务
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO org_mission (auto_id, status, org_id, create_bind_id, bind_id, other_bind_ids, title, des, des_files, start_at, end_at, tip_id, parent_id, level, sort_id, tags, params) VALUES (:auto_id,0,:org_id,:create_bind_id,:bind_id,:other_bind_ids,:title,:des,:des_files,:start_at,:end_at,:tip_id,0,:level,:sort_id,:tags,:params)", map[string]interface{}{
		"auto_id":        autoData.ID,
		"org_id":         autoData.OrgID,
		"create_bind_id": autoData.CreateBindID,
		"bind_id":        autoData.BindID,
		"other_bind_ids": autoData.OtherBindIDs,
		"title":          autoData.Title,
		"des":            autoData.Des,
		"des_files":      autoData.DesFiles,
		"start_at":       startAt.Time,
		"end_at":         endAt.Time,
		"tip_id":         autoData.TipID,
		"level":          autoData.Level,
		"sort_id":        autoData.SortID,
		"tags":           autoData.Tags,
		"params":         autoData.Params,
	})
	if err != nil {
		err = errors.New("create mission, " + err.Error())
		return
	}
	//更新下一次执行时间
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_mission_auto SET update_at = NOW(), next_at = :next_at WHERE id = :id", map[string]interface{}{
		"id":      autoData.ID,
		"next_at": nextAt,
	})
	if err != nil {
		err = errors.New("update mission auto next at, " + err.Error())
		return
	}
	return
}
