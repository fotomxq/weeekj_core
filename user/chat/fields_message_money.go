package UserChat

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

//FieldsMessageMoney 红包消息类型绑定关系
type FieldsMessageMoney struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//发起用户
	// 该用户必须具有大于资金量
	UserID int64 `db:"user_id" json:"userID"`
	//参与聊天室
	GroupID int64 `db:"group_id" json:"groupID"`
	//消息ID
	MessageID int64 `db:"message_id" json:"messageID"`
	//储蓄配置
	ConfigMark string `db:"config_mark" json:"configMark" check:"mark" empty:"true"`
	//金额
	Price int64 `db:"price" json:"price"`
	//发放办法
	// 0 全部发放给领取的第一个人; 1 随机发放(0.01 - 总金额-尚未领取人数x0.01)
	TakeType int `db:"take_type" json:"takeType"`
	//领取人数限制
	CountLimit int `db:"count_limit" json:"countLimit"`
	//领取人列表
	TakeList FieldsMessageMoneyTakeList `db:"take_list" json:"takeList"`
}

type FieldsMessageMoneyTakeList []FieldsMessageMoneyTake

//Value sql底层处理器
func (t FieldsMessageMoneyTakeList) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsMessageMoneyTakeList) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsMessageMoneyTake struct {
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//领取用户
	UserID int64 `db:"user_id" json:"userID"`
	//支付ID
	PayID int64 `db:"pay_id" json:"payID"`
	//金额
	Price int64 `db:"price" json:"price"`
}

//Value sql底层处理器
func (t FieldsMessageMoneyTake) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsMessageMoneyTake) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}
