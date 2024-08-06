package ServiceOrder

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	CoreSQLAddress "github.com/fotomxq/weeekj_core/v5/core/sql/address"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	"github.com/lib/pq"
	"time"
)

// FieldsOrder 订单
type FieldsOrder struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	// TODO: 逐步取消该设计，当前取消订单会标记状态，不会删除订单
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//过期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	//来源系统
	// 该订单创建来源的系统
	// eg: user_sub 用户订阅 / org_sub 商户订阅 / mall 普通商城 / core_api API服务
	SystemMark string `db:"system_mark" json:"systemMark"`
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//公司ID
	CompanyID int64 `db:"company_id" json:"companyID"`
	//创建来源
	// 0 商户平台下单、1 用户APP下单、2 微信小程序下单、3 支付宝小程序下单、4 电话系统下单、5 线下柜台下单、6 辅助设备下单、7 其他渠道
	CreateFrom int `db:"create_from" json:"createFrom"`
	//编号器，提供累计编号
	SerialNumber int64 `db:"serial_number" json:"serialNumber"`
	//当天的编号
	SerialNumberDay int64 `db:"serial_number_day" json:"serialNumberDay"`
	//状态
	// 0 草稿等待提交
	// 1 提交等待审核中
	// 2 送货中，内部状态根据配送状态确认
	// 3 送货完成，可能包含货到付款
	// 4 送货完成且付款完成
	// 5 订单失败，发货失败等因素
	// 6 取消，包括超时、人为因素
	Status int `db:"status" json:"status"`
	//退货状态
	// 0 没有退货申请
	// 1 提交退货申请
	// 2 退货中
	// 3 退货完成，退款需配合pay_status进行
	// 4 退货失败
	RefundStatus int `db:"refund_status" json:"refundStatus"`
	//退货原因
	RefundWay string `db:"refund_way" json:"refundWay" check:"des" min:"1" max:"600" empty:"true"`
	//退货备注
	RefundDes string `db:"refund_des" json:"refundDes" check:"des" min:"1" max:"1000" empty:"true"`
	//退货图片列
	RefundFileIDs pq.Int64Array `db:"refund_file_ids" json:"refundFileIDs" check:"ids" empty:"true"`
	//退货是否收到货物
	RefundHaveGood bool `db:"refund_have_good" json:"refundHaveGood" check:"bool"`
	//退货快递类型
	// 0 self 其他配送; 1 take 自提; 2 transport 自运营配送; 3 running 跑腿服务; 4 housekeeping 家政服务
	RefundTransportSystem string `db:"refund_transport_system" json:"refundTransportSystem"`
	//退货快递单号
	RefundTransportSN string `db:"refund_transport_sn" json:"refundTransportSN"`
	//配送服务的状态信息
	RefundTransportInfo string `db:"refund_transport_info" json:"refundTransportInfo"`
	//退货支付ID
	RefundPayID int64 `db:"refund_pay_id" json:"refundPayID" check:"id"`
	//退货金额
	RefundPrice int64 `db:"refund_price" json:"refundPrice" check:"price"`
	//退款是否完成
	RefundPayFinish time.Time `db:"refund_pay_finish" json:"refundPayFinish"`
	//退货到期时间
	RefundExpireAt time.Time `db:"refund_expire_at" json:"refundExpireAt"`
	//退货催促时间
	RefundTipAt time.Time `db:"refund_tip_at" json:"refundTipAt"`
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
	//配送ID
	TransportID int64 `db:"transport_id" json:"transportID"`
	//允许自动配送
	TransportAllowAuto bool `db:"transport_allow_auto" json:"transportAllowAuto"`
	//期望送货时间
	TransportTaskAt time.Time `db:"transport_task_at" json:"transportTaskAt"`
	//是否允许货到付款？
	TransportPayAfter bool `db:"transport_pay_after" json:"transportPayAfter"`
	//历史配送ID序列
	TransportIDs pq.Int64Array `db:"transport_ids" json:"transportIDs"`
	//配送服务系统
	// 0 self 其他配送; 1 take 自提; 2 transport 自运营配送; 3 running 跑腿服务; 4 housekeeping 家政服务
	TransportSystem string `db:"transport_system" json:"transportSystem"`
	//配送单号
	TransportSN string `db:"transport_sn" json:"transportSN"`
	//配送服务的状态信息
	TransportInfo string `db:"transport_info" json:"transportInfo"`
	//配送状态
	// 0 等待分配人员; 1 取货中; 2 送货中; 3 完成配送
	TransportStatus int `db:"transport_status" json:"transportStatus"`
	//费用组成
	PriceList FieldsPrices `db:"price_list" json:"priceList"`
	//订单总费用
	// 总费用是否支付
	// 该设计和payStatus并列，但不冲突。因为payStatus可能为退款状态
	PricePay bool `db:"price_pay" json:"pricePay"`
	// 货币
	Currency int `db:"currency" json:"currency"`
	// 总费用金额
	Price int64 `db:"price" json:"price"`
	//折扣前费用
	PriceTotal int64 `db:"price_total" json:"priceTotal"`
	//付费状态
	// 0 尚未付款
	// 1 已经付款
	// 2 发起退款
	// 3 完成退款
	// 4 支付失败
	// 5 退款失败
	PayStatus int `db:"pay_status" json:"payStatus"`
	//当前匹配的支付ID
	PayID int64 `db:"pay_id" json:"payID"`
	//支付ID列
	// 所有关联请求，最后一条为最新的匹配数据
	PayList pq.Int64Array `db:"pay_list" json:"payList"`
	//支付渠道
	PayFrom string `db:"pay_from" json:"payFrom"`
	//支付时间
	PayAt string `db:"pay_at" json:"payAt"`
	//备注信息
	Des string `db:"des" json:"des"`
	//日志
	Logs FieldsLogs `db:"logs" json:"logs"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
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
	// system: user_sub 用户会员 / org_sub 组织会员 / mall 商品 / core_api API服务
	From CoreSQLFrom.FieldsFrom `db:"from" json:"from"`
	//选项Key
	// 如果给空，则该商品必须也不包含选项
	OptionKey string `db:"option_key" json:"optionKey" check:"mark" empty:"true"`
	//货物个数
	Count int64 `db:"count" json:"count"`
	//获取价值
	// 单个商品价值
	Price int64 `db:"price" json:"price"`
	//抵扣
	Exemptions FieldsExemptions `db:"exemptions" json:"exemptions"`
	//是否买家评价
	CommentBuyer bool `db:"comment_buyer" json:"commentBuyer"`
	//买家评论ID
	CommentBuyerID int64 `db:"comment_buyer_id" json:"commentBuyerID"`
	//是否卖家评价
	CommentSeller bool `db:"comment_seller" json:"commentSeller"`
	//卖家评论ID
	CommentSellerID int64 `db:"comment_seller_id" json:"commentSellerID"`
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
	// 0 货物费用/预付款；1 配送费用；2 保险费用; 3 尾款
	PriceType int `db:"price_type" json:"priceType"`
	//支付ID
	PayID int64 `db:"pay_id" json:"payID"`
	//支付失败信息
	PayFailed string `db:"pay_failed" json:"payFailed"`
	//是否缴费
	IsPay bool `db:"is_pay" json:"isPay"`
	//金额
	Price int64 `db:"price" json:"price"`
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
	Price int64 `db:"price" json:"price"`
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

// FieldsLog 日志
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

// Value sql底层处理器
func (t FieldsLog) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsLog) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}
