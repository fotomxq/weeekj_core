package MapUserArea

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsUpdateMonitor 修改自动化参数
type ArgsUpdateMonitor struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
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

// UpdateMonitor 修改自动化
func UpdateMonitor(args *ArgsUpdateMonitor) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE map_user_area SET update_at = NOW(), is_invalid = false, user_info_id = :user_info_id, device_id = :device_id, area_id = :area_id, org_group_id = :org_group_id, send_mission = false, params = :params WHERE id = :id AND org_id = :org_id AND delete_at < to_timestamp(1000000)", args)
	return
}
