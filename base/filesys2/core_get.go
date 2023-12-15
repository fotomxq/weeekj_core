package BaseFileSys2

import (
	"fmt"
	CoreCache "github.com/fotomxq/weeekj_core/v5/core/cache"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetCoreList 获取文件列表参数
type ArgsGetCoreList struct {
	//分页
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//组织ID
	OrgID int64 `json:"orgID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `json:"userID" check:"id" empty:"true"`
	//存储方式
	// local 本地化单一服务器存储; qiniu 七牛云存储
	SaveSystem string `db:"save_system" json:"saveSystem" check:"mark" empty:"true"`
	//存储块
	SaveMark string `db:"save_mark" json:"saveMark" check:"mark" empty:"true"`
}

// GetCoreList 获取文件列表
func GetCoreList(args *ArgsGetCoreList) (dataList []FieldsFile, dataCount int64, err error) {
	dataCount, err = coreDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "create_at", "update_at", "file_size"}).SetPages(args.Pages).SelectList("(org_id = $1 OR $1 < 0) OR (user_id = $2 OR $2 < 0) OR (save_system = $3 OR $3 = '') OR (save_mark = $4 OR $4 = '')", args.OrgID, args.UserID, args.SaveSystem, args.SaveMark).ResultAndCount(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		vData := getCoreByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	return
}

// GetCore 获取文件数据
func GetCore(id int64) (data FieldsFile) {
	data = getCoreByID(id)
	return
}

// 获取缓冲
func getCoreCacheMark(id int64) string {
	return fmt.Sprint("base:filesys2:core:id.", id)
}

func getCoreByID(id int64) (data FieldsFile) {
	cacheMark := getCoreCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = coreDB.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "update_hash", "create_ip", "org_id", "user_id", "file_size", "file_type", "file_hash", "file_src", "save_system", "save_mark", "save_success", "infos"}).Result(&data)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, CoreCache.CacheTime1Hour)
	return
}

func deleteCoreCache(id int64) {
	cacheMark := getCoreCacheMark(id)
	Router2SystemConfig.MainCache.DeleteMark(cacheMark)
}
