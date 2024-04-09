package BaseSystemMission

import (
	"fmt"
	CoreCache "github.com/fotomxq/weeekj_core/v5/core/cache"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

type argsCreateMission struct {
	//组织ID
	// 如果为0则为系统服务
	OrgID int64 `db:"org_id" json:"orgID"`
	//任务名称
	Name string `db:"name" json:"name"`
	//标识码
	Mark string `db:"mark" json:"mark"`
	//计划执行时间
	NextTime string `db:"next_time" json:"nextTime"`
}

func createMission(args *argsCreateMission) (err error) {
	data := getMissionByMark(args.Mark, args.OrgID)
	if data.ID < 1 {
		err = missionDB.Insert().SetFields([]string{"org_id", "name", "mark", "start_at", "next_time"}).Add(map[string]interface{}{
			"org_id":    args.OrgID,
			"name":      args.Name,
			"mark":      args.Mark,
			"start_at":  CoreFilter.GetNowTime(),
			"next_time": args.NextTime,
		}).ExecAndCheckID()
		if err != nil {
			return
		}
	} else {
		err = missionDB.Update().SetFields([]string{"update_at", "org_id", "name", "mark", "start_at", "now_tip", "stop_at", "pause_at", "location", "all_count", "run_count", "run_all_sec", "next_time"}).AddWhereID(data.ID).NamedExec(map[string]interface{}{
			"id":          data.ID,
			"update_at":   CoreFilter.GetNowTime(),
			"org_id":      args.OrgID,
			"name":        args.Name,
			"mark":        args.Mark,
			"start_at":    CoreFilter.GetNowTime(),
			"now_tip":     "",
			"stop_at":     time.Time{},
			"pause_at":    time.Time{},
			"location":    "",
			"all_count":   0,
			"run_count":   0,
			"run_all_sec": 0,
			"next_time":   args.NextTime,
		})
		if err != nil {
			return
		}
		deleteMissionCache(data.ID)
	}
	return
}

func startMission(id int64, nowTip string, location string, allCount int64) (err error) {
	err = missionDB.Update().SetFields([]string{"update_at", "now_tip", "location", "all_count"}).AddWhereID(id).NamedExec(map[string]interface{}{
		"id":        id,
		"update_at": CoreFilter.GetNowTime(),
		"now_tip":   nowTip,
		"location":  location,
		"all_count": allCount,
	})
	if err != nil {
		return
	}
	deleteMissionCache(id)
	return
}

func updateMissionTotal(id int64, allCount int64) (err error) {
	err = missionDB.Update().SetFields([]string{"all_count"}).AddWhereID(id).NamedExec(map[string]interface{}{
		"id":        id,
		"all_count": allCount,
	})
	if err != nil {
		return
	}
	deleteMissionCache(id)
	return
}

func updateMissionAddTotal(id int64, allCount int64) (err error) {
	data := getMission(id)
	data.AllCount += allCount
	err = missionDB.Update().SetFields([]string{"all_count"}).AddWhereID(id).NamedExec(map[string]interface{}{
		"id":        id,
		"all_count": data.AllCount,
	})
	if err != nil {
		return
	}
	deleteMissionCache(id)
	return
}

func updateMission(id int64, nowTip string, location string, runCount int64, runSec int64) (err error) {
	data := getMission(id)
	data.RunCount = data.RunCount + runCount
	data.RunAllSec = data.RunAllSec + runSec
	if data.RunCount > data.AllCount {
		data.RunCount = data.AllCount
	}
	err = missionDB.Update().SetFields([]string{"update_at", "now_tip", "location", "run_count", "run_all_sec", "all_count"}).AddWhereID(id).NamedExec(map[string]interface{}{
		"id":          id,
		"update_at":   CoreFilter.GetNowTime(),
		"now_tip":     nowTip,
		"location":    location,
		"run_count":   data.RunCount,
		"run_all_sec": data.RunAllSec,
		"all_count":   data.AllCount,
	})
	if err != nil {
		return
	}
	deleteMissionCache(id)
	return
}

func pauseMission(id int64) (err error) {
	err = missionDB.Update().SetFields([]string{"pause_at"}).AddWhereID(id).NamedExec(map[string]interface{}{
		"id":       id,
		"pause_at": CoreFilter.GetNowTime(),
	})
	if err != nil {
		return
	}
	deleteMissionCache(id)
	return
}

func stopMission(id int64) (err error) {
	err = missionDB.Update().SetFields([]string{"stop_at"}).AddWhereID(id).NamedExec(map[string]interface{}{
		"id":      id,
		"stop_at": CoreFilter.GetNowTime(),
	})
	if err != nil {
		return
	}
	deleteMissionCache(id)
	CoreNats.PushDataNoErr("base_system_mission_stop", "/base/system_mission/stop", "", id, "", nil)
	return
}

// getMissionByMark 获取指定任务
func getMissionByMark(mark string, orgID int64) (data FieldsMission) {
	cacheMark := getMissionByMarkCacheMark(mark, orgID)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := missionDB.Get().GetByMarkAndOrgID(mark, orgID).Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, CoreCache.CacheTime1Day)
	return
}

// 获取任务
func getMission(id int64) (data FieldsMission) {
	cacheMark := getMissionCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := missionDB.Get().GetByID(id).Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, CoreCache.CacheTime1Day)
	return
}

// 缓冲
func getMissionCacheMark(id int64) string {
	return fmt.Sprint("base:system:mission:id:", id)
}
func getMissionByMarkCacheMark(mark string, orgID int64) string {
	return fmt.Sprint("base:system:mission:mark:", orgID, ".", mark)
}

func deleteMissionCache(id int64) {
	data := getMission(id)
	Router2SystemConfig.MainCache.DeleteMark(getMissionCacheMark(id))
	if data.ID > 0 {
		getMissionByMarkCacheMark(data.Mark, data.OrgID)
	}
}
