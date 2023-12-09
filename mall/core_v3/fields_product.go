package MallCoreV3

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	CoreSQLAddress "gitee.com/weeekj/weeekj_core/v5/core/sql/address"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

// FieldsProduct 商品核心表
type FieldsProduct struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//上级ID
	// 标记为历史记录，<1为最高级，其他级别代表历史记录
	ParentID int64 `db:"parent_id" json:"parentID"`
	//商品类型
	// 0 普通商品; 1 关联选项商品
	ProductType int `db:"product_type" json:"productType" check:"intThan0" empty:"true"`
	//是否为虚拟商品
	// 不会发生配送处理
	IsVirtual bool `db:"is_virtual" json:"isVirtual" check:"bool"`
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID"`
	//排序
	Sort int `db:"sort" json:"sort"`
	//发布时间
	PublishAt time.Time `db:"publish_at" json:"publishAt"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID"`
	//标签ID列
	Tags pq.Int64Array `db:"tags" json:"tags"`
	//商品条形码编码
	Code string `db:"code" json:"code"`
	//标题
	Title string `db:"title" json:"title"`
	//副标题
	TitleDes string `db:"title_des" json:"titleDes"`
	//商品描述
	Des string `db:"des" json:"des"`
	//封面ID
	CoverFileIDs pq.Int64Array `db:"cover_file_ids" json:"coverFileIDs"`
	//描述图组
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles"`
	//货币
	Currency int `db:"currency" json:"currency"`
	//实际费用
	PriceReal int64 `db:"price_real" json:"priceReal"`
	//折扣截止
	PriceExpireAt time.Time `db:"price_expire_at" json:"priceExpireAt"`
	//折扣前费用
	Price int64 `db:"price" json:"price"`
	//积分价格
	Integral int64 `db:"integral" json:"integral"`
	//积分最多抵扣费用
	IntegralPrice int64 `db:"integral_price" json:"integralPrice"`
	//积分包含配送费免费
	IntegralTransportFree bool `db:"integral_transport_free" json:"integralTransportFree"`
	//会员价格
	// 会员配置分平台和商户，平台会员需参与活动才能使用，否则将禁止设置和后期使用
	UserSubPrice FieldsUserSubPrices `db:"user_sub_price" json:"userSubPrice"`
	//票据
	// 可以使用的票据列，具体的配置在票据配置内进行设置
	// 票据分平台和商户，平台票据需参与活动才能使用，否则将自动禁止设置和后期使用
	UserTicket pq.Int64Array `db:"user_ticket" json:"userTicket"`
	//配送费计费模版ID
	TransportID int64 `db:"transport_id" json:"transportID"`
	//唯一送货地址
	Address CoreSQLAddress.FieldsAddress `db:"address" json:"address"`
	//关联仓库的产品
	WarehouseProductID int64 `db:"warehouse_product_id" json:"warehouseProductID"`
	//货物重量
	// 单位g
	Weight int `db:"weight" json:"weight"`
	//总库存
	Count int `db:"count" json:"count"`
	//销量
	BuyCount int `db:"buy_count" json:"buyCount"`
	//关联附加选项
	OtherOptions FieldsOtherOption `db:"other_options" json:"otherOptions"`
	//给与票据列
	// 和赠礼区别在于，赠礼不可退，此票据会跟随订单取消设置是否退还
	GivingTickets FieldsGivingTickets `db:"giving_tickets" json:"givingTickets"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// FieldsUserSubPrice 会员价格
type FieldsUserSubPrice struct {
	//会员ID
	ID int64 `db:"id" json:"id"`
	//标定价格
	Price int64 `db:"price" json:"price"`
}

// Value sql底层处理器
func (t FieldsUserSubPrice) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsUserSubPrice) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsUserSubPrices []FieldsUserSubPrice

// Value sql底层处理器
func (t FieldsUserSubPrices) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsUserSubPrices) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

// FieldsGivingTickets 给与票据
type FieldsGivingTickets []FieldsGivingTicket

// Value sql底层处理器
func (t FieldsGivingTickets) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsGivingTickets) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsGivingTicket struct {
	//票据ID
	TicketConfigID int64 `db:"ticket_config_id" json:"ticketConfigID" check:"id"`
	//张数
	Count int64 `db:"count" json:"count" check:"int64Than0"`
	//是否可退
	IsRefund bool `db:"is_refund" json:"isRefund" check:"bool"`
}

// Value sql底层处理器
func (t FieldsGivingTicket) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsGivingTicket) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}
