package MapRoom

import (
	CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"
	ServiceUserInfo "gitee.com/weeekj/weeekj_core/v5/service/user_info"
	"github.com/nats-io/nats.go"
	"time"
)

func subNats() {
	//档案变更
	CoreNats.SubDataByteNoErr("/service/user_info/status", subNatsInfoUpdateStatus)
	//请求房间统计
	CoreNats.SubDataByteNoErr("/map/room/analysis", subNatsRoomAnalysis)
	//自动退出服务状态
	CoreNats.SubDataByteNoErr("/map/room/service_status", subNatsServiceStatus)
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
