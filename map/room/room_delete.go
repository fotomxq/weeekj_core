package MapRoom

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsDeleteRoom 删除房间参数
type ArgsDeleteRoom struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteRoom 删除房间
func DeleteRoom(args *ArgsDeleteRoom) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "map_room", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	if args.OrgID < 1 {
		var data FieldsRoom
		data, err = GetRoomID(&ArgsGetRoomID{
			ID:    args.ID,
			OrgID: -1,
		})
		if err == nil {
			args.OrgID = data.OrgID
		}
	}
	//删除缓冲
	deleteRoomCache(args.ID)
	//推送nats
	pushNatsUpdateStatus(args.ID, "delete", "")
	//更新统计
	pushNatsUpdateAnalysis(args.OrgID)
	//反馈
	return
}

// 重新核准房间的档案ID信息
// 剔除指定的档案
func updateRoomOut(infoID int64) {
	roomList, err := GetRoomListByInfo(infoID)
	if err != nil {
		return
	}
	for _, v := range roomList {
		var newInfos pq.Int64Array
		for _, v2 := range v.Infos {
			if v2 == infoID {
				continue
			}
			newInfos = append(newInfos, v2)
		}
		_, err := CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE map_room SET infos = :infos WHERE id = :id", map[string]interface{}{
			"id":    v.ID,
			"infos": newInfos,
		})
		if err != nil {
			continue
		}
		//删除缓冲
		deleteRoomCache(v.ID)
		//推送nats
		pushNatsUpdateStatus(v.ID, "update_user_info", "")
		//更新统计
		pushNatsUpdateAnalysis(v.OrgID)
		//请求更新入驻名单数据
		updateInfoNoRoom(v.OrgID)
	}
}
