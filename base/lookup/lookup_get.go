package BaseLookup

import (
	CoreCache "github.com/fotomxq/weeekj_core/v5/core/cache"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

type ArgsGetLookupList struct {
	//分页参数
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//是否为系统预设
	IsSys bool `db:"is_sys" json:"isSys" check:"bool"`
	//领域ID
	DomainID int64 `db:"domain_id" json:"domainID" check:"id" empty:"true"`
	//管理单元
	UnitID int64 `db:"unit_id" json:"unitID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

func GetLookupList(args *ArgsGetLookupList) (dataList []FieldsLookup, dataCount int64, err error) {
	dataCount, err = lookupDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "create_at", "update_at", "delete_at"}).SetPages(args.Pages).SelectList("((is_sys = $1) AND (domain_id = $2 OR $2 < 0) AND (unit_id = $3 OR $3 < 0) AND ((delete_at < to_timestamp(1000000) AND $4 = false) OR (delete_at >= to_timestamp(1000000) AND $4 = true)) AND (name LIKE $5 OR $5 = ''))", args.IsSys, args.DomainID, args.UnitID, args.IsRemove, "%"+args.Search+"%").ResultAndCount(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		vData := getLookupID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	return
}

// GetLookupAll 获取所有数据
func GetLookupAll(domainID int64, unitID int64) (dataList []FieldsLookup, err error) {
	err = lookupDB.DB.GetPostgresql().Select(&dataList, "SELECT id FROM base_lookup_child WHERE domain_id = $1 AND unit_id = $2 AND delete_at < to_timestamp(1000000)", domainID, unitID)
	if err != nil {
		return
	}
	if len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		vData := getLookupID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	return
}

func GetLookupID(id int64) (data FieldsLookup) {
	data = getLookupID(id)
	if data.ID < 1 {
		return
	}
	return
}

func GetLookupCode(code string) (data FieldsLookup) {
	cacheMark := getLookupCodeCacheMark(code)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		data = getLookupID(data.ID)
		if data.ID > 0 {
			return
		}
	}
	_ = lookupDB.Get().SetFieldsOne([]string{"id"}).NeedLimit().AppendWhere("code = $1", code).Result(&data)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, CoreCache.CacheTime1Hour)
	data = getLookupID(data.ID)
	if data.ID < 1 {
		data = FieldsLookup{}
		return
	}
	return
}

func getLookupID(id int64) (data FieldsLookup) {
	cacheMark := getLookupCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = lookupDB.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "delete_at", "is_sys", "domain_id", "unit_id", "code", "name"}).NeedLimit().GetByID(id).Result(&data)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, CoreCache.CacheTime1Hour)
	return
}
