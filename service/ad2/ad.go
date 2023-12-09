package ServiceAD2

import (
	"fmt"
	AnalysisAny2 "gitee.com/weeekj/weeekj_core/v5/analysis/any2"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"time"
)

// GetAD 获取指定广告
func GetAD(orgID int64, mark string) (data FieldsAD) {
	cacheMark := getADCacheMark(orgID, mark)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, org_id, mark, data FROM service_ad2 WHERE org_id = $1 AND mark = $2", orgID, mark)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 43200)
	return
}

// PutAD 投放广告
func PutAD(orgID int64, mark string) (data FieldsAD) {
	data = GetAD(orgID, mark)
	if data.ID < 1 {
		if orgID > 0 {
			data = GetAD(0, mark)
		}
	}
	if data.ID < 1 {
		return
	}
	AnalysisAny2.AppendData("add", "service_ad2_ad_put", time.Time{}, orgID, 0, data.ID, 0, 0, 1)
	return
}

// ClickAD 点击广告
func ClickAD(orgID int64, mark string, key int) {
	data := GetAD(orgID, mark)
	if data.ID < 1 {
		if orgID > 0 {
			data = GetAD(0, mark)
		}
	}
	if data.ID < 1 {
		return
	}
	isFind := false
	for k, _ := range data.Data {
		if k == key {
			isFind = true
			break
		}
	}
	if !isFind {
		return
	}
	AnalysisAny2.AppendData("add", "service_ad2_ad_click", time.Time{}, orgID, 0, data.ID, int64(key), 0, 1)
	return
}

// ArgsSetAD 设置广告参数
type ArgsSetAD struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//分区标识码
	// 同一个组织下唯一，具体行为交给mode识别和前端处理
	Mark string `db:"mark" json:"mark" check:"mark"`
	//结构
	Data FieldsADChildList `json:"data"`
}

// SetAD 设置广告
func SetAD(args *ArgsSetAD) (err error) {
	data := GetAD(args.OrgID, args.Mark)
	if data.ID > 0 && args.OrgID == data.OrgID {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE service_ad2 SET update_at = NOW(), data = :data WHERE id = :id", map[string]interface{}{
			"data": args.Data,
			"id":   data.ID,
		})
		if err != nil {
			return
		}
	} else {
		_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO service_ad2 (org_id, mark, data) VALUES (:org_id,:mark,:data)", map[string]interface{}{
			"org_id": args.OrgID,
			"mark":   args.Mark,
			"data":   args.Data,
		})
		if err != nil {
			return
		}
	}
	deleteADCache(args.OrgID, args.Mark)
	return
}

// 缓冲
func getADCacheMark(orgID int64, mark string) string {
	return fmt.Sprint("service:ad2:mark:", orgID, ".", mark)
}

func deleteADCache(orgID int64, mark string) {
	Router2SystemConfig.MainCache.DeleteMark(getADCacheMark(orgID, mark))
}
