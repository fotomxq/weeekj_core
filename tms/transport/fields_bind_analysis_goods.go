package TMSTransport

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

//FieldsBindAnalysisGoods 配送员配送货物统计
type FieldsBindAnalysisGoods struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	// 每个配送人员每天每个货物，只会产生一条数据
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//配送人员
	BindID int64 `db:"bind_id" json:"bindID"`
	//是否为退单
	IsRefund bool `db:"is_refund" json:"isRefund"`
	//货物列
	Goods FieldsBindAnalysisGoodsGoods `db:"goods" json:"goods"`
}

type FieldsBindAnalysisGoodsGoods []FieldsBindAnalysisGoodsGood

//Value sql底层处理器
func (t FieldsBindAnalysisGoodsGoods) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsBindAnalysisGoodsGoods) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsBindAnalysisGoodsGood struct {
	//来源
	System string `db:"system" json:"system"`
	//ID
	ID int64 `db:"id" json:"id"`
	//订单创建渠道
	FromSystem int `db:"from_system" json:"fromSystem"`
	//支付渠道
	// 支付system+_+mark
	PaySystem string `db:"pay_system" json:"paySystem"`
	//数量
	Count int64 `db:"count" json:"count"`
	//订单付款金额
	Price int64 `db:"price" json:"price"`
}

//Value sql底层处理器
func (t FieldsBindAnalysisGoodsGood) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsBindAnalysisGoodsGood) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}