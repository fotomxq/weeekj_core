package ToolsCommunication

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"time"
)

// FieldsFrom 参与来源
type FieldsFrom struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//到期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	//房间ID
	RoomID int64 `db:"room_id" json:"roomID"`
	//来源系统
	// 1 用户; 2 设备
	FromSystem int `db:"from_system" json:"fromSystem"`
	//来源ID
	FromID int64 `db:"from_id" json:"fromID"`
	//昵称
	Name string `db:"name" json:"name"`
	//链接token
	// 用于第三方链接用
	Token string `db:"token" json:"token"`
	//是否允许发言
	AllowSend bool `db:"allow_send" json:"allowSend"`
	//角色类型
	// 0 普通; 1 房主; 2 副房主
	Role int `db:"role" json:"role"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
