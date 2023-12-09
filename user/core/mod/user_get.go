package UserCoreMod

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// GetAllUserList 获取所有用户
func GetAllUserList(orgID int64, status int64, sortID int64, tags pq.Int64Array, step, limit int) (dataList []FieldsUserType) {
	where := ""
	maps := map[string]interface{}{
		"limit": limit,
		"step":  step,
	}
	where = CoreSQL.GetDeleteSQL(false, where)
	if status > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = orgID
	}
	if status > -1 {
		where = where + " AND status = :status"
		maps["status"] = status
	}
	if sortID > 0 {
		where = where + " AND sort_id = :sort_id"
		maps["sort_id"] = sortID
	}
	if len(tags) > 0 {
		where = where + " AND tags @> :tags"
		maps["tags"] = tags
	}
	var rawList []FieldsUserType
	_ = CoreSQL.GetList(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		"SELECT id FROM user_core WHERE "+where+" LIMIT :limit OFFSET :step",
		maps,
	)
	for _, v := range rawList {
		vData := getUserByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

func GetUserByID(userID int64, orgID int64) (data FieldsUserType) {
	data = getUserByID(userID)
	if data.ID < 1 || !CoreFilter.EqID2(orgID, data.OrgID) {
		data = FieldsUserType{}
		return
	}
	return
}

func GetUserNiceNameByID(userID int64) string {
	data := getUserByID(userID)
	return data.Name
}

// 获取指定用户数据
func getUserByID(userID int64) (data FieldsUserType) {
	cacheMark := getUserCacheMark(userID)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, status, org_id, name, password, nation_code, phone, email, username, avatar, parents, groups, infos, logins, sort_id, tags, phone_verify, email_verify FROM user_core WHERE id = $1", userID)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheUserTime)
	return
}
