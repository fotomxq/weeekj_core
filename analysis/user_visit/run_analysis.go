package AnalysisUserVisit

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/golang-module/carbon"
)

//runAnalysis 计算数据包
/**
该模块用于计算：
 所有数据，以最开始存在数据时启动计算，直到数据截止日期时。
1/ 用户进入和离开间隔，如果2次没有找到离开数据，将交给移动数据包处理
2/ 配合用户移动数据，来分析用户可能停留的时间长度
3/ 统计不同模块下、不同时间段下，访客的访问数量和停留时间平均时长
*/
func runAnalysis() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("analysis user visit run, analysis, ", r)
		}
	}()
	//所有数据核对上一个小时开始到结束的数据
	startAt := CoreFilter.GetNowTimeCarbon().SubHour().StartOfHour()
	endAt := CoreFilter.GetNowTimeCarbon().SubHour().EndOfHour()
	//计算用户的访问停留时间
	runAnalysisWaitTime(startAt, endAt)
	//点击重要按钮次数
	runAnalysisClick(startAt, endAt)
	//点击购物按钮次数
	runAnalysisClickBuy(startAt, endAt)
	//潜在黑客行为次数
	runAnalysisHack(startAt, endAt)
}

func runAnalysisWaitTime(startAt, endAt carbon.Carbon) {
	var dataList []FieldsVisit
	if err := Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, create_at, org_id, create_info, user_id, country, phone, ip, mark, action, params FROM analysis_user_visit WHERE create_at >= $1 AND create_at <= $2 AND action = 'insert'", startAt.Time, endAt.Time); err != nil {
		return
	}
	if len(dataList) < 1 {
		return
	}
	for _, v := range dataList {
		//进入人次
		if err := CreateCount(&ArgsCreateCount{
			OrgID: v.OrgID,
			Mark:  200,
			Count: 1,
		}); err != nil {
			CoreLog.Error("analysis user visit run, analysis, add insert count, err: ", err)
		}
		//获取该客户最早离开时间
		var lastOut FieldsVisit
		if err := Router2SystemConfig.MainDB.Get(&lastOut, "SELECT id, create_at FROM analysis_user_visit WHERE org_id = $1 AND create_at > $2 AND user_id = $3 AND phone = $4 AND ip = $5 AND action = 'out' ORDER BY id LIMIT 1", v.OrgID, v.CreateAt, v.UserID, v.Phone, v.IP); err != nil {
			continue
		}
		//获取数据
		var newData FieldsWaitTime
		if err := Router2SystemConfig.MainDB.Get(&newData, "SELECT id FROM analysis_user_wait_time WHERE org_id = $1 AND system = $2 AND from_mark = $3 AND from_id = $4 AND create_at >= $5 AND create_at <= $6", v.OrgID, v.CreateInfo.System, v.CreateInfo.Mark, v.CreateInfo.ID, startAt.Time, endAt.Time); err == nil && newData.ID > 0 {
			_, err := CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE analysis_user_wait_time SET count = count + 1, wait_time = wait_time + :wait_time WHERE id = :id", map[string]interface{}{
				"id":        newData.ID,
				"wait_time": lastOut.CreateAt.Unix() - v.CreateAt.Unix(),
			})
			if err != nil {
				CoreLog.Error("analysis user visit run, analysis, update id: ", newData.ID, ", err: ", err)
			}
			continue
		}
		_, err := CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO analysis_user_wait_time (create_at, org_id, system, from_mark, from_id, count, wait_time) VALUES (:create_at,:org_id,:system,:from_mark,:from_id,:count,:wait_time)", map[string]interface{}{
			"create_at": endAt.Time,
			"org_id":    v.OrgID,
			"system":    v.CreateInfo.System,
			"from_mark": v.CreateInfo.Mark,
			"from_id":   v.CreateInfo.ID,
			"count":     1,
			"wait_time": lastOut.CreateAt.Unix() - v.CreateAt.Unix(),
		})
		if err != nil {
			CoreLog.Error("analysis user visit run, analysis, create new data, err: ", err)
			continue
		}
	}
}

func runAnalysisClick(startAt, endAt carbon.Carbon) {
	var dataList []FieldsVisit
	if err := Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, create_at, org_id, create_info, user_id, country, phone, ip, mark, action, params FROM analysis_user_visit WHERE create_at >= $1 AND create_at <= $2 AND action = 'click'", startAt.Time, endAt.Time); err != nil {
		return
	}
	if len(dataList) < 1 {
		return
	}
	for _, v := range dataList {
		//记录次数
		if err := CreateCount(&ArgsCreateCount{
			OrgID: v.OrgID,
			Mark:  201,
			Count: 1,
		}); err != nil {
			CoreLog.Error("analysis user visit run, analysis, add click count, err: ", err)
		}
	}
}

func runAnalysisClickBuy(startAt, endAt carbon.Carbon) {
	var dataList []FieldsVisit
	if err := Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, create_at, org_id, create_info, user_id, country, phone, ip, mark, action, params FROM analysis_user_visit WHERE create_at >= $1 AND create_at <= $2 AND action = 'click_buy'", startAt.Time, endAt.Time); err != nil {
		return
	}
	if len(dataList) < 1 {
		return
	}
	for _, v := range dataList {
		//记录次数
		if err := CreateCount(&ArgsCreateCount{
			OrgID: v.OrgID,
			Mark:  202,
			Count: 1,
		}); err != nil {
			CoreLog.Error("analysis user visit run, analysis, add click buy count, err: ", err)
		}
	}
}

func runAnalysisHack(startAt, endAt carbon.Carbon) {
	var dataList []FieldsVisit
	if err := Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, create_at, org_id, create_info, user_id, country, phone, ip, mark, action, params FROM analysis_user_visit WHERE create_at >= $1 AND create_at <= $2 AND action = 'ha'", startAt.Time, endAt.Time); err != nil {
		return
	}
	if len(dataList) < 1 {
		return
	}
	for _, v := range dataList {
		//记录次数
		if err := CreateCount(&ArgsCreateCount{
			OrgID: v.OrgID,
			Mark:  203,
			Count: 1,
		}); err != nil {
			CoreLog.Error("analysis user visit run, analysis, add hack count, err: ", err)
		}
	}
}
