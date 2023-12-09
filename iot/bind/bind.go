package IOTBind

import (
	"errors"
	"fmt"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsGetBindList 获取绑定列表参数
type ArgsGetBindList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id" empty:"true"`
	//附加模块
	FromInfo CoreSQLFrom.FieldsFrom `db:"from_info" json:"fromInfo"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
}

// GetBindList 获取绑定列表
func GetBindList(args *ArgsGetBindList) (dataList []FieldsBind, dataCount int64, err error) {
	maps := map[string]interface{}{}
	where := ""
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.DeviceID > -1 {
		where = where + " AND device_id = :device_id"
		maps["device_id"] = args.DeviceID
	}
	where, maps, err = args.FromInfo.GetListAnd("from_info", "from_info", where, maps)
	if err != nil {
		return
	}
	tableName := "iot_core_bind"
	var rawList []FieldsBind
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		"SELECT id FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getBindByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// ArgsGetBindFrom 获取来源的所有绑定设备参数
type ArgsGetBindFrom struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//附加模块
	FromInfo CoreSQLFrom.FieldsFrom `db:"from_info" json:"fromInfo"`
}

// GetBindFrom 获取来源的所有绑定设备
func GetBindFrom(args *ArgsGetBindFrom) (dataList []FieldsBind, err error) {
	rawList := getBindByFrom(args.OrgID, args.FromInfo)
	for _, v := range rawList {
		if v.ID < 1 || CoreSQL.CheckTimeHaveData(v.DeleteAt) {
			continue
		}
		dataList = append(dataList, v)
	}
	if len(dataList) < 1 {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsCheckBind 验证设备绑定关系是否存在参数
type ArgsCheckBind struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id"`
	//附加模块
	FromInfo CoreSQLFrom.FieldsFrom `db:"from_info" json:"fromInfo"`
}

// CheckBind 验证设备绑定关系是否存在
func CheckBind(args *ArgsCheckBind) (err error) {
	dataList := getBindByFrom(args.OrgID, args.FromInfo)
	if len(dataList) < 1 {
		err = errors.New("no data")
		return
	}
	for _, v := range dataList {
		if v.ID < 1 || CoreSQL.CheckTimeHaveData(v.DeleteAt) || v.DeviceID != args.DeviceID {
			continue
		}
		//找到匹配直接反馈成功
		return
	}
	err = errors.New("no data")
	return
}

// ArgsGetBindFromGroup 获取来源的所有绑定设备带设备分组过滤参数
type ArgsGetBindFromGroup struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//附加模块
	FromInfo CoreSQLFrom.FieldsFrom `db:"from_info" json:"fromInfo"`
	//设备分组标识码
	DeviceGroupMark string `db:"device_group_mark" json:"deviceGroupMark" check:"mark"`
}

// GetBindFromGroup 获取来源的所有绑定设备带设备分组过滤
func GetBindFromGroup(args *ArgsGetBindFromGroup) (dataList []FieldsBind, err error) {
	var fromData string
	fromData, err = args.FromInfo.GetRawNoName()
	if err != nil {
		return
	}
	var rawList []FieldsBind
	err = Router2SystemConfig.MainDB.Select(&rawList, "SELECT b.id as id FROM iot_core_bind as b INNER JOIN iot_core_device as d ON d.id = b.device_id INNER JOIN iot_core_group as g ON g.id = d.group_id WHERE b.org_id = $1 AND b.from_info @> $2 AND g.mark = $3 AND b.delete_at < to_timestamp(1000000)", args.OrgID, fromData, args.DeviceGroupMark)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getBindByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// ArgsGetBindDevice 获取设备的所有绑定关系参数
type ArgsGetBindDevice struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id"`
}

// GetBindDevice 获取设备的所有绑定关系
func GetBindDevice(orgID int64, deviceID int64) (data FieldsBind, err error) {
	data = getBindByOrg(orgID, deviceID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	return
}

func GetBindByDeviceID(deviceID int64) (data FieldsBind, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM iot_core_bind WHERE device_id = $1 AND delete_at < to_timestamp(1000000) ORDER BY id DESC LIMIT 1", deviceID)
	if err != nil || data.ID < 1 {
		err = errors.New("no data")
		return
	}
	data = getBindByID(data.ID)
	if data.ID < 1 || CoreSQL.CheckTimeHaveData(data.DeleteAt) {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsSetBind 设置绑定关系参数
type ArgsSetBind struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id"`
	//附加模块
	FromInfo CoreSQLFrom.FieldsFrom `db:"from_info" json:"fromInfo"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// SetBind 设置绑定关系
func SetBind(args *ArgsSetBind) (data FieldsBind, err error) {
	data = getBindByOrg(args.OrgID, args.DeviceID)
	if data.ID < 1 {
		err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "iot_core_bind", "INSERT INTO iot_core_bind (org_id, device_id, from_info, params) VALUES (:org_id,:device_id,:from_info,:params)", args, &data)
		if err != nil {
			return
		}
	} else {
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE iot_core_bind SET update_at = NOW(), from_info = :from_info, params = :params WHERE id = :id", map[string]interface{}{
			"id":        data.ID,
			"from_info": args.FromInfo,
			"params":    args.Params,
		})
		if err != nil {
			return
		}
		data.FromInfo = args.FromInfo
		data.Params = args.Params
		deleteBindCache(data.ID)
	}
	//清洗其他来源
	var fromData string
	fromData, err = args.FromInfo.GetRawNoName()
	if err != nil {
		err = nil
	} else {
		var otherFromList []FieldsBind
		err = Router2SystemConfig.MainDB.Select(&otherFromList, "SELECT id, create_at, update_at, delete_at, org_id, device_id, from_info, params FROM iot_core_bind WHERE org_id = $1 AND from_info @> $2 AND id != $3 AND delete_at < to_timestamp(1000000)", args.OrgID, fromData, data.ID)
		if err != nil {
			err = nil
		} else {
			for _, v := range otherFromList {
				_ = DeleteBind(&ArgsDeleteBind{
					ID:    v.ID,
					OrgID: v.OrgID,
				})
			}
		}
	}
	//反馈
	return
}

// ArgsDeleteBind 删除绑定参数
type ArgsDeleteBind struct {
	//ID
	ID int64 `json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteBind 删除绑定
func DeleteBind(args *ArgsDeleteBind) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "iot_core_bind", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	deleteBindCache(args.ID)
	return
}

// 获取指定ID
func getBindByID(id int64) (data FieldsBind) {
	cacheMark := getBindCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, device_id, from_info, params FROM iot_core_bind WHERE id = $1", id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 1800)
	return
}

// getBindByOrg 仅用于查询组织和设备ID非删除类绑定关系
func getBindByOrg(orgID int64, deviceID int64) (data FieldsBind) {
	cacheMark := getBindOrgCacheMark(orgID, deviceID)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, device_id, from_info, params FROM iot_core_bind WHERE org_id = $1 AND device_id = $2 and delete_at < to_timestamp(1000000)", orgID, deviceID)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 1800)
	return
}

func getBindByFrom(orgID int64, fromInfo CoreSQLFrom.FieldsFrom) (dataList []FieldsBind) {
	cacheMark := getBindFromCacheMark(orgID, fromInfo)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &dataList); err == nil && len(dataList) > 0 {
		return
	}
	fromData, err := fromInfo.GetRawNoName()
	if err != nil {
		return
	}
	_ = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, create_at, update_at, delete_at, org_id, device_id, from_info, params FROM iot_core_bind WHERE org_id = $1 AND from_info @> $2", orgID, fromData)
	if len(dataList) < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, dataList, 1800)
	return
}

// 缓冲
func getBindCacheMark(id int64) string {
	return fmt.Sprint("iot:bind:id:", id)
}

func getBindOrgCacheMark(orgID int64, deviceID int64) string {
	return fmt.Sprint("iot:bind:org:v2:", orgID, ".", deviceID)
}

func getBindFromCacheMark(orgID int64, fromInfo CoreSQLFrom.FieldsFrom) string {
	return fmt.Sprint("iot:bind:from:", orgID, ".", fromInfo.System, ".", fromInfo.Mark, ".", fromInfo.ID)
}

func deleteBindCache(id int64) {
	data := getBindByID(id)
	Router2SystemConfig.MainCache.DeleteMark(getBindCacheMark(id))
	if data.ID > 0 {
		Router2SystemConfig.MainCache.DeleteMark(getBindOrgCacheMark(data.OrgID, data.DeviceID))
		Router2SystemConfig.MainCache.DeleteMark(getBindFromCacheMark(data.OrgID, data.FromInfo))
	}
}
