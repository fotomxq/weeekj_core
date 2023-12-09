package OrgUser

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	ServiceOrder "gitee.com/weeekj/weeekj_core/v5/service/order"
	"time"
)

// FieldsOrgUser 足迹记录
// 数据将在每次购物等特定行为时更新，确保用户数据被动处于较为新的时间段
// 该数据集合同时方便检索
type FieldsOrgUser struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	// 每次购物、生成订单、购买订阅、获得票据、积分变动，将更新本数据集合
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//用户昵称
	Name string `db:"name" json:"name"`
	//用户联系电话
	Phone string `db:"phone" json:"phone"`
	//地址结构
	AddressList FieldsOrgUserAddress `db:"address_list" json:"addressList"`
	//用户积分总数
	UserIntegral int64 `db:"user_integral" json:"userIntegral"`
	//用户订阅状态
	UserSubs FieldsOrgUserSubs `db:"user_subs" json:"userSubs"`
	//用户票据状态
	UserTickets FieldsOrgUserTickets `db:"user_tickets" json:"userTickets"`
	//储蓄状态
	DepositData FieldsOrgUserDeposits `db:"deposit_data" json:"depositData"`
	//最后一次订单结构体
	LastOrder FieldsOrgUserOrder `db:"last_order" json:"lastOrder"`
	//附加参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// FieldsOrgUserDeposits 储蓄数据集合
type FieldsOrgUserDeposits []FieldsOrgUserDeposit

type FieldsOrgUserDeposit struct {
	//标识码
	Mark string `db:"mark" json:"mark"`
	//金额
	Price int64 `db:"price" json:"price"`
}

func (t FieldsOrgUserDeposits) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsOrgUserDeposits) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

func (t FieldsOrgUserDeposit) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsOrgUserDeposit) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

// FieldsOrgUserAddress 地址信息
type FieldsOrgUserAddress []FieldsAddress

func (t FieldsOrgUserAddress) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsOrgUserAddress) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

// FieldsAddress 通用地址结构
type FieldsAddress struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//所属国家 国家代码
	// eg: china => 86
	Country int `db:"country" json:"country"`
	//省份 编码
	// eg: 710000
	Province int `db:"province" json:"province"`
	//所属城市
	City int `db:"city" json:"city"`
	//街道详细信息
	Address string `db:"address" json:"address"`
	//地图制式
	// WGS-84 / GCJ-02 / BD-09
	MapType int `db:"map_type" json:"mapType"`
	//坐标位置
	Longitude float64 `db:"longitude" json:"longitude"`
	Latitude  float64 `db:"latitude" json:"latitude"`
	//联系人姓名
	Name string `db:"name" json:"name"`
	//联系人国家代码
	NationCode string `db:"nation_code" json:"nationCode"`
	//联系人手机号
	Phone string `db:"phone" json:"phone"`
}

func (t FieldsAddress) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsAddress) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

// FieldsOrgUserSubs 用户订阅
type FieldsOrgUserSubs []FieldsOrgUserSub

type FieldsOrgUserSub struct {
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID"`
	//到期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
}

func (t FieldsOrgUserSubs) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsOrgUserSubs) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

func (t FieldsOrgUserSub) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsOrgUserSub) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

// FieldsOrgUserTickets 用户票据
type FieldsOrgUserTickets []FieldsOrgUserTicket

type FieldsOrgUserTicket struct {
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID"`
	//张数
	Count int64 `db:"count" json:"count"`
	//到期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
}

func (t FieldsOrgUserTickets) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsOrgUserTickets) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

func (t FieldsOrgUserTicket) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsOrgUserTicket) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

// FieldsOrgUserOrder 订单结构
type FieldsOrgUserOrder ServiceOrder.FieldsOrder

func (t FieldsOrgUserOrder) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsOrgUserOrder) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}
