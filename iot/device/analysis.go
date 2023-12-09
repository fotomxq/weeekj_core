package IOTDevice

import (
	"errors"
	"fmt"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetAnalysisIsOnline 获取设备在线情况参数
type ArgsGetAnalysisIsOnline struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//设备分组
	GroupID int64 `db:"group_id" json:"groupID" check:"id" empty:"true"`
	//组织分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//组织标签ID组
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
}

// GetAnalysisIsOnlineData 获取设备在线情况结构
type GetAnalysisIsOnlineData struct {
	//设备总数
	DeviceCount int64 `db:"device_count" json:"deviceCount"`
	//在线数量
	OnlineCount int64 `db:"online_count" json:"onlineCount"`
}

// GetAnalysisIsOnline 获取设备在线情况
func GetAnalysisIsOnline(args *ArgsGetAnalysisIsOnline) (data GetAnalysisIsOnlineData, err error) {
	maps := map[string]interface{}{}
	where := ""
	if args.OrgID > 0 {
		where = "d.delete_at < to_timestamp(1000000)"
		maps["org_id"] = args.OrgID
		if args.GroupID > 0 {
			where = where + " AND d.group_id = :group_id"
			maps["group_id"] = args.GroupID
		}
		if args.SortID > 0 {
			where = where + " AND o.sort_id = :sort_id"
			maps["sort_id"] = args.SortID
		}
		if len(args.Tags) > 0 {
			where = where + " AND o.tags @> :tags"
			maps["tags"] = args.Tags
		}
		data.DeviceCount, err = CoreSQL.GetAllCountMapTables(
			Router2SystemConfig.MainDB.DB,
			"d.id",
			"iot_core_device as d INNER JOIN iot_core_operate as o ON d.id = o.device_id AND o.delete_at < to_timestamp(1000000) WHERE o.org_id = :org_id AND "+where,
			maps,
		)
		data.OnlineCount, err = CoreSQL.GetAllCountMapTables(
			Router2SystemConfig.MainDB.DB,
			"d.id",
			"iot_core_device as d INNER JOIN iot_core_operate as o ON d.id = o.device_id AND o.delete_at < to_timestamp(1000000) WHERE o.org_id = :org_id AND d.is_online = true AND "+where,
			maps,
		)
	} else {
		where = "delete_at < to_timestamp(1000000)"
		if args.GroupID > 0 {
			where = where + " AND group_id = :group_id"
			maps["group_id"] = args.GroupID
		}
		data.DeviceCount, err = CoreSQL.GetAllCountMap(
			Router2SystemConfig.MainDB.DB,
			"iot_core_device",
			"id",
			where,
			maps,
		)
		if err != nil {
			err = errors.New(fmt.Sprint("get all device count, ", err))
			return
		}
		data.OnlineCount, err = CoreSQL.GetAllCountMap(
			Router2SystemConfig.MainDB.DB,
			"iot_core_device",
			"id",
			"is_online = true AND "+where,
			maps,
		)
		if err != nil {
			err = errors.New(fmt.Sprint("get all online count, ", err))
			return
		}
	}
	return
}
