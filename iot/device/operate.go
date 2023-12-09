package IOTDevice

import (
	"errors"
	"fmt"
	BaseExpireTip "gitee.com/weeekj/weeekj_core/v5/base/expire_tip"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"time"

	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	"github.com/lib/pq"
)

// ArgsGetOperateList 获取授权列表参数
type ArgsGetOperateList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetOperateList 获取授权列表参数
func GetOperateList(args *ArgsGetOperateList) (dataList []FieldsOperate, dataCount int64, err error) {
	maps := map[string]interface{}{}
	where := ""
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > 0 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.Search != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "(address ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	if where == "" {
		where = "true"
	}
	tableName := "iot_core_operate"
	var rawList []FieldsOperate
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		"SELECT id FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at", "expire_at"},
	)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getOperateByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// ArgsGetOperateAndDevice 获取设备ID分组和授权组织ID参数
type ArgsGetOperateAndDevice struct {
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id"`
}

// GetOperateAndDeviceData 获取设备ID分组和授权组织ID数据
type GetOperateAndDeviceData struct {
	//设备ID
	ID int64 `db:"id" json:"id"`
	//设备分组
	GroupID int64 `db:"group_id" json:"groupID"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
}

// GetOperateAndDevice 获取设备ID分组和授权组织ID
func GetOperateAndDevice(args *ArgsGetOperateAndDevice) (data []GetOperateAndDeviceData, err error) {
	err = Router2SystemConfig.MainDB.Select(&data, "SELECT d.id as id, d.group_id as group_id, o.org_id as org_id FROM iot_core_operate as o INNER JOIN iot_core_device as d ON o.device_id = d.id WHERE d.delete_at < to_timestamp(1000000) AND o.device_id = $1 AND o.delete_at < to_timestamp(1000000)", args.DeviceID)
	return
}

// ArgsCheckOperate 检查组织是否具备设备参数
type ArgsCheckOperate struct {
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// CheckOperate 检查组织是否具备设备
// 同时反馈该操作权限
func CheckOperate(args *ArgsCheckOperate) (data FieldsOperate, err error) {
	data = getOperateByDeviceID(args.DeviceID)
	if data.ID < 1 || CoreSQL.CheckTimeHaveData(data.DeleteAt) || !CoreFilter.EqID2(args.OrgID, data.OrgID) {
		err = errors.New(fmt.Sprint("device id: ", args.DeviceID, ", org id: ", args.OrgID, ", err: ", err))
		return
	}
	if CoreSQL.CheckTimeHaveData(data.ExpireAt) && !CoreSQL.CheckTimeThanNow(data.ExpireAt) {
		err = errors.New("no data")
		return
	}
	return
}

func CheckOperateNoData(args *ArgsCheckOperate) (err error) {
	data := getOperateByDeviceID(args.DeviceID)
	if data.ID < 1 || CoreSQL.CheckTimeHaveData(data.DeleteAt) || !CoreFilter.EqID2(args.OrgID, data.OrgID) {
		err = errors.New(fmt.Sprint("device id: ", args.DeviceID, ", org id: ", args.OrgID, ", err: ", err))
		return
	}
	if CoreSQL.CheckTimeHaveData(data.ExpireAt) && !CoreSQL.CheckTimeThanNow(data.ExpireAt) {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsCheckOperates 检查组织是否具备一组设备参数
type ArgsCheckOperates struct {
	//设备IDs
	DeviceIDs pq.Int64Array `db:"device_ids" json:"deviceIDs" check:"ids"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// CheckOperates 检查组织是否具备一组设备
func CheckOperates(args *ArgsCheckOperates) (err error) {
	for _, v := range args.DeviceIDs {
		vData := getOperateByDeviceID(v)
		if vData.ID < 1 || CoreSQL.CheckTimeHaveData(vData.DeleteAt) {
			continue
		}
		if vData.OrgID != args.OrgID {
			err = errors.New(fmt.Sprint("device id: ", vData.DeviceID, " no this org id: ", args.OrgID))
			return
		}
	}
	return
}

// ArgsGetOperateByDeviceID 获取设备的控制人参数
type ArgsGetOperateByDeviceID struct {
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id"`
}

// GetOperateByDeviceID 获取设备的控制人
func GetOperateByDeviceID(args *ArgsGetOperateByDeviceID) (dataList []FieldsOperate, err error) {
	var rawList []FieldsOperate
	err = Router2SystemConfig.MainDB.Select(&rawList, "SELECT id FROM iot_core_operate WHERE device_id = $1 AND delete_at < to_timestamp(1000000)", args.DeviceID)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getOperateByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// ArgsGetOperateActionByDeviceID 获取设备支持的动作参数
type ArgsGetOperateActionByDeviceID struct {
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// GetOperateActionByDeviceID 获取设备支持的动作
func GetOperateActionByDeviceID(args *ArgsGetOperateActionByDeviceID) (actionList []FieldsAction, err error) {
	if args.OrgID > 0 {
		//获取组织授权的动作集合
		// 通过组织授权，结合设备组结构交叉获取数据
		var operateData FieldsOperate
		err = Router2SystemConfig.MainDB.Get(&operateData, "SELECT a.id as id, a.action as action, a.permissions as permissions FROM iot_core_operate as a WHERE a.device_id = $1 AND a.org_id = $2 AND delete_at < to_timestamp(1000000)", args.DeviceID, args.OrgID)
		if err != nil || operateData.ID < 1 {
			err = errors.New(fmt.Sprint("operate not exist, ", err))
			return
		}
		//配对检查权限
		isFind := false
		for _, v := range operateData.Permissions {
			if v == "all" || v == "mission" {
				isFind = true
				break
			}
		}
		if !isFind {
			err = errors.New("org no permission")
			return
		}
		var newActionList []FieldsAction
		for _, v := range actionList {
			for _, v2 := range operateData.Action {
				if v.ID == v2 {
					newActionList = append(newActionList, v)
					break
				}
			}
		}
		if len(newActionList) < 1 {
			err = errors.New("no action")
			return
		}
		return
	} else {
		//获取设备动作集合
		// 直接获取设备设备组的数据包
		err = Router2SystemConfig.MainDB.Select(&actionList, "SELECT a.id as id, a.create_at as create_at, a.update_at as update_at, a.delete_at as delete_at, a.mark as mark, a.name as name, a.des as des, a.expire_time as expire_time, a.connect_type as connect_type, a.configs as configs FROM iot_core_action as a INNER JOIN iot_core_device as d ON d.id = $1 AND d.delete_at < to_timestamp(1000000) INNER JOIN iot_core_group as g ON g.id = d.group_id AND g.delete_at < to_timestamp(1000000) WHERE a.delete_at < to_timestamp(1000000) AND a.id = ANY(g.action)", args.DeviceID)
		if err != nil || len(actionList) < 1 {
			err = errors.New(fmt.Sprint("action empty, ", err))
			return
		}
		return
	}
}

// ArgsSetOperate 设置授权参数
type ArgsSetOperate struct {
	//授权过期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt" check:"isoTime"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//权限类型
	// all 全部权限
	// read 允许查看设备信息
	// write 允许编辑设备信息
	// mission 任务下达权限
	// operate 修改授权关系
	// associated 关联设备
	Permissions pq.StringArray `db:"permissions" json:"permissions" check:"marks"`
	//允许执行的动作
	// 将根据设备组的动作查询，如果存在则允许，否则将禁止执行该类动作
	Action pq.Int64Array `db:"action" json:"action" check:"ids" empty:"true"`
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id"`
	//注册地
	Address string `db:"address" json:"address" check:"address" empty:"true"`
	//组织分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//组织标签ID组
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// SetOperate 设置授权
func SetOperate(args *ArgsSetOperate) (err error) {
	var data FieldsOperate
	//设备不能被其他组织绑定
	operateData := getOperateByDeviceID(args.DeviceID)
	if operateData.ID > 0 && operateData.OrgID != args.OrgID && !CoreSQL.CheckTimeHaveData(operateData.DeleteAt) {
		err = errors.New("device have bind and not this org")
		return
	}
	//检查动作是否超出范围
	// 获取设备所属的设备组
	var groupData FieldsGroup
	err = Router2SystemConfig.MainDB.Get(&groupData, "SELECT g.id as id, g.action as action FROM iot_core_group as g INNER JOIN iot_core_device as d ON d.group_id = g.id WHERE d.id = $1", args.DeviceID)
	if err != nil || groupData.ID < 1 {
		err = errors.New("not find group")
		return
	}
	for _, v := range args.Action {
		isFind := false
		for _, v2 := range groupData.Action {
			if v == v2 {
				isFind = true
				break
			}
		}
		if !isFind {
			err = errors.New("action not support")
			return
		}
	}
	//查询组织和设备是否存在？
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM iot_core_operate WHERE device_id = $1 AND org_id = $2 AND delete_at < to_timestamp(1000000)", args.DeviceID, args.OrgID)
	if err == nil && data.ID > 0 {
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE iot_core_operate SET update_at = NOW(), expire_at = :expire_at, permissions = :permissions, action = :action, address = :address, sort_id = :sort_id, tags = :tags, params = :params WHERE org_id = :org_id AND device_id = :device_id", args)
		if err != nil {
			return
		}
		deleteOperateCache(args.DeviceID)
	} else {
		err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "iot_core_operate", "INSERT INTO iot_core_operate (expire_at, org_id, permissions, action, device_id, address, sort_id, tags, params) VALUES (:expire_at,:org_id,:permissions,:action,:device_id,:address,:sort_id,:tags,:params)", args, &data)
		if err != nil {
			return
		}
	}
	//过期数据删除处理
	if CoreSQL.CheckTimeHaveData(data.ExpireAt) && CoreSQL.CheckTimeThanNow(data.ExpireAt) {
		BaseExpireTip.AppendTipNoErr(&BaseExpireTip.ArgsAppendTip{
			OrgID:      0,
			UserID:     0,
			SystemMark: "iot_device_operate",
			BindID:     data.ID,
			Hash:       "",
			ExpireAt:   data.ExpireAt,
		})
	}
	return
}

// ArgsUpdateOperateByOrg 组织修改绑定信息参数
type ArgsUpdateOperateByOrg struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id"`
	//注册地
	Address string `db:"address" json:"address" check:"address" empty:"true"`
	//组织分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//组织标签ID组
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
}

// UpdateOperateByOrg 组织修改绑定信息参数
func UpdateOperateByOrg(args *ArgsUpdateOperateByOrg) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE iot_core_operate SET update_at = NOW(), address = :address, sort_id = :sort_id, tags = :tags WHERE org_id = :org_id AND device_id = :device_id", args)
	if err != nil {
		return
	}
	deleteOperateCache(args.DeviceID)
	return
}

// ArgsDeleteOperate 删除授权参数
type ArgsDeleteOperate struct {
	//设备IDs
	DeviceIDs pq.Int64Array `db:"device_ids" json:"deviceIDs" check:"ids"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// DeleteOperate 删除授权
func DeleteOperate(args *ArgsDeleteOperate) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "iot_core_operate", "device_id = ANY(:device_ids) AND org_id = :org_id", args)
	if err != nil {
		return
	}
	for _, v := range args.DeviceIDs {
		deleteOperateCache(v)
	}
	return
}

func deleteOperateByID(id int64) (err error) {
	_, err = CoreSQL.DeleteOneSoft(Router2SystemConfig.MainDB.DB, "iot_core_operate", "id", map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return
	}
	var data FieldsOperate
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, device_id FROM iot_core_operate WHERE id = $1", id)
	if err != nil {
		return
	}
	deleteOperateCache(data.DeviceID)
	return
}

// 获取指定ID
func getOperateByID(id int64) (data FieldsOperate) {
	cacheMark := getOperateIDCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, expire_at, org_id, permissions, action, device_id, address, sort_id, tags, params FROM iot_core_operate WHERE id = $1", id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 1800)
	return
}

func getOperateByDeviceID(deviceID int64) (data FieldsOperate) {
	cacheMark := getOperateCacheMark(deviceID)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, expire_at, org_id, permissions, action, device_id, address, sort_id, tags, params FROM iot_core_operate WHERE device_id = $1 AND delete_at < to_timestamp(1000000) ORDER BY id DESC LIMIT 1", deviceID)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 1800)
	return
}

// 缓冲
func getOperateCacheMark(deviceID int64) string {
	return fmt.Sprint("iot:device:operate:device:", deviceID)
}

func getOperateIDCacheMark(id int64) string {
	return fmt.Sprint("iot:device:operate:id:", id)
}

func deleteOperateCache(deviceID int64) {
	data := getOperateByDeviceID(deviceID)
	Router2SystemConfig.MainCache.DeleteMark(getOperateCacheMark(deviceID))
	if data.ID > 0 {
		Router2SystemConfig.MainCache.DeleteMark(getOperateIDCacheMark(data.ID))
	}
}
