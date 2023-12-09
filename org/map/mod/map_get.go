package OrgMapMod

import (
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// GetMapByID 获取地图数据包
func GetMapByID(id int64) (data FieldsMap) {
	cacheMark := getMapCacheMark(id)
	_ = Router2SystemConfig.MainCache.GetStruct(cacheMark, &data)
	if data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, audit_at, org_id, user_id, parent_id, cover_file_id, name, des, country, province, city, address, map_type, longitude, latitude, ad_count, ad_count_limit, view_time_limit, params FROM org_map WHERE id = $1 LIMIT 1", id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 86400)
	return
}
