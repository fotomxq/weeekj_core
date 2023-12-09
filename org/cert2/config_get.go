package OrgCert2

import (
	"errors"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLIDs "gitee.com/weeekj/weeekj_core/v5/core/sql/ids"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetConfigList 获取证件配置列表参数
type ArgsGetConfigList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//绑定来源
	// user 用户 / org 商户 / org_bind 商户成员 / finance_assets 财务资产 /
	BindFrom string `db:"bind_from" json:"bindFrom" check:"mark" empty:"true"`
	//标识码
	// 用于程序化识别处理机制
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
	//是否被删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetConfigList 获取证件配置列表
func GetConfigList(args *ArgsGetConfigList) (dataList []FieldsConfig, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.BindFrom != "" {
		where = where + " AND bind_from = :bind_from"
		maps["bind_from"] = args.BindFrom
	}
	if args.Mark != "" {
		where = where + " AND mark = :mark"
		maps["mark"] = args.Mark
	}
	if args.Search != "" {
		where = where + " AND (name ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	var rawList []FieldsConfig
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		"org_cert_config2",
		"id",
		"SELECT id FROM org_cert_config2 WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getConfigByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// ArgsGetConfigByID 获取指定配置ID参数
type ArgsGetConfigByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 可选
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetConfigByID 获取指定配置ID
func GetConfigByID(args *ArgsGetConfigByID) (data FieldsConfig, err error) {
	data = getConfigByID(args.ID)
	if data.ID < 1 || CoreSQL.CheckTimeHaveData(data.DeleteAt) || !CoreFilter.EqID2(args.OrgID, data.OrgID) {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsGetConfigByMark 获取指定配置Mark参数
type ArgsGetConfigByMark struct {
	//mark
	Mark string `db:"mark" json:"mark" check:"mark"`
	//组织ID
	// 可选
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetConfigByMark 获取指定配置Mark
func GetConfigByMark(args *ArgsGetConfigByMark) (data FieldsConfig, err error) {
	data = getConfigByMark(args.OrgID, args.Mark)
	if data.ID < 1 || CoreSQL.CheckTimeHaveData(data.DeleteAt) {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsGetConfigMore 获取一组配置参数
type ArgsGetConfigMore struct {
	//ID列
	IDs pq.Int64Array `json:"ids" check:"ids"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove" check:"bool"`
}

// GetConfigMore 获取一组配置
func GetConfigMore(args *ArgsGetConfigMore) (dataList []FieldsConfig, err error) {
	for _, v := range args.IDs {
		vData := getConfigByID(v)
		if vData.ID < 1 || (!args.HaveRemove && CoreSQL.CheckTimeHaveData(vData.DeleteAt)) {
			continue
		}
		dataList = append(dataList, vData)

	}
	return
}

// GetConfigName 获取配置名称
func GetConfigName(configID int64) string {
	data := getConfigByID(configID)
	if data.ID < 1 {
		return ""
	}
	return data.Name
}

// GetConfigFrom 获取绑定来源
func GetConfigFrom(configID int64) string {
	data := getConfigByID(configID)
	if data.ID < 1 {
		return ""
	}
	return data.BindFrom
}

// GetConfigMoreMap 获取一组配置名称组
func GetConfigMoreMap(args *ArgsGetConfigMore) (data map[int64]string, err error) {
	data, err = CoreSQLIDs.GetIDsNameAndDelete("org_cert_config2", args.IDs, args.HaveRemove)
	return
}

// 根据mark获取配置
func getConfigByMark(orgID int64, mark string) (data FieldsConfig) {
	if mark == "" {
		return
	}
	cacheMark := getConfigMarkCacheMark(orgID, mark)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, default_expire, org_id, bind_from, mark, name, des, cover_file_id, des_files, audit_type, currency, price, sn_len, tip_type, params FROM org_cert_config2 WHERE org_id = $1 AND mark = $2 AND delete_at < to_timestamp(1000000) LIMIT 1", orgID, mark)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 21600)
	return
}

// 获取指定配置ID
func getConfigByID(id int64) (data FieldsConfig) {
	cacheMark := getConfigCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, default_expire, org_id, bind_from, mark, name, des, cover_file_id, des_files, audit_type, currency, price, sn_len, tip_type, params FROM org_cert_config2 WHERE id = $1", id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 21600)
	return
}
