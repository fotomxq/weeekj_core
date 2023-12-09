package BlogCoreMod

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// GetContentByIDNoErr 无错误获取文章信息
func GetContentByIDNoErr(id int64, orgID int64) (data FieldsContent) {
	data = getContentID(id)
	if data.ID < 1 || !CoreFilter.EqID2(orgID, data.OrgID) {
		data = FieldsContent{}
		return
	}
	return
}

// 获取文章
func getContentID(id int64) (data FieldsContent) {
	cacheMark := getContentCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, content_type, audit_at, audit_des, org_id, user_id, bind_id, param1, param2, param3, visit_count, key, parent_id, publish_at, is_top, sort_id, tags, title, title_des, cover_file_id, des_files, des, params FROM blog_core_content WHERE id = $1", id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 604800)
	return
}
