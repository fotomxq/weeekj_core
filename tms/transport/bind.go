package TMSTransport

import (
	"errors"
	"fmt"
	ClassConfig "github.com/fotomxq/weeekj_core/v5/class/config"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLGPS "github.com/fotomxq/weeekj_core/v5/core/sql/gps"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	MapArea "github.com/fotomxq/weeekj_core/v5/map/area"
	OrgCore "github.com/fotomxq/weeekj_core/v5/org/core"
	OrgTime "github.com/fotomxq/weeekj_core/v5/org/time"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetBindList 获取绑定列表参数
type ArgsGetBindList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//配送人员
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//分区ID
	MapAreaID int64 `db:"map_area_id" json:"mapAreaID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
}

// GetBindList 获取绑定列表
func GetBindList(args *ArgsGetBindList) (dataList []FieldsBind, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.BindID > -1 {
		where = where + " AND bind_id = :bind_id"
		maps["bind_id"] = args.BindID
	}
	if args.MapAreaID > 0 {
		where = where + " AND map_area_id = :map_area_id"
		maps["map_area_id"] = args.MapAreaID
	}
	tableName := "tms_transport_bind"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, update_at, delete_at, org_id, bind_id, map_area_id, more_map_area_ids, km_30_day, level_30_day, time_30_day, count_30_day, count_finish_30_day, km_1_day, level_1_day, time_1_day, count_1_day, count_finish_1_day, un_finish_count, params FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at", "km_30_day", "level_30_day", "time_30_day", "count_30_day", "count_finish_30_day"},
	)
	return
}

// ArgsGetBind 获取指定ID参数
type ArgsGetBind struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetBind 获取指定ID
func GetBind(args *ArgsGetBind) (data FieldsBind, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, bind_id, map_area_id, more_map_area_ids, km_30_day, level_30_day, time_30_day, count_30_day, count_finish_30_day, km_1_day, level_1_day, time_1_day, count_1_day, count_finish_1_day, un_finish_count, params FROM tms_transport_bind WHERE id = $1 AND (org_id = $2 OR $2 < 1)", args.ID, args.OrgID)
	if err == nil && data.ID < 1 {
		err = errors.New("bind not exist")
	}
	if err != nil {
		err = errors.New(fmt.Sprint("get bind data, bind id: ", args.ID, ", ", err))
	}
	return
}

// getBindByBindID 通过成员ID获取配送员信息
func getBindByBindID(bindID int64) (data FieldsBind) {
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, bind_id, map_area_id, more_map_area_ids, km_30_day, level_30_day, time_30_day, count_30_day, count_finish_30_day, km_1_day, level_1_day, time_1_day, count_1_day, count_finish_1_day, un_finish_count, params FROM tms_transport_bind WHERE bind_id = $1", bindID)
	if err != nil || data.ID < 1 {
		return
	}
	return
}

// GetBindByBind 获取指定ID
func GetBindByBind(args *ArgsGetBind) (data FieldsBind, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, bind_id, map_area_id, more_map_area_ids, km_30_day, level_30_day, time_30_day, count_30_day, count_finish_30_day, km_1_day, level_1_day, time_1_day, count_1_day, count_finish_1_day, un_finish_count, params FROM tms_transport_bind WHERE bind_id = $1 AND (org_id = $2 OR $2 < 1)", args.ID, args.OrgID)
	if err == nil && data.ID < 1 {
		err = errors.New("bind not exist")
	}
	if err != nil {
		err = errors.New(fmt.Sprint("get bind data, bind id: ", args.ID, ", ", err))
	}
	return
}

// ArgsGetBindByBindID 通过组织成员ID获取绑定关系参数
type ArgsGetBindByBindID struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//组织成员ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id"`
}

// GetBindByBindID 通过组织成员ID获取绑定关系
func GetBindByBindID(args *ArgsGetBindByBindID) (data FieldsBind, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, bind_id, map_area_id, km_30_day, level_30_day, time_30_day, count_30_day, count_finish_30_day, km_1_day, level_1_day, time_1_day, count_1_day, count_finish_1_day, un_finish_count, params FROM tms_transport_bind WHERE bind_id = $1 AND (org_id = $2 OR $2 < 1)", args.BindID, args.OrgID)
	if err == nil && data.ID < 1 {
		err = errors.New("bind not exist")
	}
	return
}

// ArgsSetBind 设置绑定关系参数
type ArgsSetBind struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//组织成员ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id"`
	//分区ID
	MapAreaID int64 `db:"map_area_id" json:"mapAreaID" check:"id" empty:"true"`
	//更多分区
	// 可以绑定更多分区，但性能会下降
	MoreMapAreaIDs pq.Int64Array `db:"more_map_area_ids" json:"moreMapAreaIDs" check:"ids" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// SetBind 设置绑定关系
func SetBind(args *ArgsSetBind) (data FieldsBind, err error) {
	data, err = GetBindByBindID(&ArgsGetBindByBindID{
		OrgID:  args.OrgID,
		BindID: args.BindID,
	})
	if err != nil {
		err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "tms_transport_bind", "INSERT INTO tms_transport_bind (org_id, bind_id, map_area_id, more_map_area_ids, km_30_day, level_30_day, time_30_day, count_30_day, count_finish_30_day, km_1_day, level_1_day, time_1_day, count_1_day, count_finish_1_day, params) VALUES (:org_id,:bind_id,:map_area_id,:more_map_area_ids,0,0,0,0,0,0,0,0,0,0,:params)", args, &data)
		return
	} else {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE tms_transport_bind SET delete_at = to_timestamp(0), update_at = NOW(), map_area_id = :map_area_id, more_map_area_ids = :more_map_area_ids, params = :params WHERE id = :id", map[string]interface{}{
			"id":                data.ID,
			"map_area_id":       args.MapAreaID,
			"more_map_area_ids": args.MoreMapAreaIDs,
			"params":            args.Params,
		})
		if err == nil {
			data.MapAreaID = args.MapAreaID
			data.MoreMapAreaIDs = args.MoreMapAreaIDs
			data.Params = args.Params
		}
		return
	}
}

// ArgsDeleteBind 删除绑定关系参数
type ArgsDeleteBind struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteBind 删除绑定关系
func DeleteBind(args *ArgsDeleteBind) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "tms_transport_bind", "id = :id AND (org_id = :org_id OR :org_id < 1)", args)
	return
}

// argsGetBindToTransport 获取符合配送条件的绑定关系参数
type argsGetBindToTransport struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//地图制式
	// WGS-84 / GCJ-02 / BD-09
	MapType int `db:"map_type" json:"mapType"`
	//坐标位置
	Longitude float64 `db:"longitude" json:"longitude"`
	Latitude  float64 `db:"latitude" json:"latitude"`
}

// getBindToTransport 获取符合配送条件的绑定关系
func getBindToTransport(args *argsGetBindToTransport) (resultBind FieldsBind, areaID int64, errCode string, err error) {
	//获取全局权重值
	var globAutoWeightMission, globAutoWeightSpeed, globAutoWeightLevel float64
	globAutoWeightMission, err = OrgCore.Config.GetConfigValFloat64(&ClassConfig.ArgsGetConfig{
		BindID:    args.OrgID,
		Mark:      "TransportAutoWeightMission",
		VisitType: "admin",
	})
	if err != nil {
		globAutoWeightMission = 0.6
	}
	globAutoWeightSpeed, err = OrgCore.Config.GetConfigValFloat64(&ClassConfig.ArgsGetConfig{
		BindID:    args.OrgID,
		Mark:      "TransportAutoWeightSpeed",
		VisitType: "admin",
	})
	if err != nil {
		globAutoWeightMission = 0.3
	}
	globAutoWeightLevel, err = OrgCore.Config.GetConfigValFloat64(&ClassConfig.ArgsGetConfig{
		BindID:    args.OrgID,
		Mark:      "TransportAutoWeightLevel",
		VisitType: "admin",
	})
	if err != nil {
		globAutoWeightMission = 0.1
	}
	//根据定位找到所在分区
	if args.MapType > -1 {
		var areaList []MapArea.FieldsArea
		areaList, err = MapArea.CheckPointInAreas(&MapArea.ArgsCheckPointInAreas{
			MapType: args.MapType,
			Point: CoreSQLGPS.FieldsPoint{
				Longitude: args.Longitude,
				Latitude:  args.Latitude,
			},
			OrgID:    args.OrgID,
			IsParent: false,
			Mark:     "tms",
		})
		if err != nil {
			//找不到分组，则检查该组织下是否不存在任何分区？
			areaList, _, err = MapArea.GetList(&MapArea.ArgsGetList{
				Pages: CoreSQLPages.ArgsDataList{
					Page: 1,
					Max:  1,
					Sort: "id",
					Desc: false,
				},
				OrgID:    args.OrgID,
				Mark:     "tms",
				ParentID: -1,
				Country:  -1,
				City:     -1,
				MapType:  -1,
				IsRemove: false,
				Search:   "",
			})
			if err == nil && len(areaList) > 0 {
				errCode = "not_in_area"
				err = errors.New("no in area")
				return
			} else {
				//说明该组织没有任何分区
				areaID = 0
				// 交给后续继续处理
			}
		} else {
			notAnyBind := true
			var lastBindList []FieldsBind
			//遍历所有分区，找到所有该分区的成员
			for _, v := range areaList {
				var bindList []FieldsBind
				err = Router2SystemConfig.MainDB.Select(&bindList, "SELECT id, create_at, update_at, delete_at, org_id, bind_id, map_area_id, km_30_day, level_30_day, time_30_day, count_30_day, count_finish_30_day, km_1_day, level_1_day, time_1_day, count_1_day, count_finish_1_day, params FROM tms_transport_bind WHERE org_id = $1 AND (map_area_id = $2 OR $2 = ANY(more_map_area_ids)) AND delete_at < to_timestamp(1000000)", args.OrgID, v.ID)
				if err != nil {
					continue
				}
				var vAutoWeightMission, vAutoWeightSpeed, vAutoWeightLevel float64
				var b bool
				vAutoWeightMission, b = v.Params.GetValFloat64("autoWeightMission")
				if !b {
					vAutoWeightMission = globAutoWeightMission
				}
				vAutoWeightSpeed, b = v.Params.GetValFloat64("autoWeightSpeed")
				if !b {
					vAutoWeightSpeed = globAutoWeightSpeed
				}
				vAutoWeightLevel, b = v.Params.GetValFloat64("autoWeightLevel")
				if !b {
					vAutoWeightLevel = globAutoWeightLevel
				}
				var vBind FieldsBind
				vBind, err = getBindToTransportByBindList(args.OrgID, bindList, vAutoWeightMission, vAutoWeightSpeed, vAutoWeightLevel)
				if err != nil {
					break
				}
				lastBindList = append(lastBindList, vBind)
				areaID = v.ID
				notAnyBind = false
			}
			if notAnyBind || len(lastBindList) < 1 {
				errCode = "area_no_bind"
				err = errors.New("area not have any bind")
				return
			}
			//处理最终候选人
			resultBind, err = getBindToTransportByBindList(args.OrgID, lastBindList, globAutoWeightMission, globAutoWeightSpeed, globAutoWeightLevel)
			return
		}
	}
	/**
	//抽取组织所有人员信息，并抽取合适的人员
	// 将组织内所有关联成员写入数据集合内
	var bindList []FieldsBind
	err = Router2SystemConfig.MainDB.Select(&bindList, "SELECT id, create_at, update_at, delete_at, org_id, bind_id, map_area_id, km_30_day, level_30_day, time_30_day, count_30_day, count_finish_30_day, km_1_day, level_1_day, time_1_day, count_1_day, count_finish_1_day, params FROM tms_transport_bind WHERE org_id = $1 AND delete_at < to_timestamp(1000000)", args.OrgID)
	if err != nil || len(bindList) < 1 {
		errCode = "org_no_bind"
		err = errors.New(fmt.Sprint("no in area and org not have any bind, ", err))
		return
	}
	resultBind, err = getBindToTransportByBindList(args.OrgID, bindList, globAutoWeightMission, globAutoWeightSpeed, globAutoWeightLevel)
	*/
	//反馈
	return
}

// 根据需求计算权重值
func getBindToTransportByBindList(orgID int64, bindList []FieldsBind, autoWeightMission float64, autoWeightSpeed float64, autoWeightLevel float64) (data FieldsBind, err error) {
	//如果不存在成员列表
	if len(bindList) < 1 {
		err = errors.New("no bind")
		return
	}
	//检查上班情况
	var bindTransportMustWork bool
	bindTransportMustWork, err = OrgCore.Config.GetConfigValBool(&ClassConfig.ArgsGetConfig{
		BindID:    orgID,
		Mark:      "BindTransportMustWork",
		VisitType: "admin",
	})
	if bindTransportMustWork {
		var newBindList []FieldsBind
		for _, v := range bindList {
			if b := OrgTime.CheckIsWorkByOrgBindID(v.BindID); b {
				newBindList = append(newBindList, v)
			}
		}
		bindList = newBindList
	}
	//如果不存在成员列表
	if len(bindList) < 1 {
		err = errors.New("no bind or no work")
		return
	}
	if len(bindList) == 1 {
		data = bindList[0]
		return
	}
	//比分值
	var lastWeightCount float64 = 0
	for _, v := range bindList {
		var vWeightCount float64 = 0
		if v.Time1Day == 0 {
			vWeightCount = float64(v.Count1Day)*autoWeightMission + float64(v.Level1Day)*autoWeightLevel
		} else {
			vWeightCount = float64(v.Count1Day)*autoWeightMission + float64(int64(v.KM1Day)/v.Time1Day)*autoWeightSpeed + float64(v.Level1Day)*autoWeightLevel
		}
		if lastWeightCount == 0 || data.ID < 1 {
			lastWeightCount = vWeightCount
			data = v
		} else {
			if lastWeightCount > vWeightCount {
				data = v
			}
		}
	}
	return
}
