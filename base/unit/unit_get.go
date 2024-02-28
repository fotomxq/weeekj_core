package BaseUnit

import (
	CoreCache "github.com/fotomxq/weeekj_core/v5/core/cache"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

type ArgsGetUnitList struct {
	//分页参数
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

func GetUnitList(args *ArgsGetUnitList) (dataList []FieldsUnit, dataCount int64, err error) {
	dataCount, err = unitDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "create_at", "update_at", "delete_at"}).SetPages(args.Pages).SelectList("((delete_at < to_timestamp(1000000) AND $1 = false) OR (delete_at >= to_timestamp(1000000) AND $1 = true)) AND (name LIKE $2 OR $2 = '')", args.IsRemove, "%"+args.Search+"%").ResultAndCount(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		vData := GetUnitByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	return
}

func GetUnitByID(id int64) (data FieldsUnit) {
	data = getUnitByID(id)
	if data.ID < 1 {
		return
	}
	return
}

func GetUnitNameByID(id int64) (name string) {
	data := GetUnitByID(id)
	if data.ID < 1 {
		return
	}
	return data.Name
}

func GetUnitByCode(code string) (data FieldsUnit) {
	cacheMark := getUnitCodeCacheMark(code)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		data = getUnitByID(data.ID)
		if data.ID > 0 {
			return
		}
	}
	_ = unitDB.Get().SetFieldsOne([]string{"id"}).NeedLimit().AppendWhere("code = $1", code).Result(&data)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, CoreCache.CacheTime1Hour)
	data = getUnitByID(data.ID)
	if data.ID < 1 {
		data = FieldsUnit{}
		return
	}
	return
}

func GetUnitNoDeleteByCode(code string) (data FieldsUnit) {
	data = GetUnitByCode(code)
	if data.ID < 1 {
		return
	}
	if CoreFilter.CheckHaveTime(data.DeleteAt) {
		data = FieldsUnit{}
		return
	}
	return
}

func getUnitByID(id int64) (data FieldsUnit) {
	cacheMark := getUnitCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = unitDB.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "delete_at", "code", "name"}).NeedLimit().GetByID(id).Result(&data)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, CoreCache.CacheTime1Hour)
	return
}
