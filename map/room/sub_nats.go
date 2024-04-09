package MapRoom

import (
	BaseService "github.com/fotomxq/weeekj_core/v5/base/service"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	ServiceUserInfo "github.com/fotomxq/weeekj_core/v5/service/user_info"
	"github.com/nats-io/nats.go"
	"time"
)

func subNats() {
	//档案变更
	CoreNats.SubDataByteNoErr("service_user_info_status", "/service/user_info/status", subNatsInfoUpdateStatus)
	//请求房间统计
	_ = BaseService.SetService(&BaseService.ArgsSetService{
		ExpireAt:     CoreFilter.GetNowTimeCarbon().AddDay().Time,
		Name:         "地图房间统计",
		Description:  "",
		EventSubType: "all",
		Code:         "map_room_analysis",
		EventType:    "nats",
		EventURL:     "/map/room/analysis",
		//TODO:待补充
		EventParams: "",
	})
	CoreNats.SubDataByteNoErr("map_room_analysis", "/map/room/analysis", subNatsRoomAnalysis)
	//自动退出服务状态
	_ = BaseService.SetService(&BaseService.ArgsSetService{
		ExpireAt:     CoreFilter.GetNowTimeCarbon().AddDay().Time,
		Name:         "地图房间服务状态变更",
		Description:  "",
		EventSubType: "all",
		Code:         "map_room_service_status",
		EventType:    "nats",
		EventURL:     "/map/room/service_status",
		//TODO:待补充
		EventParams: "",
	})
	CoreNats.SubDataByteNoErr("map_room_service_status", "/map/room/service_status", subNatsServiceStatus)
	//注册服务
	_ = BaseService.SetService(&BaseService.ArgsSetService{
		ExpireAt:     CoreFilter.GetNowTimeCarbon().AddDay().Time,
		Name:         "地图房间状态变更",
		Description:  "",
		EventSubType: "all",
		Code:         "map_room_status",
		EventType:    "nats",
		EventURL:     "/map/room/status",
		//TODO:待补充
		EventParams: "",
	})
}

// 信息档案发生变更
func subNatsInfoUpdateStatus(_ *nats.Msg, action string, id int64, _ string, _ []byte) {
	switch action {
	case "create":
		infoData, _ := ServiceUserInfo.GetInfoID(&ServiceUserInfo.ArgsGetInfoID{
			ID:    id,
			OrgID: -1,
		})
		if infoData.ID < 1 {
			return
		}
		updateInfoNoRoom(infoData.OrgID)
	case "update":
	case "die":
		updateRoomOut(id)
	case "out":
		updateRoomOut(id)
	case "delete":
		updateRoomOut(id)
	}
}

// 房间统计更新
func subNatsRoomAnalysis(_ *nats.Msg, _ string, id int64, _ string, _ []byte) {
	analysisBlockerWait.CheckWait(id, "", func(modID int64, _ string) {
		updateRoomAnalysis(modID)
	})
}

// 自动退出服务状态
func subNatsServiceStatus(_ *nats.Msg, action string, id int64, _ string, _ []byte) {
	serviceStatusBlockWait.CheckWait(id, action, func(modID int64, modMark string) {
		switch modMark {
		case "no":
			//延迟10秒执行
			time.Sleep(time.Second * 10)
			_, _ = UpdateServiceStatus(&ArgsUpdateServiceStatus{
				ID:               modID,
				OrgID:            -1,
				ServiceStatus:    0,
				ServiceBindID:    0,
				ServiceMissionID: 0,
			})
		case "call":
			_, _ = UpdateServiceStatus(&ArgsUpdateServiceStatus{
				ID:               modID,
				OrgID:            -1,
				ServiceStatus:    1,
				ServiceBindID:    0,
				ServiceMissionID: 0,
			})
		case "ok":
			_, _ = UpdateServiceStatus(&ArgsUpdateServiceStatus{
				ID:               modID,
				OrgID:            -1,
				ServiceStatus:    2,
				ServiceBindID:    0,
				ServiceMissionID: 0,
			})
		}
	})
}
