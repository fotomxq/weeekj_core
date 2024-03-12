package IOTQuickRecord

import (
	"errors"
	"fmt"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	IOTDevice "github.com/fotomxq/weeekj_core/v5/iot/device"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetList 获取列表参数
type ArgsGetList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetList 获取列表
func GetList(args *ArgsGetList) (dataList []FieldsRecord, dataCount int64, err error) {
	maps := map[string]interface{}{}
	where := ""
	if args.Search != "" {
		where = where + "(device_code ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	if where == "" {
		where = "true"
	}
	tableName := "iot_quick_record"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, device_code, device_id, params FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	return
}

// ArgsCreate 请求生成设备数据参数
type ArgsCreate struct {
	//设备标识码
	DeviceCode string `db:"device_code" json:"deviceCode"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// Create 请求生成设备数据
func Create(args *ArgsCreate) (data FieldsRecord, err error) {
	//检查是否存在记录
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, device_code, device_id, params FROM iot_quick_record WHERE device_code = $1", args.DeviceCode)
	if err == nil && data.ID > 0 {
		return
	}
	//创建数据
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "iot_quick_record", "INSERT INTO iot_quick_record (device_code, params) VALUES (:device_code,:params)", args, &data)
	return
}

// ArgsAudit 审核通过并生成设备参数
type ArgsAudit struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//设备ID
	// 匹配好的设备，保留直到设备领取数据
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id"`
}

// Audit 审核通过并生成设备参数
func Audit(args *ArgsAudit) (err error) {
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE iot_quick_record SET device_id = :device_id WHERE id = :id", args)
	return
}

// ArgsGetResult 拉取请求结果参数
type ArgsGetResult struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

type DataGetResult struct {
	//设备组标识码
	GroupMark string `db:"groupMark" json:"groupMark"`
	//设备编号
	// 同一个分组下，必须唯一
	Code string `db:"code" json:"code"`
	//连接密钥
	// 设备连接使用的唯一密钥
	// 设备需使用该key+code+时间戳+随机码混合计算，作为握手的识别码
	Key string `db:"key" json:"key"`
}

// GetResult 拉取请求结果
func GetResult(args *ArgsGetResult) (data DataGetResult, err error) {
	//检查是否存在记录
	var recordData FieldsRecord
	err = Router2SystemConfig.MainDB.Get(&recordData, "SELECT id, create_at, device_code, device_id, params FROM iot_quick_record WHERE id = $1 AND device_id > 0", args.ID)
	if err != nil || recordData.ID < 1 {
		err = errors.New(fmt.Sprint("no data, arg id: ", args.ID, ", err: ", err))
		return
	}
	//获取数据并反馈
	var deviceData IOTDevice.FieldsDevice
	deviceData, err = IOTDevice.GetDeviceByID(&IOTDevice.ArgsGetDeviceByID{
		ID:    recordData.DeviceID,
		OrgID: -1,
	})
	if err != nil {
		err = errors.New(fmt.Sprint("get device data error, ", err))
		return
	}
	var groupData IOTDevice.FieldsGroup
	groupData, err = IOTDevice.GetGroupByID(&IOTDevice.ArgsGetGroupByID{
		ID: deviceData.GroupID,
	})
	if err != nil {
		err = errors.New(fmt.Sprint("get group data error, ", err))
		return
	}
	var deviceKey string
	deviceKey, err = IOTDevice.GetDeviceKey(&IOTDevice.ArgsGetDeviceKey{
		ID: recordData.DeviceID,
	})
	if err != nil {
		err = errors.New(fmt.Sprint("get device key error, ", err))
		return
	}
	data.GroupMark = groupData.Mark
	data.Key = deviceKey
	data.Code = recordData.DeviceCode
	if err == nil {
		_, _ = CoreSQL.DeleteOne(Router2SystemConfig.MainDB.DB, "iot_quick_record", "id", map[string]interface{}{
			"id": recordData.ID,
		})
	}
	return
}
