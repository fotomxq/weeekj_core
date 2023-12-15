package MapUserArea

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"time"
)

// FieldsMonitor 重点看护人GPS、设备绑定
type FieldsMonitor struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//是否失效
	IsInvalid bool `db:"is_invalid" json:"isInvalid"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//看护人档案ID
	// 如果没有设备，将根据档案ID查询用户GPS讯号
	UserInfoID int64 `db:"user_info_id" json:"userInfoID"`
	//绑定的设备ID
	// 该设备被视为此人的GPS移动讯号
	DeviceID int64 `db:"device_id" json:"deviceID"`
	//电子围栏ID
	// 超出该围栏范围将推送预警消息
	AreaID int64 `db:"area_id" json:"areaID"`
	//当前是否超出区域？
	InRange bool `db:"in_range" json:"inRange"`
	//任务推送给哪个组的成员？
	OrgGroupID int64 `db:"org_group_id" json:"orgGroupID"`
	//是否已经推送了任务
	SendMission bool `db:"send_mission" json:"sendMission"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
