package BaseLookup

import (
	CoreCache "github.com/fotomxq/weeekj_core/v5/core/cache"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

type ArgsGetDomainList struct {
	//分页参数
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

func GetDomainList(args *ArgsGetDomainList) (dataList []FieldsDomain, dataCount int64, err error) {
	dataCount, err = domainDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "create_at", "update_at", "delete_at", "name"}).SetPages(args.Pages).SelectList("((delete_at < to_timestamp(1000000) AND $1 = false) OR (delete_at >= to_timestamp(1000000) AND $1 = true)) AND (name LIKE $2 OR $2 = '')", args.IsRemove, "%"+args.Search+"%").ResultAndCount(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		vData := getDomainID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	return
}

func GetDomainID(id int64) (data FieldsDomain) {
	data = getDomainID(id)
	if data.ID < 1 {
		return
	}
	return
}

func GetDomainNameByID(id int64) (name string) {
	data := getDomainID(id)
	if data.ID < 1 {
		return
	}
	return data.Name
}

func GetDomainByName(name string) (data FieldsDomain) {
	_ = domainDB.DB.GetPostgresql().DB.Get(&data, "SELECT id FROM base_lookup_domain WHERE name = $1 and delete_at < to_timestamp(1000000)", name)
	if data.ID < 1 {
		return
	}
	data = GetDomainID(data.ID)
	return
}

func getDomainID(id int64) (data FieldsDomain) {
	cacheMark := getDomainCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = domainDB.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "delete_at", "name"}).NeedLimit().GetByID(id).Result(&data)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, CoreCache.CacheTime1Hour)
	return
}
