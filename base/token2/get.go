package BaseToken2

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsGetList 获取列表参数
type ArgsGetList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//用户ID
	UserID int64 `json:"userID" check:"id" empty:"true"`
	//组织ID
	OrgID int64 `json:"orgID" check:"id" empty:"true"`
	//组织成员ID
	OrgBindID int64 `json:"orgBindID" check:"id" empty:"true"`
	//设备ID
	DeviceID int64 `json:"deviceID" check:"id" empty:"true"`
	//登录渠道
	LoginFrom string `json:"loginFrom" check:"mark" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetList 获取列表
func GetList(args *ArgsGetList) (dataList []FieldsToken, dataCount int64, err error) {
	maps := map[string]interface{}{}
	where := ""
	if args.UserID > -1 {
		where = where + "user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.OrgID > -1 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.OrgBindID > -1 {
		where = where + "org_bind_id = :org_bind_id"
		maps["org_bind_id"] = args.OrgBindID
	}
	if args.DeviceID > -1 {
		where = where + "device_id = :device_id"
		maps["device_id"] = args.DeviceID
	}
	if args.LoginFrom != "" {
		where = where + "login_from = :login_from"
		maps["login_from"] = args.LoginFrom
	}
	if args.Search != "" {
		where = where + " AND (ip ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	if where == "" {
		where = "true"
	}
	var rawList []FieldsToken
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		"core_token2",
		"id",
		"SELECT id FROM core_token2 WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "expire_at", "ip"},
	)
	if err != nil {
		return
	}
	//覆盖数据
	for _, v := range rawList {
		vData := getByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// GetByID 获取指定token
func GetByID(id int64) (data FieldsToken) {
	data = getByID(id)
	return
}

// GetByFrom 根据来源获取token
func GetByFrom(userID, deviceID int64, loginFrom string) (data FieldsToken) {
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM core_token2 WHERE user_id = $1 AND device_id = $2 AND ($3 = '' OR login_from = $3) ORDER BY id DESC LIMIT 1", userID, deviceID, loginFrom)
	if data.ID < 1 {
		return
	}
	data = getByID(data.ID)
	if data.ID < 1 {
		return
	}
	return
}

// getByID 获取token
func getByID(id int64) (data FieldsToken) {
	cacheMark := getTokenCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, expire_at, key, user_id, org_id, org_bind_id, device_id, login_from, ip, is_remember FROM core_token2 WHERE id = $1", id)
	if err != nil {
		return
	}
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 1800)
	return
}
