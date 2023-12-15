package UserSub2

import (
	CoreCache "github.com/fotomxq/weeekj_core/v5/core/cache"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// GetConfigByOrg 获取指定组织的配置
func GetConfigByOrg(orgID int64) (data FieldsConfig) {
	cacheMark := getConfigCacheMark(orgID)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = configSQL.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "delete_at", "org_id", "config_data", "params"}).AppendWhere("org_id = $1", orgID).NeedLimit().Result(&data)
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, CoreCache.CacheTime1Day)
	return
}
