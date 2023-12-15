package UserChat

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"time"
)

// FieldsMessage 聊天离线消息
type FieldsMessage struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//参与聊天室
	GroupID int64 `db:"group_id" json:"groupID"`
	//发起用户
	UserID int64 `db:"user_id" json:"userID"`
	//消息类型
	// 0 普通消息; 1 红包领取; 2 优惠券; 3 定位坐标; 4 语音消息; 5 提示信息（发起了视频、语音通话）
	MessageType int `db:"message_type" json:"messageType"`
	//消息内容
	/**
	0 普通消息，存储消息内容
	4 语音消息，存储语音转化为base64文本数据包
	*/
	Message string `db:"message" json:"message"`
	//扩展参数
	/**
	3 定位坐标中，此处将存储坐标系统的address_xy: xy位置、address: 地址详情
	4 语音消息，此处将存储语音转译后的文字信息message_text: 语音消息文本
	*/
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
