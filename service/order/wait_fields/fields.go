package ServiceOrderWaitFields

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	CoreSQLAddress "github.com/fotomxq/weeekj_core/v5/core/sql/address"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	"time"
)

type FieldsWait struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//订单ID
	OrderID int64 `db:"order_id" json:"orderID"`
	//来源系统
	// 该订单创建来源的系统
	// eg: user_sub / org_sub / mall
	SystemMark string `db:"system_mark" json:"systemMark"`
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//创建来源
	// 0 商户平台下单、1 用户APP下单、2 微信小程序下单、3 支付宝小程序下单、4 电话系统下单、5 线下柜台下单、6 辅助设备下单、7 其他渠道
	CreateFrom int `db:"create_from" json:"createFrom"`
	//hash
	Hash string `db:"hash" json:"hash"`
	//收取货物地址
	AddressFrom CoreSQLAddress.FieldsAddress `db:"address_from" json:"addressFrom"`
	//送货地址
	AddressTo CoreSQLAddress.FieldsAddress `db:"address_to" json:"addressTo"`
	//货物清单
	Goods FieldsGoods `db:"goods" json:"goods"`
	//订单总的抵扣
	// 例如满减活动，不局限于个别商品的活动
	Exemptions FieldsExemptions `db:"exemptions" json:"exemptions"`
	//是否允许自动审核
	// 客户提交订单后，将自动审核该订单。订单如果存在至少一件未开启的商品，将禁止该操作
	AllowAutoAudit bool `db:"allow_auto_audit" json:"allowAutoAudit"`
	//允许自动配送
	TransportAllowAuto bool `db:"transport_allow_auto" json:"transportAllowAuto"`
	//期望送货时间
	TransportTaskAt time.Time `db:"transport_task_at" json:"transportTaskAt"`
	//是否允许货到付款？
	TransportPayAfter bool `db:"transport_pay_after" json:"transportPayAfter"`
	//配送服务系统
	// 0 self 其他配送; 1 take 自提; 2 transport 自运营配送; 3 running 跑腿服务; 4 housekeeping 家政服务
	TransportSystem string `db:"transport_system" json:"transportSystem"`
	//费用组成
	PriceList FieldsPrices `db:"price_list" json:"priceList"`
	//订单总费用
	// 总费用是否支付
	PricePay bool `db:"price_pay" json:"pricePay"`
	// 货币
	Currency int `db:"currency" json:"currency"`
	// 总费用金额
	Price int64 `db:"price" json:"price"`
	//折扣前费用
	PriceTotal int64 `db:"price_total" json:"priceTotal"`
	//备注信息
	Des string `db:"des" json:"des"`
	//日志
	Logs FieldsLogs `db:"logs" json:"logs"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
	//失败信息
	ErrCode string `db:"err_code" json:"errCode"`
	ErrMsg  string `db:"err_msg" json:"errMsg"`
}

// FieldsGoods 货物
type FieldsGoods []FieldsGood

// Value sql底层处理器
func (t FieldsGoods) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsGoods) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsGood struct {
	//获取来源
	// 如果商品mark带有virtual标记，且订单商品全部带有该标记，订单将在付款后直接完成
	From CoreSQLFrom.FieldsFrom `db:"from" json:"from"`
	//选项Key
	// 如果给空，则该商品必须也不包含选项
	OptionKey string `db:"option_key" json:"optionKey" check:"mark" empty:"true"`
	//货物个数
	Count int64 `db:"count" json:"count" check:"intThan0" empty:"true"`
	//获取价值
	// 单个商品价值
	Price int64 `db:"price" json:"price"`
	//抵扣
	Exemptions FieldsExemptions `db:"exemptions" json:"exemptions"`
}

// FieldsPrices 费用组成
type FieldsPrices []FieldsPrice

// Value sql底层处理器
func (t FieldsPrices) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsPrices) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsPrice struct {
	//费用类型
	// 0 货物费用；1 配送费用；2 保险费用; 3 跑腿费用
	PriceType int `db:"price_type" json:"priceType" check:"mark"`
	//是否缴费
	IsPay bool `db:"is_pay" json:"isPay" check:"bool"`
	//总金额
	Price int64 `db:"price" json:"price" check:"price"`
}

// FieldsExemptions 抵扣结构
type FieldsExemptions []FieldsExemption

// Value sql底层处理器
func (t FieldsExemptions) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsExemptions) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsExemption struct {
	//抵扣系统来源
	// integral 积分; ticket 票据; sub 订阅
	System string `db:"system" json:"system"`
	//抵扣配置ID
	// 可能不存在，如积分没有配置ID
	ConfigID int64 `db:"config_id" json:"configID"`
	//抵扣名称
	// eg: 订阅X
	Name string `db:"name" json:"name"`
	//抵扣描述信息
	// eg: 票据X使用3张，减免13元
	Des string `db:"des" json:"des"`
	//使用数量
	// 使用的张数、或使用积分的个数
	Count int64 `db:"count" json:"count"`
	//抵扣费用
	// 总的费用，含商品复数。例如抵扣多个商品N元；或抵扣该多个商品的百分比
	Price int64 `db:"price" json:"price"`
}

// FieldsPays 支付请求结构体
type FieldsPays []FieldsPay

// Value sql底层处理器
func (t FieldsPays) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsPays) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsPay struct {
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//支付ID
	PayID int64 `db:"pay_id" json:"payID" check:"id"`
	//付费状态
	PayStatus int `db:"pay_status" json:"payStatus"`
	//缴费费用
	Price int64 `db:"price" json:"price" check:"price" empty:"true"`
	//支付错误信息
	PayErrorCode string `db:"pay_error_code" json:"payErrorCode"`
	PayErrorMsg  string `db:"pay_error_msg" json:"payErrorMsg"`
}

// FieldsLogs 日志记录
type FieldsLogs []FieldsLog

// Value sql底层处理器
func (t FieldsLogs) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsLogs) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsLog struct {
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//操作用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//操作组织人员ID
	// 如果留空则说明为系统自动调整或创建人产生
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID"`
	//调整说明标识码
	Mark string `db:"mark" json:"mark"`
	//调整说明描述
	Des string `db:"des" json:"des"`
}
