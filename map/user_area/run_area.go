package MapUserArea

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLGPS "github.com/fotomxq/weeekj_core/v5/core/sql/gps"
	IOTTrack "github.com/fotomxq/weeekj_core/v5/iot/track"
	MapArea "github.com/fotomxq/weeekj_core/v5/map/area"
	OrgCore "github.com/fotomxq/weeekj_core/v5/org/core"
	OrgMission "github.com/fotomxq/weeekj_core/v5/org/mission"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	ServiceUserInfo "github.com/fotomxq/weeekj_core/v5/service/user_info"
	UserGPS "github.com/fotomxq/weeekj_core/v5/user/gps"
	"time"
)

// 超出判断和记录处理
// 自动根据需求，推送任务列
func runArea() {
	//遍历所有带监测数据
	limit := 100
	step := 0
	for {
		//获取列表
		var dataList []FieldsMonitor
		if err := Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, org_id, device_id, user_info_id, area_id, in_range, send_mission, org_group_id FROM map_user_area WHERE delete_at < to_timestamp(1000000) AND is_invalid = false LIMIT $1 OFFSET $2", limit, step); err != nil {
			break
		}
		if len(dataList) < 1 {
			break
		}
		for _, v := range dataList {
			//先检查该用户所持设备的位置
			isFindDeviceTrack, inRange := runAreaDevice(v.DeviceID, v.AreaID)
			var isFindUserTrack bool
			if !isFindDeviceTrack {
				//通过用户查询定位数据包
				isFindUserTrack, inRange = runAreaUser(v.UserInfoID, v.AreaID)
				if !isFindUserTrack {
					if _, err := CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE map_user_area SET update_at = NOW(), is_invalid = true WHERE id = :id", map[string]interface{}{
						"id": v.ID,
					}); err != nil {
						CoreLog.Error("map user area, data invalid, id: ", v.ID, ", err: ", err)
					}
					continue
				}
			}
			//范围不一致，则需更新位置数据包
			if inRange != v.InRange {
				//更新该信息条目，用户超出区域
				if _, err := CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE map_user_area SET update_at = NOW(), in_range = :in_range, send_mission = false", map[string]interface{}{
					"id":       v.ID,
					"in_range": inRange,
				}); err != nil {
					CoreLog.Error("map user area, update map user area, id: ", v.ID, ", err: ", err)
					continue
				}
				//推送预警消息
				if !inRange && v.OrgGroupID > 0 {
					bindData, err := OrgCore.GetBindLast(&OrgCore.ArgsGetBindLast{
						OrgID:   v.OrgID,
						GroupID: v.OrgGroupID,
						Mark:    "map_user_area",
						Params:  []CoreSQLConfig.FieldsConfigType{},
					})
					if err != nil {
						CoreLog.Error("map user area, get bind data, id: ", v.ID, ", err: ", err)
						continue
					}
					//推送行政任务
					_, err = OrgMission.CreateMission(&OrgMission.ArgsCreateMission{
						OrgID:        v.OrgID,
						CreateBindID: bindData.ID,
						BindID:       bindData.ID,
						OtherBindIDs: []int64{},
						Title:        "监护人员超出围栏范围",
						Des:          "监护人员超出范围，请尽快查看地图并前往当前定位区域",
						DesFiles:     []int64{},
						StartAt:      CoreFilter.GetNowTime(),
						EndAt:        CoreFilter.GetNowTime().Add(time.Hour * 6),
						TipID:        0,
						ParentID:     0,
						Level:        2,
						SortID:       0,
						Tags:         []int64{},
						Params:       []CoreSQLConfig.FieldsConfigType{},
					})
					if err != nil {
						CoreLog.Error("map user area, create org mission, id: ", v.ID, ", err: ", err)
						continue
					}
					if _, err := CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE map_user_area SET update_at = NOW(), send_mission = true", map[string]interface{}{
						"id": v.ID,
					}); err != nil {
						CoreLog.Error("map user area, update map user area, send mission, id: ", v.ID, ", err: ", err)
						continue
					}
				}
			}
		}
		step += limit
	}
}

// 获取设备的定位数据包
func runAreaDevice(deviceID int64, areaID int64) (isFindDeviceTrack bool, inArea bool) {
	//获取设备定位数据
	if deviceID < 1 {
		return
	}
	deviceTrack, err := IOTTrack.GetLast(&IOTTrack.ArgsGetLast{
		DeviceID: deviceID,
	})
	if err != nil {
		return
	}
	//拿到定位数据后，检查是否超出区域
	inArea = MapArea.CheckPointInArea(&MapArea.ArgsCheckPointInArea{
		MapType: deviceTrack.MapType,
		Point: CoreSQLGPS.FieldsPoint{
			Longitude: deviceTrack.Longitude,
			Latitude:  deviceTrack.Latitude,
		},
		AreaID: areaID,
	})
	return
}

// 获取用户定位数据
func runAreaUser(userInfoID int64, areaID int64) (isFindUserTrack bool, inArea bool) {
	userInfoData, err := ServiceUserInfo.GetInfoID(&ServiceUserInfo.ArgsGetInfoID{
		ID:    userInfoID,
		OrgID: -1,
	})
	if err != nil {
		return
	}
	//找到用户定位数据
	userTrack, err := UserGPS.GetLast(&UserGPS.ArgsGetLast{
		UserID: userInfoData.UserID,
	})
	if err != nil {
		return
	}
	//检查定位是否在范围内？
	inArea = MapArea.CheckPointInArea(&MapArea.ArgsCheckPointInArea{
		MapType: userTrack.MapType,
		Point: CoreSQLGPS.FieldsPoint{
			Longitude: userTrack.Longitude,
			Latitude:  userTrack.Latitude,
		},
		AreaID: areaID,
	})
	return
}
