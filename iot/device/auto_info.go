package IOTDevice

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsOpenAutoInfo 触发指定的info参数
type ArgsOpenAutoInfo struct {
	//组织ID
	// 如果留空，方法内将检索该设备控制方
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//触发设备ID
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id"`
	//触发条件
	// 扩展参数mark
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
	//值
	Val string `db:"val" json:"val"`
}

// OpenAutoInfo 触发指定的info
func OpenAutoInfo(args *ArgsOpenAutoInfo) (reportInfoList []FieldsAutoInfo, needReport bool, err error) {
	//检查组织ID
	var orgList []int64
	if args.OrgID < 1 {
		var operateList []FieldsOperate
		operateList, err = GetOperateByDeviceID(&ArgsGetOperateByDeviceID{
			DeviceID: args.DeviceID,
		})
		if err != nil {
			err = errors.New(fmt.Sprint("no any operate, ", err))
			return
		}
		for _, v := range operateList {
			if v.OrgID < 1 {
				continue
			}
			orgList = append(orgList, v.OrgID)
		}
	} else {
		orgList = append(orgList, args.OrgID)
	}
	//查询所有和该mark关联的组织下的处理方法
	for _, vOrgID := range orgList {
		var infoList []FieldsAutoInfo
		err = Router2SystemConfig.MainDB.Select(&infoList, "SELECT id FROM iot_core_auto_info WHERE org_id = $1 AND device_id = $2 AND (mark = $3 OR template_id > 0)", vOrgID, args.DeviceID, args.Mark)
		if err != nil || len(infoList) < 1 {
			//CoreLog.Error("get auto info data, no info data, org id: ", vOrgID, ", device id: ", args.DeviceID, ", mark: ", args.Mark, ", err: ", err)
			err = nil
			continue
		}
		//遍历集合，依次触发
		for _, v := range infoList {
			v = getAutoInfoByID(v.ID)
			if v.ID < 1 {
				continue
			}
			//如果存在模版，获取模版数据
			if v.TemplateID > 0 {
				var templateData FieldsAutoInfoTemplate
				templateData, err = GetAutoInfoTemplate(&ArgsGetAutoInfoTemplate{
					ID: v.TemplateID,
				})
				if err != nil || templateData.ID < 1 {
					err = errors.New(fmt.Sprint("template lost, ", err, ", id: ", v.TemplateID))
					return
				}
				if templateData.Mark != args.Mark {
					continue
				}
				v.WaitTime = templateData.WaitTime
				v.Mark = templateData.Mark
				v.Eq = templateData.Eq
				v.Val = templateData.Val
				v.SendAction = templateData.SendAction
				if len(v.ParamsData) < 1 {
					v.ParamsData = templateData.ParamsData
				}
			}
			//检查间隔时间，如果在此时间触发过，则跳过
			var logData FieldsAutoLog
			err = Router2SystemConfig.MainDB.Get(&logData, "SELECT id FROM iot_core_auto_log WHERE create_at > $1", CoreFilter.GetNowTimeCarbon().SubSeconds(int(v.WaitTime)))
			if err == nil && logData.ID > 0 {
				err = nil
				continue
			}
			//引发触发器
			switch v.Eq {
			case 0:
				if args.Val == v.Val {
					//检查触发器的基本条件
					// 如果不存在目标设备ID，则直接记录日志
					if v.ReportDeviceID < 1 {
						//记录日志
						err = createAutoLog(&argsCreateAutoLog{
							OrgID:    v.OrgID,
							DeviceID: v.DeviceID,
							InfoID:   v.ID,
							Mark:     v.Mark,
							Eq:       v.Eq,
							EqVal:    v.Val,
							Val:      fmt.Sprint(args.Val),
						})
					} else {
						err = sendAutoToDevice(&argsSendAutoToDevice{}, &v)
					}
					if err != nil {
						err = errors.New(fmt.Sprint("send auto to device, ", err))
						return
					}
					reportInfoList = append(reportInfoList, v)
					needReport = true
				}
			case 1:
				var v1, v2 int64
				v1, err = CoreFilter.GetInt64ByString(args.Val)
				if err != nil {
					err = nil
					continue
				}
				v2, err = CoreFilter.GetInt64ByString(v.Val)
				if err != nil {
					err = nil
					continue
				}
				if v1 < v2 {
					//检查触发器的基本条件
					// 如果不存在目标设备ID，则直接记录日志
					if v.ReportDeviceID < 1 {
						//记录日志
						err = createAutoLog(&argsCreateAutoLog{
							OrgID:    v.OrgID,
							DeviceID: v.DeviceID,
							InfoID:   v.ID,
							Mark:     v.Mark,
							Eq:       v.Eq,
							EqVal:    v.Val,
							Val:      fmt.Sprint(args.Val),
						})
					} else {
						err = sendAutoToDevice(&argsSendAutoToDevice{}, &v)
					}
					if err != nil {
						err = errors.New(fmt.Sprint("send auto to device, ", err))
						return
					}
					reportInfoList = append(reportInfoList, v)
					needReport = true
				}
			case 2:
				var v1, v2 int64
				v1, err = CoreFilter.GetInt64ByString(args.Val)
				if err != nil {
					err = nil
					continue
				}
				v2, err = CoreFilter.GetInt64ByString(v.Val)
				if err != nil {
					err = nil
					continue
				}
				if v1 > v2 {
					//检查触发器的基本条件
					// 如果不存在目标设备ID，则直接记录日志
					if v.ReportDeviceID < 1 {
						//记录日志
						err = createAutoLog(&argsCreateAutoLog{
							OrgID:    v.OrgID,
							DeviceID: v.DeviceID,
							InfoID:   v.ID,
							Mark:     v.Mark,
							Eq:       v.Eq,
							EqVal:    v.Val,
							Val:      fmt.Sprint(args.Val),
						})
					} else {
						err = sendAutoToDevice(&argsSendAutoToDevice{}, &v)
					}
					if err != nil {
						err = errors.New(fmt.Sprint("send auto to device, ", err))
						return
					}
					reportInfoList = append(reportInfoList, v)
					needReport = true
				}
			case 3:
				if args.Val != v.Val {
					//检查触发器的基本条件
					// 如果不存在目标设备ID，则直接记录日志
					if v.ReportDeviceID < 1 {
						//记录日志
						err = createAutoLog(&argsCreateAutoLog{
							OrgID:    v.OrgID,
							DeviceID: v.DeviceID,
							InfoID:   v.ID,
							Mark:     v.Mark,
							Eq:       v.Eq,
							EqVal:    v.Val,
							Val:      fmt.Sprint(args.Val),
						})
					} else {
						err = sendAutoToDevice(&argsSendAutoToDevice{}, &v)
					}
					if err != nil {
						err = errors.New(fmt.Sprint("send auto to device, ", err))
						return
					}
					reportInfoList = append(reportInfoList, v)
					needReport = true
				}
			default:
				continue
			}
		}
	}
	//反馈结果
	return
}

// 给设备推送信息
type argsSendAutoToDevice struct {
	//实际发生值
	Val string `json:"val"`
}

func sendAutoToDevice(args *argsSendAutoToDevice, infoData *FieldsAutoInfo) (err error) {
	//记录日志
	err = createAutoLog(&argsCreateAutoLog{
		OrgID:    infoData.OrgID,
		DeviceID: infoData.DeviceID,
		InfoID:   infoData.ID,
		Mark:     infoData.Mark,
		Eq:       infoData.Eq,
		EqVal:    infoData.Val,
		Val:      fmt.Sprint(args.Val),
	})
	return
}

// ArgsGetAutoInfoList 获取关联列表参数
type ArgsGetAutoInfoList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//触发条件
	// 扩展参数mark
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
	//任务动作
	SendAction string `db:"send_action" json:"sendAction" check:"mark" empty:"true"`
}

// GetAutoInfoList 获取关联列表
func GetAutoInfoList(args *ArgsGetAutoInfoList) (dataList []FieldsAutoInfo, dataCount int64, err error) {
	maps := map[string]interface{}{}
	where := ""
	if args.OrgID > 0 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.Mark != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "mark = :mark"
		maps["mark"] = args.Mark
	}
	if args.SendAction != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "send_action = :send_action"
		maps["send_action"] = args.SendAction
	}
	if where == "" {
		where = "true"
	}
	tableName := "iot_core_auto_info"
	var rawList []FieldsAutoInfo
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		"SELECT id FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getAutoInfoByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// ArgsGetAutoInfo 获取指定关联ID参数
type ArgsGetAutoInfo struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetAutoInfo 获取指定关联ID
func GetAutoInfo(args *ArgsGetAutoInfo) (data FieldsAutoInfo, err error) {
	data = getAutoInfoByID(args.ID)
	if data.ID < 1 || CoreFilter.EqID2(args.OrgID, data.OrgID) {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsCreateAutoInfo 创建新的关联参数
type ArgsCreateAutoInfo struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//触发设备ID
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id"`
	//采用模版
	// 如果存在模版ID，自定义触发条件将无效
	TemplateID int64 `db:"template_id" json:"templateID" check:"id" empty:"true"`
	//触发条件
	// 扩展参数mark
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
	// 等式
	// 0 等于; 1 小于; 2 大于; 3 不等于
	Eq int `db:"eq" json:"eq" check:"intThan0" empty:"true"`
	//值
	Val string `db:"val" json:"val"`
	//冷却时间
	WaitTime int64 `db:"wait_time" json:"waitTime" check:"int64Than0"`
	//反馈设备ID
	// 如果没有指定反馈设备ID，将该记录自动归档到log表中，方便外部模块查询
	ReportDeviceID int64 `db:"report_device_id" json:"reportDeviceID" check:"id" empty:"true"`
	//发送任务指令
	// 留空则发送触发条件的数据包
	SendAction string `db:"send_action" json:"sendAction" check:"mark" empty:"true"`
	//发送参数
	ParamsData []byte `db:"params_data" json:"paramsData"`
}

// CreateAutoInfo 创建新的关联
func CreateAutoInfo(args *ArgsCreateAutoInfo) (data FieldsAutoInfo, err error) {
	if args.TemplateID > 0 {
		_, err = GetAutoInfoTemplate(&ArgsGetAutoInfoTemplate{
			ID: args.TemplateID,
		})
		if err != nil {
			err = errors.New(fmt.Sprint("get info template data, ", err))
			return
		}
	}
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "iot_core_auto_info", "INSERT INTO iot_core_auto_info (org_id, device_id, template_id, mark, eq, val, wait_time, report_device_id, send_action, params_data) VALUES (:org_id,:device_id,:template_id,:mark,:eq,:val,:wait_time,:report_device_id,:send_action,:params_data)", args, &data)
	return
}

// ArgsUpdateAutoInfo 修改关联参数
type ArgsUpdateAutoInfo struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//触发条件
	// 扩展参数mark
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
	// 等式
	// 0 等于; 1 小于; 2 大于; 3 不等于
	Eq int `db:"eq" json:"eq" check:"intThan0" empty:"true"`
	//值
	Val string `db:"val" json:"val"`
	//冷却时间
	WaitTime int64 `db:"wait_time" json:"waitTime" check:"int64Than0"`
	//反馈设备ID
	// 如果没有指定反馈设备ID，将该记录自动归档到log表中，方便外部模块查询
	ReportDeviceID int64 `db:"report_device_id" json:"reportDeviceID" check:"id" empty:"true"`
	//发送任务指令
	// 留空则发送触发条件的数据包
	SendAction string `db:"send_action" json:"sendAction" check:"mark" empty:"true"`
	//发送参数
	ParamsData []byte `db:"params_data" json:"paramsData"`
}

// UpdateAutoInfo 修改关联
func UpdateAutoInfo(args *ArgsUpdateAutoInfo) (err error) {
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE iot_core_auto_info SET mark = :mark, eq = :eq, val = :val, wait_time = :wait_time, report_device_id = :report_device_id, send_action = :send_action, params_data = :params_data WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	deleteAutoInfoCache(args.ID)
	return
}

// ArgsDeleteAutoInfo 删除关联参数
type ArgsDeleteAutoInfo struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteAutoInfo 删除关联
func DeleteAutoInfo(args *ArgsDeleteAutoInfo) (err error) {
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "iot_core_auto_info", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	deleteAutoInfoCache(args.ID)
	return
}

// 获取指定ID
func getAutoInfoByID(id int64) (data FieldsAutoInfo) {
	cacheMark := getAutoInfoCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, org_id, device_id, template_id, mark, eq, val, wait_time, report_device_id, send_action, params_data FROM iot_core_auto_info WHERE id = $1", id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 1800)
	return
}

// 缓冲
func getAutoInfoCacheMark(id int64) string {
	return fmt.Sprint("iot:device:auto:template:id:", id)
}

func deleteAutoInfoCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getAutoInfoCacheMark(id))
}
