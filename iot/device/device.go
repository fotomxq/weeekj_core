package IOTDevice

import (
	"fmt"
	BaseExpireTip "github.com/fotomxq/weeekj_core/v5/base/expire_tip"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"

	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	"github.com/lib/pq"
	"gopkg.in/errgo.v2/fmt/errors"
)

// ArgsGetDeviceList 获取设备列表参数
type ArgsGetDeviceList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//设备分组
	GroupID int64 `db:"group_id" json:"groupID" check:"id" empty:"true"`
	//组织分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//组织标签ID组
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetDeviceList 获取设备列表
func GetDeviceList(args *ArgsGetDeviceList) (dataList []FieldsDevice, dataCount int64, err error) {
	maps := map[string]interface{}{}
	where := ""
	var rawList []FieldsDevice
	if args.OrgID > 0 {
		if args.IsRemove {
			where = "d.delete_at > to_timestamp(1000000)"
		} else {
			where = "d.delete_at < to_timestamp(1000000)"
		}
		maps["org_id"] = args.OrgID
		if args.GroupID > 0 {
			where = where + " AND d.group_id = :group_id"
			maps["group_id"] = args.GroupID
		}
		if args.Search != "" {
			where = where + " AND (d.name ILIKE '%' || :search || '%' OR d.des ILIKE '%' || :search || '%')"
			maps["search"] = args.Search
		}
		maps["sort_id"] = args.SortID
		if len(args.Tags) > 0 {
			where = where + " AND o.tags @> :tags"
			maps["tags"] = args.Tags
		}
		//where = where + " AND o.org_id = :org_id AND d.id = o.device_id"
		if args.Pages.Sort != "" {
			args.Pages.Sort = "d." + args.Pages.Sort
		}
		dataCount, err = CoreSQL.GetListPageAndCount(
			Router2SystemConfig.MainDB.DB,
			&rawList,
			"iot_core_device as d",
			"d.id",
			"SELECT d.id as id FROM iot_core_device as d INNER JOIN iot_core_operate as o ON d.id = o.device_id AND o.delete_at < to_timestamp(1000000) AND (:sort_id < 1 OR o.sort_id = :sort_id) WHERE o.org_id = :org_id AND "+where,
			where,
			maps,
			&args.Pages,
			[]string{"d.id", "d.create_at", "d.update_at", "d.delete_at", "d.last_at"},
		)
		if err != nil {
			return
		}
		/**
				var wherePage string
				wherePage, maps = CoreSQLPages.GetMaps(&args.Pages, maps)
				err = CoreSQL.GetList(
					Router2SystemConfig.MainDB.DB,
					&dataList,
					fmt.Sprint("SELECT d.id as id, d.create_at as create_at, d.update_at as update_at, d.delete_at as delete_at, d.status as status, d.is_online as is_online, d.last_at as last_at, d.name as name, d.des as des, d.cover_files as cover_files, d.des_files as des_files, d.group_id as group_id, d.code as code, d.address as address, d.params as params FROM iot_core_device as d INNER JOIN iot_core_operate as o ON d.id = o.device_id AND o.delete_at < to_timestamp(1000000) WHERE o.org_id = :org_id AND ", where, " GROUP BY d.id ", wherePage),
					maps,
				)
				dataCount, err = CoreSQL.GetAllCountMap(
					Router2SystemConfig.MainDB.DB,
					"iot_core_device as d",
					"d.id",
					fmt.Sprint("SELECT d.id as id, d.create_at as create_at, d.update_at as update_at, d.delete_at as delete_at, d.status as status, d.is_online as is_online, d.last_at as last_at, d.name as name, d.des as des, d.cover_files as cover_files, d.des_files as des_files, d.group_id as group_id, d.code as code, d.address as address, d.params as params FROM iot_core_device as d INNER JOIN iot_core_operate as o ON d.id = o.device_id AND o.delete_at < to_timestamp(1000000) WHERE o.org_id = :org_id AND ", where, " GROUP BY d.id"),
					maps,
				)
		**/
	} else {
		if args.IsRemove {
			where = "delete_at > to_timestamp(1000000)"
		} else {
			where = "delete_at < to_timestamp(1000000)"
		}
		if args.GroupID > 0 {
			where = where + " AND group_id = :group_id"
			maps["group_id"] = args.GroupID
		}
		if args.Search != "" {
			where = where + " AND (name ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%')"
			maps["search"] = args.Search
		}
		dataCount, err = CoreSQL.GetListPageAndCount(
			Router2SystemConfig.MainDB.DB,
			&rawList,
			"iot_core_device",
			"id",
			"SELECT id FROM iot_core_device WHERE "+where,
			where,
			maps,
			&args.Pages,
			[]string{"id", "create_at", "update_at", "delete_at", "last_at"},
		)
		if err != nil {
			return
		}
	}
	for _, v := range rawList {
		vData := getDeviceByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// ArgsGetDeviceByID 获取设备ID参数
type ArgsGetDeviceByID struct {
	//ID
	ID int64 `json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetDeviceByID 获取设备ID
func GetDeviceByID(args *ArgsGetDeviceByID) (data FieldsDevice, err error) {
	data = getDeviceByID(args.ID)
	if data.ID < 1 || CoreSQL.CheckTimeHaveData(data.DeleteAt) {
		err = errors.New("no data")
		return
	}
	operateData := getOperateByDeviceID(data.ID)
	if operateData.ID < 1 || CoreSQL.CheckTimeHaveData(operateData.DeleteAt) {
		err = errors.New("no data")
		return
	}
	if operateData.OrgID != args.OrgID {
		err = errors.New("no permission by org id")
		return
	}
	isFind := false
	for _, v := range operateData.Permissions {
		if v == "read" || v == "all" {
			isFind = true
		}
	}
	if !isFind {
		err = errors.New("no permission by operate")
		return
	}
	return
}

// ArgsGetDeviceKey 获取设备的key参数
type ArgsGetDeviceKey struct {
	//ID
	ID int64 `json:"id" check:"id"`
}

// GetDeviceKey 获取设备的key
func GetDeviceKey(args *ArgsGetDeviceKey) (key string, err error) {
	data := getDeviceByID(args.ID)
	if data.ID < 1 || CoreSQL.CheckTimeHaveData(data.DeleteAt) {
		err = errors.New("no data")
		return
	}
	key = data.Key
	return
}

// ArgsGetDeviceByCode 通过分组mark和设备code查询设备参数
type ArgsGetDeviceByCode struct {
	//设备分组
	GroupMark string `db:"group_mark" json:"groupMark" check:"mark"`
	//设备编号
	// 同一个分组下，必须唯一
	Code string `db:"code" json:"code" check:"mark"`
}

// GetDeviceByCode 通过分组mark和设备code查询设备
func GetDeviceByCode(args *ArgsGetDeviceByCode) (data FieldsDevice, err error) {
	groupData := getGroupByMark(args.GroupMark)
	if groupData.ID < 1 || CoreSQL.CheckTimeHaveData(groupData.DeleteAt) {
		err = errors.New("no data")
		return
	}
	data = getDeviceByCode(groupData.ID, args.Code)
	if data.ID < 1 || CoreSQL.CheckTimeHaveData(data.DeleteAt) {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsCheckDeviceCode 检查设备ID和code是否对应参数
type ArgsCheckDeviceCode struct {
	//所属设备
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id"`
	//设备分组
	GroupMark string `db:"group_mark" json:"groupMark" check:"mark"`
	//设备编号
	// 同一个分组下，必须唯一
	Code string `db:"code" json:"code" check:"mark"`
}

// CheckDeviceCode 检查设备ID和code是否对应
func CheckDeviceCode(args *ArgsCheckDeviceCode) (err error) {
	var data FieldsDevice
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT d.id as id FROM iot_core_device as d INNER JOIN iot_core_group as g ON d.group_id = g.id AND g.delete_at < to_timestamp(1000000) WHERE d.id = $1 AND g.mark = $2 AND d.code = $3 AND d.delete_at < to_timestamp(1000000) LIMIT 1", args.DeviceID, args.GroupMark, args.Code)
	if err != nil || data.ID < 1 {
		err = errors.New(fmt.Sprint("device not exist, ", err))
		return
	}
	return
}

// ArgsCheckDeviceKey 握手设备处理机制参数
type ArgsCheckDeviceKey struct {
	//设备分组
	GroupMark string `db:"group_mark" json:"groupMark" check:"mark"`
	//设备编号
	// 同一个分组下，必须唯一
	Code string `db:"code" json:"code" check:"mark"`
	//时间戳
	NowTime int64 `db:"now_time" json:"nowTime"`
	//随机码
	Rand string `db:"rand" json:"rand"`
	//key计算结果
	// key+code+时间戳+随机码
	Key string `db:"key" json:"key"`
}

// CheckDeviceKey 握手设备处理机制参数
func CheckDeviceKey(args *ArgsCheckDeviceKey) (err error) {
	_, err = CheckDeviceKeyAndDeviceID(args)
	return
}

func CheckDeviceKeyAndDeviceID(args *ArgsCheckDeviceKey) (deviceID int64, err error) {
	groupData := getGroupByMark(args.GroupMark)
	if groupData.ID < 1 || CoreSQL.CheckTimeHaveData(groupData.DeleteAt) {
		err = errors.New("group no data")
		return
	}
	deviceData := getDeviceByCode(groupData.ID, args.Code)
	if deviceData.ID < 1 || CoreSQL.CheckTimeHaveData(deviceData.DeleteAt) {
		err = errors.New(fmt.Sprint("not find device, group mark: ", args.GroupMark, ", code: ", args.Code))
		return
	}
	var key string
	key, err = CoreFilter.GetSha1ByString(fmt.Sprint(groupData.Mark, deviceData.Code, deviceData.Key, args.NowTime, args.Rand))
	if err != nil {
		err = errors.New("get sha1 by keys, " + err.Error())
		return
	}
	if key != args.Key {
		err = errors.New(fmt.Sprint("key error, post key: sum: ", fmt.Sprint(groupData.Mark, deviceData.Code, deviceData.Key, args.NowTime, args.Rand), ", key: ", args.Key, ", device key: ", key))
		return
	}
	deviceID = deviceData.ID
	return
}

// ArgsGetDeviceMore 获取一组设备参数
type ArgsGetDeviceMore struct {
	//ID列
	IDs pq.Int64Array `json:"ids"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove"`
}

// GetDeviceMore 获取一组设备
func GetDeviceMore(args *ArgsGetDeviceMore) (dataList []FieldsDevice, err error) {
	for _, v := range args.IDs {
		vData := getDeviceByID(v)
		if vData.ID < 1 {
			continue
		}
		if !args.HaveRemove && CoreSQL.CheckTimeHaveData(vData.DeleteAt) {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

func GetDeviceMoreMap(args *ArgsGetDeviceMore) (data map[int64]string, err error) {
	data = map[int64]string{}
	for _, v := range args.IDs {
		vData := getDeviceByID(v)
		if vData.ID < 1 {
			continue
		}
		if !args.HaveRemove && CoreSQL.CheckTimeHaveData(vData.DeleteAt) {
			continue
		}
		data[vData.ID] = vData.Name
	}
	return
}

// ArgsGetDeviceGroup 获取设备及设备组参数
type ArgsGetDeviceGroup struct {
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id"`
}

// GetDeviceGroupData 获取设备及设备组
type GetDeviceGroupData struct {
	//设备ID
	ID int64 `db:"id" json:"id"`
	//设备编码
	Code string `db:"code" json:"code"`
	//分组ID
	GroupID int64 `db:"group_id" json:"groupID"`
	//分组标识码
	GroupMark string `db:"group_mark" json:"groupMark"`
}

// GetDeviceGroup 获取设备及设备组
func GetDeviceGroup(args *ArgsGetDeviceGroup) (data GetDeviceGroupData, err error) {
	deviceData := getDeviceByID(args.DeviceID)
	if deviceData.ID < 1 || CoreSQL.CheckTimeHaveData(deviceData.DeleteAt) {
		err = errors.New("no data")
		return
	}
	data = GetDeviceGroupData{
		ID:        deviceData.ID,
		Code:      deviceData.Code,
		GroupID:   deviceData.GroupID,
		GroupMark: "",
	}
	if deviceData.GroupID > 0 {
		groupData := getGroupByID(deviceData.GroupID)
		data.GroupMark = groupData.Mark
	}
	return
}

// ArgsCreateDevice 创建新的设备参数
type ArgsCreateDevice struct {
	//状态
	// 0 public 公共可用 / 1 private 私有 / 2 ban 停用
	Status int `db:"status" json:"status" check:"intThan0" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"1000" empty:"true"`
	//封面列
	// 第一张作为封面
	CoverFiles pq.Int64Array `db:"cover_files" json:"coverFiles" check:"ids" empty:"true"`
	//描述信息
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//设备分组
	GroupID int64 `db:"group_id" json:"groupID" check:"id"`
	//设备编号
	// 同一个分组下，必须唯一
	Code string `db:"code" json:"code" check:"mark"`
	//连接密钥
	// 设备连接使用的唯一密钥
	// 设备需使用该key+code+时间戳+随机码混合计算，作为握手的识别码
	Key string `db:"key" json:"key"`
	//注册地
	// 如果设置将优先使用设备注册地，而不是管辖注册地
	Address string `db:"address" json:"address" check:"address" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateDevice 创建新的设备
func CreateDevice(args *ArgsCreateDevice) (data FieldsDevice, err error) {
	if args.GroupID > 0 {
		var groupData FieldsGroup
		groupData = getGroupByID(args.GroupID)
		if groupData.ID < 1 || CoreSQL.CheckTimeHaveData(groupData.DeleteAt) {
			err = errors.New("group not exist")
			return
		}
	}
	data = getDeviceByCode(args.GroupID, args.Code)
	if data.ID > 0 && !CoreSQL.CheckTimeHaveData(data.DeleteAt) {
		err = errors.New("code is replace")
		return
	}
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "iot_core_device", "INSERT INTO iot_core_device (status, is_online, last_at, name, des, cover_files, des_files, group_id, code, key, address, params) VALUES (:status,false,to_timestamp(0),:name,:des,:cover_files,:des_files,:group_id,:code,:key,:address,:params)", args, &data)
	if err != nil {
		return
	}
	return
}

// ArgsUpdateDevice 修改设备信息参数
type ArgsUpdateDevice struct {
	//ID
	ID int64 `json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//状态
	// 0 public 公共可用 / 1 private 私有 / 2 ban 停用
	Status int `db:"status" json:"status" check:"intThan0" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"1000" empty:"true"`
	//封面列
	// 第一张作为封面
	CoverFiles pq.Int64Array `db:"cover_files" json:"coverFiles" check:"ids" empty:"true"`
	//描述信息
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//设备分组
	GroupID int64 `db:"group_id" json:"groupID" check:"id"`
	//设备编号
	// 同一个分组下，必须唯一
	Code string `db:"code" json:"code" check:"mark"`
	//连接密钥
	// 设备连接使用的唯一密钥
	// 设备需使用该key+code+时间戳+随机码混合计算，作为握手的识别码
	Key string `db:"key" json:"key"`
	//注册地
	// 如果设置将优先使用设备注册地，而不是管辖注册地
	Address string `db:"address" json:"address" check:"address" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateDevice 修改设备信息
func UpdateDevice(args *ArgsUpdateDevice) (err error) {
	if args.GroupID > 0 {
		var groupData FieldsGroup
		groupData, err = GetGroupByID(&ArgsGetGroupByID{
			ID: args.GroupID,
		})
		if err != nil || groupData.ID < 1 {
			err = errors.New("group not exist")
			return
		}
	}
	var data FieldsDevice
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM iot_core_device WHERE group_id = $1 AND code = $2 AND id != $3 AND delete_at < to_timestamp(1000000)", args.GroupID, args.Code, args.ID)
	if err == nil && data.ID > 0 {
		err = errors.New("code is replace")
		return
	}
	data = getDeviceByID(args.ID)
	if data.ID < 1 || CoreSQL.CheckTimeHaveData(data.DeleteAt) {
		err = errors.New("no data")
		return
	}
	if args.Key == "" {
		args.Key = data.Key
	}
	if args.OrgID > 0 {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE iot_core_device as d INNER JOIN iot_core_operate as o ON o.device_id = d.id SET d.update_at = NOW(), d.status = :status, d.name = :name, d.des = :des, d.cover_files = :cover_files, d.des_files = :des_files, d.group_id = :group_id, d.code = :code, d.key = :key, d.address = :address, d.params = :params WHERE d.id = :id AND o.org_id = :org_id AND ('write' = ANY(o.permissions) OR 'all' = ANY(o.permissions)) AND d.delete_at < to_timestamp(1000000) AND o.delete_at < to_timestamp(1000000)", args)
		if err != nil {
			return
		}
	} else {
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE iot_core_device SET update_at = NOW(), status = :status, name = :name, des = :des, cover_files = :cover_files, des_files = :des_files, group_id = :group_id, code = :code, key = :key, address = :address, params = :params WHERE id = :id", args)
		if err != nil {
			return
		}
	}
	deleteDeviceCache(args.ID)
	return
}

// ArgsUpdateDeviceOnline 更新设备的在线状态
type ArgsUpdateDeviceOnline struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//状态
	IsOnline bool `db:"is_online" json:"isOnline" check:"bool"`
}

// UpdateDeviceOnline 更新设备的在线状态
func UpdateDeviceOnline(args *ArgsUpdateDeviceOnline) (err error) {
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE iot_core_device SET last_at = NOW(), is_online = :is_online WHERE id = :id", args)
	if err != nil {
		return
	}
	deleteDeviceCache(args.ID)
	deviceData := getDeviceByID(args.ID)
	if deviceData.ID > 0 {
		if deviceData.GroupID > 0 {
			groupData := getGroupByID(deviceData.GroupID)
			//计算设备组掉线时间长度
			groupExpireAt := CoreFilter.GetNowTimeCarbon().AddSeconds(int(groupData.ExpireTime))
			//发送通知
			BaseExpireTip.AppendTipNoErr(&BaseExpireTip.ArgsAppendTip{
				OrgID:      0,
				UserID:     0,
				SystemMark: "iot_device_group_online",
				BindID:     deviceData.ID,
				Hash:       "",
				ExpireAt:   groupExpireAt.Time,
			})
		}
	}
	return
}

// ArgsUpdateDeviceInfos 全量更新数据集合参数
type ArgsUpdateDeviceInfos struct {
	//ID
	ID int64 `json:"id" check:"id"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateDeviceInfos 全量更新数据集合
func UpdateDeviceInfos(args *ArgsUpdateDeviceInfos) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE iot_core_device SET update_at = NOW(), params = :params WHERE id = :id", args)
	if err != nil {
		return
	}
	deleteDeviceCache(args.ID)
	return
}

// ArgsUpdateDeviceInfo 分量更新数据集合参数
type ArgsUpdateDeviceInfo struct {
	//ID
	ID int64 `json:"id" check:"id"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateDeviceInfo 分量更新数据集合
func UpdateDeviceInfo(args *ArgsUpdateDeviceInfo) (err error) {
	data := getDeviceByID(args.ID)
	if data.ID < 1 || CoreSQL.CheckTimeHaveData(data.DeleteAt) {
		err = errors.New("no data")
		return
	}
	for _, v := range args.Params {
		isFind := false
		for k2, v2 := range data.Params {
			if v.Mark == v2.Mark {
				isFind = true
				data.Params[k2].Val = v.Val
				break
			}
		}
		if !isFind {
			data.Params = append(data.Params, v)
		}
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE iot_core_device SET update_at = NOW(), params = :params WHERE id = :id", map[string]interface{}{
		"id":     args.ID,
		"params": data.Params,
	})
	if err != nil {
		return
	}
	deleteDeviceCache(args.ID)
	return
}

// ArgsDeleteDevice 删除设备参数
type ArgsDeleteDevice struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteDevice 删除设备
func DeleteDevice(args *ArgsDeleteDevice) (err error) {
	if args.OrgID > 0 {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE iot_core_device as d INNER JOIN iot_core_operate as o ON o.device_id = d.id SET d.delete_at = NOW() WHERE d.id = :id AND o.org_id = :org_id", args)
		if err != nil {
			return
		}
	} else {
		_, err = CoreSQL.DeleteOneSoft(Router2SystemConfig.MainDB.DB, "iot_core_device", "id", map[string]interface{}{
			"id": args.ID,
		})
		if err != nil {
			return
		}
	}
	deleteDeviceCache(args.ID)
	return
}

// 获取指定ID
func getDeviceByID(id int64) (data FieldsDevice) {
	cacheMark := getDeviceCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, status, is_online, last_at, name, des, cover_files, des_files, group_id, code, key, address, params FROM iot_core_device WHERE id = $1", id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 1800)
	return
}

func getDeviceByCode(groupID int64, code string) (data FieldsDevice) {
	cacheMark := getDeviceCodeCacheMark(groupID, code)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, status, is_online, last_at, name, des, cover_files, des_files, group_id, code, key, address, params FROM iot_core_device WHERE group_id = $1 AND code = $2 AND delete_at < to_timestamp(1000000) ORDER BY id DESC LIMIT 1", groupID, code)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 1800)
	return
}

// 缓冲
func getDeviceCacheMark(id int64) string {
	return fmt.Sprint("iot:device:device:id:", id)
}

func getDeviceCodeCacheMark(groupID int64, code string) string {
	return fmt.Sprint("iot:device:device:code:", groupID, ".", code)
}

func deleteDeviceCache(id int64) {
	data := getDeviceByID(id)
	Router2SystemConfig.MainCache.DeleteMark(getDeviceCacheMark(id))
	if data.ID > 0 {
		Router2SystemConfig.MainCache.DeleteMark(getDeviceCodeCacheMark(data.GroupID, data.Code))
	}
}
