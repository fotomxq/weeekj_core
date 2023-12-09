package MapRoom

import (
	AnalysisAny2 "gitee.com/weeekj/weeekj_core/v5/analysis/any2"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"time"
)

// 更新统计信息
func updateRoomAnalysis(orgID int64) {
	var count int64
	//房间总数
	_ = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM map_room WHERE org_id = $1 AND delete_at < to_timestamp(1000000)", orgID)
	AnalysisAny2.AppendData("re", "map_room_count", time.Time{}, orgID, 0, 0, 0, 0, count)
	//按照分类记录入驻的总人数
	// 获取组织下所有可用分类
	type newSortDataType struct {
		SortID int64 `db:"sort_id"`
	}
	var newSortData []newSortDataType
	_ = Router2SystemConfig.MainDB.Select(&newSortData, "select sort_id from map_room where sort_id != 0 and org_id = $1 and delete_at < to_timestamp(1000000) group by sort_id;", orgID)
	if len(newSortData) > 0 {
		for _, vSort := range newSortData {
			//获取该分类下的入住人数
			var vInfoCount int64
			_ = Router2SystemConfig.MainDB.Get(&vInfoCount, "SELECT SUM(array_length(infos, 1)) FROM map_room WHERE org_id = $1 AND sort_id = $2 AND delete_at < to_timestamp(1000000)", orgID, vSort.SortID)
			//记录数据
			AnalysisAny2.AppendData("re", "map_room_info_count", time.Time{}, orgID, 0, vSort.SortID, 0, 0, vInfoCount)
		}
	}
	var infoCount int64
	_ = Router2SystemConfig.MainDB.Get(&infoCount, "SELECT SUM(array_length(infos, 1)) FROM map_room WHERE org_id = $1 AND delete_at < to_timestamp(1000000)", orgID)
	//记录数据
	AnalysisAny2.AppendData("re", "map_room_info_all_count", time.Time{}, orgID, 0, 0, 0, 0, infoCount)
}
