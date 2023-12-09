package BaseFileSys2

import (
	"fmt"
	CoreCache "gitee.com/weeekj/weeekj_core/v5/core/cache"
	CoreSQL2 "gitee.com/weeekj/weeekj_core/v5/core/sql2"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

type ArgsGetClaimList struct {
	//分页参数
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//组织ID
	OrgID int64 `json:"orgID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `json:"userID" check:"id" empty:"true"`
	//文件结构体
	FileID int64 `json:"fileID" check:"id" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

func GetClaimList(args *ArgsGetClaimList) (dataList []FieldsFileClaim, dataCount int64, err error) {
	dataCount, err = claimDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "create_at", "update_at", "visit_count"}).SetPages(args.Pages).SelectList("(org_id = $1 OR $1 < 0) AND (user_id = $2 OR $2 < 0) AND (file_id = $3 OR $3 < 0) AND (des LIKE $4 OR $4 = '')", args.OrgID, args.UserID, args.FileID, "%"+args.Search+"%").ResultAndCount(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		vData := getClaimByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	return
}

func GetClaim(id int64) (data FieldsFileClaim) {
	data = getClaimByID(id)
	return
}

func getClaimByID(id int64) (data FieldsFileClaim) {
	cacheMark := getClaimCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = claimDB.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "update_hash", "org_id", "user_id", "is_public", "file_id", "expire_at", "visit_last_at", "visit_count", "des", "infos"}).GetByID(id).Result(&data)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, CoreCache.CacheTime1Hour)
	return
}

func getClaimCacheMark(id int64) string {
	return fmt.Sprint("base:filesys2.claim:id.", id)
}

func deleteClaimCache(id int64) {
	cacheMark := getClaimCacheMark(id)
	Router2SystemConfig.MainCache.DeleteMark(cacheMark)
}
