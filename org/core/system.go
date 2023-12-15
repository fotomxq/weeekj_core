package OrgCoreCore

import (
	"errors"
	"fmt"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
	"strings"
)

// ArgsGetSystemList 获取系统关联列表参数
type ArgsGetSystemList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// 可选，-1忽略
	OrgID int64 `json:"orgID" check:"id" empty:"true"`
	//来源系统类型
	// eg: wxx 微信小程序
	SystemMark string `db:"system_mark" json:"systemMark" check:"mark" empty:"true"`
	//唯一的标识码
	// 例如小程序的AppID
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool" empty:"true"`
}

// GetSystemList 获取系统关联列表
func GetSystemList(args *ArgsGetSystemList) (dataList []FieldsSystem, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.SystemMark != "" {
		where = where + " AND system_mark = :system_mark"
		maps["system_mark"] = args.SystemMark
		if args.Mark != "" {
			where = where + " AND mark = :mark"
			maps["mark"] = args.Mark
		}
	}
	var rawList []FieldsSystem
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		"org_core_system",
		"id",
		"SELECT id FROM org_core_system WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getSystemByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	for k, v := range dataList {
		for k2, v2 := range v.Params {
			if strings.Count(v2.Mark, "key") > 0 && v2.Val != "" {
				dataList[k].Params[k2].Val = "***"
			}
		}
	}
	return
}

// ArgsFilterOrgIDsBySystem 根据一组组织ID和标识码系统，淘汰不匹配的数据参数
type ArgsFilterOrgIDsBySystem struct {
	//组织ID
	// 可选，-1忽略
	OrgIDs pq.Int64Array `json:"orgIDs" check:"ids"`
	//来源系统类型
	// eg: wxx 微信小程序
	SystemMark string `db:"system_mark" json:"systemMark" check:"mark"`
	//唯一的标识码
	// 例如小程序的AppID
	Mark string `db:"mark" json:"mark" check:"mark"`
}

// FilterOrgIDsBySystem 根据一组组织ID和标识码系统，淘汰不匹配的数据
func FilterOrgIDsBySystem(args *ArgsFilterOrgIDsBySystem) (ids []int64, err error) {
	type FieldsSystem struct {
		//组织ID
		OrgID int64 `db:"org_id" json:"orgID"`
	}
	var dataList []FieldsSystem
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT org_id FROM org_core_system WHERE delete_at < to_timestamp(1000000) AND system_mark = $1 AND mark = $2 AND org_id = ANY($3)", args.SystemMark, args.Mark, args.OrgIDs)
	if err == nil && len(dataList) > 0 {
		for _, v := range dataList {
			ids = append(ids, v.OrgID)
		}
		return
	}
	return
}

// GetSystemOrgIDBySystem 通过配置查询组织ID
func GetSystemOrgIDBySystem(system string, mark string) (orgID int64, err error) {
	err = Router2SystemConfig.MainDB.Get(&orgID, "SELECT org_id FROM org_core_system WHERE delete_at < to_timestamp(1000000) AND system_mark = $1 AND mark = $2", system, mark)
	if err != nil {
		return
	}
	if orgID < 1 {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsGetSystem 获取指定的组织指定标识码数据参数
type ArgsGetSystem struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//来源系统类型
	// eg: wxx 微信小程序
	SystemMark string `db:"system_mark" json:"systemMark" check:"mark"`
}

// GetSystem 获取指定的组织指定标识码数据
func GetSystem(args *ArgsGetSystem) (data FieldsSystem, err error) {
	data = getSystemByMark(args.OrgID, args.SystemMark)
	if data.ID < 1 || CoreSQL.CheckTimeHaveData(data.DeleteAt) {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsSetSystem 创建新的系统关联参数
type ArgsSetSystem struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//来源系统类型
	// eg: wxx 微信小程序
	SystemMark string `db:"system_mark" json:"systemMark" check:"mark"`
	//唯一的标识码
	// 例如小程序的AppID
	Mark string `db:"mark" json:"mark" check:"mark"`
	//附加参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// SetSystem 创建新的系统关联
func SetSystem(args *ArgsSetSystem) (data FieldsSystem, err error) {
	data, err = GetSystem(&ArgsGetSystem{
		OrgID:      args.OrgID,
		SystemMark: args.SystemMark,
	})
	if err == nil && data.ID > 0 {
		if args.SystemMark == "weixin_wxx" {
			var orgID int64
			orgID, err = GetSystemOrgIDBySystem(args.SystemMark, args.Mark)
			if orgID != args.OrgID {
				err = errors.New("app id have data")
				return
			}
		}
		for k, v := range args.Params {
			if v.Val == "***" {
				for _, v2 := range data.Params {
					if v2.Mark == v.Mark {
						args.Params[k].Val = v2.Val
					}
				}
			}
		}
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_core_system SET update_at = NOW(), mark = :mark, params = :params WHERE id = :id", map[string]interface{}{
			"id":     data.ID,
			"mark":   args.Mark,
			"params": args.Params,
		})
		if err != nil {
			return
		}
		data.Mark = args.Mark
		data.Params = args.Params
		deleteSystemCache(data.ID)
	} else {
		err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "org_core_system", "INSERT INTO org_core_system (org_id, system_mark, mark, params) VALUES (:org_id,:system_mark,:mark,:params)", args, &data)
		if err != nil {
			return
		}
	}
	return
}

// ArgsDeleteSystem 删除系统关联参数
type ArgsDeleteSystem struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 可选，验证项
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// DeleteSystem 删除系统关联
func DeleteSystem(args *ArgsDeleteSystem) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "org_core_system", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	deleteSystemCache(args.ID)
	return
}

// 获取ID
func getSystemByID(id int64) (data FieldsSystem) {
	cacheMark := getSystemCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, system_mark, mark, params FROM org_core_system WHERE id = $1", id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, systemCacheTime)
	return
}

// 获取组织下系统来源
func getSystemByMark(orgID int64, systemMark string) (data FieldsSystem) {
	cacheMark := getSystemCacheSystemMark(orgID, systemMark)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, system_mark, mark, params FROM org_core_system WHERE org_id = $1 AND system_mark = $2 AND delete_at < to_timestamp(1000000)", orgID, systemMark)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, systemCacheTime)
	return
}

// 缓冲
func getSystemCacheMark(id int64) string {
	return fmt.Sprint("org:core:system:id:", id)
}

func getSystemCacheSystemMark(orgID int64, systemMark string) string {
	return fmt.Sprint("org:core:system:org:", orgID, ".", systemMark)
}

func deleteSystemCache(id int64) {
	data := getSystemByID(id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.DeleteMark(getSystemCacheMark(data.ID))
	Router2SystemConfig.MainCache.DeleteMark(getSystemCacheSystemMark(data.OrgID, data.SystemMark))
}
