package MapRoom

import CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"

// pushNatsUpdateStatus 推送消息中间件
func pushNatsUpdateStatus(roomID int64, action string, mark string) {
	CoreNats.PushDataNoErr("map_room_status", "/map/room/status", action, roomID, mark, nil)
}

// pushNatsUpdateAnalysis 请求更新房间统计
func pushNatsUpdateAnalysis(orgID int64) {
	CoreNats.PushDataNoErr("map_room_analysis", "/map/room/analysis", "", orgID, "", nil)
}

// 请求退出服务状态
func pushNatsServiceStatus(action string, roomID int64) {
	CoreNats.PushDataNoErr("map_room_service_status", "/map/room/service_status", action, roomID, "", nil)
}
