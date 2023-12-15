package OrgCoreCoreMod

import (
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// GetBind 查看绑定关系
func GetBind(id int64, orgID int64, userID int64) (data FieldsBind) {
	data = getBindByID(id)
	if data.ID < 1 || !CoreFilter.EqID2(orgID, data.OrgID) || !CoreFilter.EqID2(userID, data.UserID) {
		data = FieldsBind{}
		return
	}
	return
}

// 获取指定ID
func getBindByID(id int64) (data FieldsBind) {
	cacheMark := getBindCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, last_at, user_id, name, org_id, group_ids, manager, params, avatar, nation_code, phone, email, sync_system, sync_id, sync_hash, role_config_ids FROM org_core_bind WHERE id = $1", id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, bindCacheTime)
	return
}

// 缓冲
func getBindCacheMark(id int64) string {
	return fmt.Sprint("org:core:bind:id:", id)
}

func deleteBindCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getBindCacheMark(id))
}
