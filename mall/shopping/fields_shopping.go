package MallShopping

import (
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DataShoppingOne 购物车数量结构
type DataShoppingOne struct {
	//数据结构
	DataList map[string]int64 `bson:"DataList" json:"dataList"`
}

// DataShopping 购物车结构
type DataShopping struct {
	//数据结构
	DataList []DataShoppingChild `bson:"DataList" json:"dataList"`
	//购买总数量
	Count int64 `bson:"Count" json:"count"`
	//总的商品标准价格
	PriceCount FieldsPrice `bson:"PriceCount" json:"priceCount"`
	//总的折扣后价格
	PriceCountLast FieldsPrice `bson:"PriceCountLast" json:"priceCountLast"`
	//总的配送费用
	PriceTransport FieldsPrice `bson:"PriceTransport" json:"priceTransport"`
	//总的保险费用
	PriceInsurance FieldsPrice `bson:"PriceInsurance" json:"priceInsurance"`
	//会员折扣总费用
	// 会员减免了多少钱
	SubPrice FieldsPrice `bson:"SubPrice" json:"subPrice"`
	//票据折扣总费用
	TicketPrice FieldsPrice `bson:"TicketPrice" json:"ticketPrice"`
	//折扣了总费用
	DiscountsPrice FieldsPrice `bson:"DiscountsPrice" json:"discountsPrice"`
}

type DataShoppingChild struct {
	//商品ID
	ID string `bson:"ID" json:"id"`
	//该商品数量
	Count int64 `bson:"Count" json:"count"`
	//商品标准价格
	PriceCount FieldsPrice `bson:"PriceCount" json:"priceCount"`
	//折扣后价格
	PriceCountLast FieldsPrice `bson:"PriceCountLast" json:"priceCountLast"`
	//标题
	Title string `bson:"Title" json:"title"`
	//封面URL
	// 商品第一张封面图
	CoverFile string `bson:"CoverFile" json:"coverFile"`
}

// DataShoppingOver 结算系统结构
type DataShoppingOver struct {
	//购物车结构
	ShoppingList []FieldsShopping `json:"shoppingList"`
	//商品标准价格
	PriceCount FieldsPrice `bson:"PriceCount" json:"priceCount"`
	//折扣费用
	// 打折后的实际费用，如果不存在请和标准价格一致
	PriceCountLast FieldsPrice `bson:"PriceCountLast" json:"priceCountLast"`
	//配送费用
	PriceTransport FieldsPrice `bson:"PriceTransport" json:"priceTransport"`
	//保险费用
	PriceInsurance FieldsPrice `bson:"PriceInsurance" json:"priceInsurance"`
	//购物车商品总数
	// 只包含允许购买的部分
	CommodityCount int64 `bson:"CommodityCount" json:"commodityCount"`
	//会员折扣总费用
	// 会员减免了多少钱
	SubPrice FieldsPrice `bson:"SubPrice" json:"subPrice"`
	//票据折扣总费用
	TicketPrice FieldsPrice `bson:"TicketPrice" json:"ticketPrice"`
	//折扣了总费用
	DiscountsPrice FieldsPrice `bson:"DiscountsPrice" json:"discountsPrice"`
}

// FieldsShopping 主表
type FieldsShopping struct {
	//基础
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	CreateTime int64              `bson:"CreateTime" json:"createTime"`
	UpdateTime int64              `bson:"UpdateTime" json:"updateTime"`
	//用户ID
	UserID string `bson:"UserID" json:"userID"`
	//所属绑定关系
	// 可以指定任意系统、ID、mark，只要有一个不同，则重新建立新的数据
	// 相同的数据将自动叠加计算
	// 该数据可用于积分商城和普通商城，以及未来其他业务逻辑的前端拆分处理
	FromInfo CoreSQLFrom.FieldsFrom `bson:"FromInfo" json:"fromInfo"`
	//购物车商品内容
	// 允许购买和符合条件的
	CommodityList []FieldsCommodity `bson:"CommodityList" json:"commodityList"`
	//废弃的商品
	// 下架的商品
	CommodityTrashList []FieldsCommodityTrash `bson:"CommodityTrashList" json:"commodityTrashList"`
	//商品标准价格
	PriceCount FieldsPrice `bson:"PriceCount" json:"priceCount"`
	//折扣费用
	// 打折后的实际费用，如果不存在请和标准价格一致
	PriceCountLast FieldsPrice `bson:"PriceCountLast" json:"priceCountLast"`
	//配送费用
	PriceTransport FieldsPrice `bson:"PriceTransport" json:"priceTransport"`
	//保险费用
	PriceInsurance FieldsPrice `bson:"PriceInsurance" json:"priceInsurance"`
	//购物车商品总数
	// 只包含允许购买的部分
	CommodityCount int64 `bson:"CommodityCount" json:"commodityCount"`
	//会员折扣总费用
	// 会员减免了多少钱
	SubPrice FieldsPrice `bson:"SubPrice" json:"subPrice"`
	//票据折扣总费用
	TicketPrice FieldsPrice `bson:"TicketPrice" json:"ticketPrice"`
	//折扣了总费用
	DiscountsPrice FieldsPrice `bson:"DiscountsPrice" json:"discountsPrice"`
}

// FieldsCommodityTrash 商品数据包
// 废弃商品清单
type FieldsCommodityTrash struct {
	//ID
	ID string `bson:"ID" json:"id"`
	//购买数量
	Count int64 `bson:"Count" json:"count"`
	//标题
	Title string `bson:"Title" json:"title"`
	//封面URL
	// 商品第一张封面图
	CoverFile string `bson:"CoverFile" json:"coverFile"`
}

// FieldsCommodity 商品数据包
type FieldsCommodity struct {
	//ID
	ID string `bson:"ID" json:"id"`
	//购买数量
	Count int64 `bson:"Count" json:"count"`
	//商品库存
	// 还是能够写入购物车，但由于库存不足，可能下单时会被拒绝
	// 添加时，如果超出该总数，将拒绝
	// 如果总数低于1，也同样拒绝
	GoodCount int64 `bson:"GoodCount" json:"goodCount"`
	//商品标准价格
	// 包含商品数量
	PriceCount FieldsPrice `bson:"PriceCount" json:"priceCount"`
	//折扣费用
	// 打折后的实际费用，如果不存在请和标准价格一致
	// 包含商品数量
	PriceCountLast FieldsPrice `bson:"PriceCountLast" json:"priceCountLast"`
	//配送费用
	// 包含商品数量
	PriceTransport FieldsPrice `bson:"PriceTransport" json:"priceTransport"`
	//保险费用
	// 包含商品数量
	PriceInsurance FieldsPrice `bson:"PriceInsurance" json:"priceInsurance"`
	//单价数据
	// 商品单价数据，不含购买数量
	PriceOne FieldsCommodityOne `bson:"PriceOne" json:"priceOne"`
	//抵扣所使用的票据
	// 注意，只计算选择票据的抵扣价格
	Tickets []FieldsTicket `bson:"Tickets" json:"tickets"`
	//抵扣的订阅服务
	Subscriptions []FieldsSubscription `bson:"Subscriptions" json:"subscriptions"`
	//标题
	Title string `bson:"Title" json:"title"`
	//封面URL
	// 商品第一张封面图
	CoverFile string `bson:"CoverFile" json:"coverFile"`
}

// FieldsCommodityOne 商品单价数据
type FieldsCommodityOne struct {
	//商品标准价格
	PriceCount FieldsPrice `bson:"PriceCount" json:"priceCount"`
	//折扣费用
	// 打折后的实际费用，如果不存在请和标准价格一致
	PriceCountLast FieldsPrice `bson:"PriceCountLast" json:"priceCountLast"`
	//配送费用
	PriceTransport FieldsPrice `bson:"PriceTransport" json:"priceTransport"`
	//保险费用
	PriceInsurance FieldsPrice `bson:"PriceInsurance" json:"priceInsurance"`
}

// FieldsTicket 抵扣票据
type FieldsTicket struct {
	//使用的票据配置ID
	ConfigID string `bson:"ConfigID" json:"configID"`
	//名称
	Name string `bson:"name" json:"name"`
	//使用的票据张数
	// 允许使用的最大票据张数
	// 抵扣订阅的服务次数
	Count int `bson:"Count" json:"count"`
	//是否选择该票据
	IsUse bool `bson:"IsUse" json:"isUse"`
	//抵扣的费用额度
	// 该额度根据货物价值及货币类型一致
	// 该抵扣额度为总数，不是每个商品的抵扣费用
	PriceCount int64 `bson:"PriceCount" json:"priceCount"`
	//抵扣的产品范围
	// 用于指定票据、订阅的使用渠道，并检查使用渠道的合法性
	From CoreSQLFrom.FieldsFrom `bson:"From" json:"from"`
	//抵扣的货物个数
	GoodCount int64 `bson:"GoodCount" json:"goodCount"`
}

// FieldsSubscription 抵扣订阅
type FieldsSubscription struct {
	//订阅配置ID
	ConfigID string `bson:"ConfigID" json:"configID"`
	//名称
	Name string `bson:"name" json:"name"`
	//抵扣的费用额度
	// 该额度根据货物价值及货币类型一致
	PriceCount int64 `bson:"PriceCount" json:"priceCount"`
	//抵扣的产品范围
	From CoreSQLFrom.FieldsFrom `bson:"From" json:"from"`
	//抵扣的货物个数
	GoodCount int64 `bson:"GoodCount" json:"goodCount"`
}

// FieldsPrice 费用
type FieldsPrice struct {
	//货物货币类型
	Currency string `bson:"Currency" json:"currency"`
	//费用
	Price int64 `bson:"Price" json:"price"`
}
