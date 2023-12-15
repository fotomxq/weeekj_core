package IOTMQTT

import (
	"fmt"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	IOTLog "github.com/fotomxq/weeekj_core/v5/iot/log"
	IOTMission "github.com/fotomxq/weeekj_core/v5/iot/mission"
)

func runMission() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("iot device mission run error, ", r)
		}
	}()
	//检查阻拦器
	if !IOTMission.RunMissionBlocker.CheckPass() {
		return
	}
	//获取任务列
	missionList, err := IOTMission.PickMissionList()
	if err != nil {
		return
	}
	for _, v := range missionList {
		switch v.ConnectType {
		case "mqtt_client":
			//推送任务
			missionData, err := IOTMission.PickMission(&IOTMission.ArgsPickMission{
				ID: v.ID,
			})
			if err != nil {
				CoreLog.Error("iot device mission run, get mission data, id: ", v.ID, ", err: ", err)
				continue
			}
			err = PushMissionToDevice(&ArgsPushMissionToDevice{
				ID:         missionData.DeviceID,
				MissionID:  missionData.ID,
				ExpireAt:   missionData.ExpireAt,
				Action:     missionData.Action,
				ParamsData: missionData.ParamsData,
			})
			if err != nil {
				CoreLog.Error("iot device mission run, push mission to device, ", err)
				err = IOTMission.UpdateMissionStatus(&IOTMission.ArgsUpdateMissionStatus{
					ID:     missionData.ID,
					Status: 3,
				})
				if err != nil {
					CoreLog.Error("iot device mission run, push mission to device, mission failed, ", err)
				} else {
					IOTLog.Append(&IOTLog.ArgsAppend{
						OrgID:    missionData.OrgID,
						GroupID:  missionData.GroupID,
						DeviceID: missionData.DeviceID,
						Mark:     "send_failed",
						Content:  fmt.Sprint("无法向设备推送任务[", missionData.ID, "]"),
					})
					_ = IOTMission.CreateLog(&IOTMission.ArgsCreateLog{
						MissionID: missionData.ID,
						Status:    3,
						Mark:      "send_device_failed",
						Content:   fmt.Sprint("无法向设备推送任务[", missionData.ID, "]"),
					})
				}
				continue
			}
			IOTLog.Append(&IOTLog.ArgsAppend{
				OrgID:    missionData.OrgID,
				GroupID:  missionData.GroupID,
				DeviceID: missionData.DeviceID,
				Mark:     "send",
				Content:  fmt.Sprint("向设备推送任务[", missionData.ID, "]"),
			})
			_ = IOTMission.CreateLog(&IOTMission.ArgsCreateLog{
				MissionID: missionData.ID,
				Status:    1,
				Mark:      "send_device",
				Content:   fmt.Sprint("向设备推送任务[", missionData.ID, "]"),
			})
		case "mqtt_group":
			//推送任务
			missionData, err := IOTMission.PickMission(&IOTMission.ArgsPickMission{
				ID: v.ID,
			})
			if err != nil {
				CoreLog.Error("iot device mission run, get mission data, id: ", v.ID, ", err: ", err)
				continue
			}
			err = PushMissionToGroup(ArgsPushMissionToGroup{
				ID:         missionData.DeviceID,
				MissionID:  missionData.ID,
				GroupID:    missionData.GroupID,
				ExpireAt:   missionData.ExpireAt,
				Action:     missionData.Action,
				ParamsData: missionData.ParamsData,
			})
			if err != nil {
				CoreLog.Error("iot device mission run, push mission to group, ", err)
				err = IOTMission.UpdateMissionStatus(&IOTMission.ArgsUpdateMissionStatus{
					ID:     missionData.ID,
					Status: 3,
				})
				if err != nil {
					CoreLog.Error("iot device mission run, push mission to device, mission failed, ", err)
				} else {
					IOTLog.Append(&IOTLog.ArgsAppend{
						OrgID:    missionData.OrgID,
						GroupID:  missionData.GroupID,
						DeviceID: missionData.DeviceID,
						Mark:     "send_failed",
						Content:  fmt.Sprint("无法向设备组推送任务[", missionData.ID, "]"),
					})
					_ = IOTMission.CreateLog(&IOTMission.ArgsCreateLog{
						MissionID: missionData.ID,
						Status:    3,
						Mark:      "send_group_failed",
						Content:   fmt.Sprint("无法向设备组推送任务[", missionData.ID, "]"),
					})
				}
				continue
			}
			IOTLog.Append(&IOTLog.ArgsAppend{
				OrgID:    missionData.OrgID,
				GroupID:  missionData.GroupID,
				DeviceID: missionData.DeviceID,
				Mark:     "send",
				Content:  fmt.Sprint("向设备分组推送任务[", missionData.ID, "]"),
			})
			_ = IOTMission.CreateLog(&IOTMission.ArgsCreateLog{
				MissionID: missionData.ID,
				Status:    1,
				Mark:      "send_group",
				Content:   fmt.Sprint("向设备分组推送任务[", missionData.ID, "]"),
			})
		case "none":
			return
		default:
			return
		}
	}
}
