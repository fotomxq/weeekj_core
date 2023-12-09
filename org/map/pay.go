package OrgMap

import (
	"errors"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// CheckMapPay 检查地图支付状态
func CheckMapPay(mapID int64) bool {
	var data FieldsMapPay
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, finish_at, org_id, user_id, map_id, pay_id, count FROM org_map_pay WHERE map_id = $1 ORDER BY id DESC LIMIT 1", mapID)
	if err != nil {
		return false
	}
	if CoreSQL.CheckTimeHaveData(data.FinishAt) {
		return true
	}
	return false
}

// CreateMapPay 创建新的支付请求
func CreateMapPay(mapID int64, payID int64, count int64) (err error) {
	mapData := getMapByID(mapID)
	if mapData.ID < 1 || CoreSQL.CheckTimeHaveData(mapData.DeleteAt) {
		err = errors.New("no data")
		return
	}
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO org_map_pay(org_id, user_id, map_id, pay_id, count) VALUES(:org_id, :user_id, :map_id, :pay_id, :count)", map[string]interface{}{
		"org_id":  mapData.OrgID,
		"user_id": mapData.UserID,
		"map_id":  mapData.ID,
		"pay_id":  payID,
		"count":   count,
	})
	if err != nil {
		return
	}
	return
}

// 更新支付状态
func updateMapPay(payID int64) (err error) {
	//获取支付ID对应的记录
	var payLog FieldsMapPay
	err = Router2SystemConfig.MainDB.Get(&payLog, "SELECT id, create_at, finish_at, org_id, user_id, map_id, pay_id, count FROM org_map_pay WHERE pay_id = $1 ORDER BY id DESC LIMIT 1", payID)
	if err != nil {
		return
	}
	//更新支付状态
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE org_map_pay SET finish_at = NOW() WHERE pay_id = :pay_id", map[string]interface{}{
		"pay_id": payID,
	})
	if err != nil {
		return
	}
	//获取地图数据
	mapData := getMapByID(payLog.MapID)
	//修改地图的限制次数
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE org_map SET ad_count_limit = ad_count_limit + :count WHERE id = :id", map[string]interface{}{
		"id":    mapData.ID,
		"count": payLog.Count,
	})
	if err != nil {
		return
	}
	return
}
