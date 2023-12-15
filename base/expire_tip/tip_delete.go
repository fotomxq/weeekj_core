package BaseExpireTip

import (
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

// 清理数据
func deleteID(id int64) (err error) {
	waitExpire1HourLock.Lock()
	defer waitExpire1HourLock.Unlock()
	_, err = CoreSQL.DeleteOne(Router2SystemConfig.MainDB.DB, "core_expire_tip", "id", map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return
	}
	var newCacheList []FieldsTip
	for _, v := range waitExpire1HourList {
		if v.ID == id {
			continue
		}
		newCacheList = append(newCacheList, v)
	}
	waitExpire1HourList = newCacheList
	Router2SystemConfig.MainCache.DeleteMark(getCacheMark(id))
	return
}

// DataGetExpireData 解析数据包数据
type DataGetExpireData struct {
	OrgID    int64     `json:"orgID"`
	UserID   int64     `json:"userID"`
	ExpireAt time.Time `json:"expireAt"`
}

// GetExpireData 解析数据包
// 通知过期后，可使用此方法解析数据包
func GetExpireData(rawData []byte) (data DataGetExpireData, err error) {
	err = CoreNats.ReflectDataByte(rawData, &data)
	return
}
