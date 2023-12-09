package OrderTake

import (
	"fmt"
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	CoreCache "gitee.com/weeekj/weeekj_core/v5/core/cache"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	OrgCore "gitee.com/weeekj/weeekj_core/v5/org/core"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// IsOpenTakeCode 检查是否打开了自提验证功能
func IsOpenTakeCode(orgID int64) (b bool) {
	//获取平台开关
	b = BaseConfig.GetDataBoolNoErr("ServiceOrderSelfTakeMustCode")
	if b {
		return
	}
	//检查组织开关
	b = OrgCore.Config.GetConfigValBoolNoErr(orgID, "OrderSelfTakeMustCode")
	//反馈
	return
}

// GetTakeCode 获取订单的自提代码
func GetTakeCode(orderID int64) (takeCode string) {
	//检查订单是否具备自提单信息
	data := GetTakeData(orderID)
	//如果没有则创建
	if data.ID < 1 {
		takeCode = CoreFilter.GetRandStr4(6)
		err := sqlTake.Insert().SetFields([]string{"order_id", "take_code"}).Add(map[string]interface{}{
			"order_id":  orderID,
			"take_code": takeCode,
		}).ExecAndCheckID()
		if err != nil {
			return
		}
		deleteTakeCache(orderID)
		return
	}
	//如果有则返回
	takeCode = data.TakeCode
	return
}

// CheckTakeCode 验证自提代码
func CheckTakeCode(orderID int64, takeCode string) (b bool) {
	data := GetTakeData(orderID)
	if data.ID < 1 {
		return
	}
	b = data.TakeCode == takeCode
	return
}

// GetTakeData 获取订单数据
func GetTakeData(orderID int64) (data FieldsTake) {
	cacheMark := getTakeCacheMark(orderID)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := sqlTake.Get().SetFieldsOne([]string{"id", "create_at", "order_id", "take_code"}).AppendWhere("order_id = $1", orderID).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, CoreCache.CacheTime1Hour)
	return
}

// 获取缓冲标识码
func getTakeCacheMark(orderID int64) (mark string) {
	return fmt.Sprint("service:order:take:order_id:", orderID)
}

// 删除缓冲
func deleteTakeCache(orderID int64) {
	Router2SystemConfig.MainCache.DeleteMark(getTakeCacheMark(orderID))
}
