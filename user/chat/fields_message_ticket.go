package UserChat

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

//FieldsMessageTicket 优惠券消息类型绑定关系
type FieldsMessageTicket struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//发起用户
	// 用户必须具备该票据配置足够的张数，否则禁止发起
	// 发起后将使用掉该用户的张数
	UserID int64 `db:"user_id" json:"userID"`
	//参与聊天室
	GroupID int64 `db:"group_id" json:"groupID"`
	//消息ID
	MessageID int64 `db:"message_id" json:"messageID"`
	//票据配置ID
	ConfigID int64 `db:"config_id" json:"configID"`
	//发放张数
	UseCount int64 `db:"use_count" json:"useCount"`
	//发放办法
	// 0 全部发放给领取的第一个人; 1 随机发放(1张 - 总张数-尚未领取人数x1张); 2 每人限制1张
	TakeType int `db:"take_type" json:"takeType"`
	//领取人数限制
	CountLimit int `db:"count_limit" json:"countLimit"`
	//领取人列表
	TakeList FieldsMessageTicketTakeList `db:"take_list" json:"takeList"`
}

type FieldsMessageTicketTakeList []FieldsMessageTicketTake

//Value sql底层处理器
func (t FieldsMessageTicketTakeList) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsMessageTicketTakeList) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsMessageTicketTake struct {
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//领取用户
	UserID int64 `db:"user_id" json:"userID"`
	//领取张数
	GetCount int64 `db:"get_count" json:"getCount"`
}

//Value sql底层处理器
func (t FieldsMessageTicketTake) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsMessageTicketTake) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}
