package MapUserArea

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsCreateMonitor 创建新的自动化参数
type ArgsCreateMonitor struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//看护人档案ID
	// 如果没有设备，将根据档案ID查询用户GPS讯号
	UserInfoID int64 `db:"user_info_id" json:"userInfoID" check:"id"`
	//绑定的设备ID
	// 该设备被视为此人的GPS移动讯号
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id" empty:"true"`
	//电子围栏ID
	// 超出该围栏范围将推送预警消息
	AreaID int64 `db:"area_id" json:"areaID" check:"id"`
	//任务推送给哪个组的成员？
	OrgGroupID int64 `db:"org_group_id" json:"orgGroupID" check:"id"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateMonitor 创建新的自动化
func CreateMonitor(args *ArgsCreateMonitor) (data FieldsMonitor, err error) {
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "map_user_area", "INSERT INTO map_user_area (org_id, is_invalid, user_info_id, device_id, area_id, in_range, org_group_id, send_mission, params) VALUES (:org_id,false,:user_info_id,:device_id,:area_id,true,:org_group_id,false,:params)", args, &data)
	return
}
