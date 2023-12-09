package BaseExpireTip

import (
	"errors"
	"fmt"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// 获取数据集合
func getID(id int64) (data FieldsTip, err error) {
	if err = Router2SystemConfig.MainCache.GetStruct(getCacheMark(id), &data); err == nil && data.ID > 0 {
		return
	}
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, org_id, user_id, system_mark, bind_id, hash, expire_at FROM core_expire_tip WHERE id = $1", id)
	if err != nil {
		return
	}
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	Router2SystemConfig.MainCache.SetStruct(getCacheMark(id), data, 3650)
	return
}

// 获取缓冲名称
func getCacheMark(id int64) string {
	return fmt.Sprint("base:expire:at:id:", id)
}
