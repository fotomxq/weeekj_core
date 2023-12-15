package BaseFileSys

import (
	"fmt"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// 根据hash获取文件
func getFileByHash(hash string) (data FieldsFileType, err error) {
	//获取缓冲
	cacheMark := fmt.Sprint(getFileCacheMark(0), ":hash:", hash)
	if err = Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	//获取数据
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, update_hash, create_ip, create_info, file_size, file_type, file_hash, file_src, from_info, infos FROM core_file WHERE file_hash = $1", hash)
	//写入缓冲
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheTime)
	//反馈
	return
}
