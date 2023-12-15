package ServiceInfoExchangeMod

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetInfoID 获取指定信息参数
type ArgsGetInfoID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
}

// GetInfoID 获取指定信息
func GetInfoID(id, orgID, userID int64) (data FieldsInfo) {
	data = getInfoByID(id)
	if data.ID < 1 || !CoreFilter.EqID2(orgID, data.OrgID) || !CoreFilter.EqID2(userID, data.UserID) {
		data = FieldsInfo{}
		return
	}
	return
}

// 获取信息交互数据
func getInfoByID(id int64) (data FieldsInfo) {
	cacheMark := getInfoCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, publish_at, audit_at, audit_des, info_type, org_id, user_id, sort_id, tags, title, title_des, des, cover_file_ids, currency, price, order_id, wait_order_id, order_finish, address, params, expire_at, limit_count FROM service_info_exchange WHERE id = $1", id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 86400)
	return
}
