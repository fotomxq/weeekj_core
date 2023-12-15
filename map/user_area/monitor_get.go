package MapUserArea

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetMonitorList 获取监控关系列表参数
type ArgsGetMonitorList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//是否失效
	NeedIsInvalid bool `db:"need_is_invalid" json:"needIsInvalid" check:"bool"`
	IsInvalid     bool `db:"is_invalid" json:"isInvalid" check:"bool"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//看护人档案ID
	// 如果没有设备，将根据档案ID查询用户GPS讯号
	UserInfoID int64 `db:"user_info_id" json:"userInfoID" check:"id" empty:"true"`
	//绑定的设备ID
	// 该设备被视为此人的GPS移动讯号
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id" empty:"true"`
	//电子围栏ID
	// 超出该围栏范围将推送预警消息
	AreaID int64 `db:"area_id" json:"areaID" check:"id" empty:"true"`
	//当前是否超出区域？
	InRange bool `db:"in_range" json:"inRange" check:"bool" empty:"true"`
	//任务推送给哪个组的成员？
	OrgGroupID int64 `db:"org_group_id" json:"orgGroupID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
}

// GetMonitorList 获取监控关系列表
func GetMonitorList(args *ArgsGetMonitorList) (dataList []FieldsMonitor, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.NeedIsInvalid {
		where = where + " AND is_invalid = :is_invalid"
		maps["is_invalid"] = args.IsInvalid
	}
	if args.UserInfoID > -1 {
		where = where + " AND user_info_id = :user_info_id"
		maps["user_info_id"] = args.UserInfoID
	}
	if args.DeviceID > -1 {
		where = where + " AND device_id = :device_id"
		maps["device_id"] = args.DeviceID
	}
	if args.AreaID > -1 {
		where = where + " AND area_id = :area_id"
		maps["area_id"] = args.AreaID
	}
	if args.InRange {
		where = where + " AND in_range = :in_range"
		maps["in_range"] = args.InRange
	}
	if args.OrgGroupID > -1 {
		where = where + " AND org_group_id = :org_group_id"
		maps["org_group_id"] = args.OrgGroupID
	}
	tableName := "map_user_area"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, update_at, delete_at, org_id, is_invalid, user_info_id, device_id, area_id, in_range, org_group_id, send_mission, params FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	return
}
