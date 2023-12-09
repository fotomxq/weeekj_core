package MapRoom

import (
	"fmt"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	ServiceUserInfo "gitee.com/weeekj/weeekj_core/v5/service/user_info"
)

type DataInfoNoInRoom struct {
	//人员ID
	InfoID int64 `json:"infoID"`
	//状态
	// 0 等待; 1 离开
	Status int `json:"status"`
}

// GetInfoNoInRoom 没有在房间的人员统计
func GetInfoNoInRoom(orgID int64, page int, limit int) (dataList []DataInfoNoInRoom) {
	//获取步数
	step := (page - 1) * limit
	if step < 1 {
		step = 0
	}
	max := step + limit
	//获取缓冲数据
	cacheMark := getInfoNoInRoomCacheMark(orgID)
	var rawList []DataInfoNoInRoom
	err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &rawList)
	if err == nil {
		//抽取数据
		if len(rawList) <= step {
			return
		}
		if len(rawList) >= max {
			dataList = rawList[step:max]
			return
		} else {
			dataList = rawList[step:]
			return
		}
	}
	//不存在数据，则建立数据
	var infoPage int64 = 1
	for {
		infoList1, _, _ := ServiceUserInfo.GetInfoList(&ServiceUserInfo.ArgsGetInfoList{
			Pages: CoreSQLPages.ArgsDataList{
				Page: infoPage,
				Max:  1000,
				Sort: "id",
				Desc: false,
			},
			OrgID:    orgID,
			UserID:   -1,
			BindID:   -1,
			Country:  -1,
			SortID:   -1,
			Tags:     []int64{},
			Director: -1,
			IsDie:    false,
			IsOut:    false,
			IsRemove: false,
			Search:   "",
		})
		for _, vInfo := range infoList1 {
			vRoom, _ := GetRoomByInfo(&ArgsGetRoomByInfo{
				OrgID:  vInfo.OrgID,
				InfoID: vInfo.ID,
				Status: -1,
			})
			if vRoom.ID < 1 {
				rawList = append(rawList, DataInfoNoInRoom{
					InfoID: vInfo.ID,
					Status: 0,
				})
			}
		}
		infoList2, _, _ := ServiceUserInfo.GetInfoList(&ServiceUserInfo.ArgsGetInfoList{
			Pages: CoreSQLPages.ArgsDataList{
				Page: infoPage,
				Max:  1000,
				Sort: "id",
				Desc: false,
			},
			OrgID:    orgID,
			UserID:   -1,
			BindID:   -1,
			Country:  -1,
			SortID:   -1,
			Tags:     []int64{},
			Director: -1,
			IsDie:    false,
			IsOut:    true,
			IsRemove: false,
			Search:   "",
		})
		for _, vInfo := range infoList2 {
			vRoom, _ := GetRoomByInfo(&ArgsGetRoomByInfo{
				OrgID:  vInfo.OrgID,
				InfoID: vInfo.ID,
				Status: -1,
			})
			if vRoom.ID < 1 {
				rawList = append(rawList, DataInfoNoInRoom{
					InfoID: vInfo.ID,
					Status: 1,
				})
			}
		}
		if len(infoList1) < 1 && len(infoList2) < 1 {
			break
		}
		//下一页
		infoPage += 1
	}
	//写入缓冲
	Router2SystemConfig.MainCache.SetStruct(cacheMark, rawList, 2419200)
	//抽取数据并反馈
	if len(rawList) <= step {
		return
	}
	if len(rawList) >= max {
		dataList = rawList[step:max]
		return
	} else {
		dataList = rawList[step:]
		return
	}
}

// 发生入驻
func updateInfoNoRoom(orgID int64) {
	deleteInfoNoInRoomCache(orgID)
	GetInfoNoInRoom(orgID, 1, 1)
}

// 缓冲
func getInfoNoInRoomCacheMark(orgID int64) string {
	return fmt.Sprint("map:room:no:info:", orgID)
}

func deleteInfoNoInRoomCache(orgID int64) {
	Router2SystemConfig.MainCache.DeleteMark(getInfoNoInRoomCacheMark(orgID))
}
