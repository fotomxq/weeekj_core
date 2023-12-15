package TMSTransport

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	CoreSQLAddress "github.com/fotomxq/weeekj_core/v5/core/sql/address"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

// FieldsTransport 配送
type FieldsTransport struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//完成时间
	FinishAt time.Time `db:"finish_at" json:"finishAt"`
	//期望送货时间
	TaskAt time.Time `db:"task_at" json:"taskAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//当前配送人员
	// 组织成员ID
	BindID int64 `db:"bind_id" json:"bindID"`
	//客户档案ID
	InfoID int64 `db:"info_id" json:"infoID"`
	//客户用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//编号
	// 商户下唯一
	SN int64 `db:"sn" json:"sn"`
	//今日编号
	SNDay int64 `db:"sn_day" json:"snDay"`
	//配送状态
	// 0 等待分配人员; 1 取货中; 2 送货中; 3 完成配送
	Status int `db:"status" json:"status"`
	//取货地址
	FromAddress CoreSQLAddress.FieldsAddress `db:"from_address" json:"fromAddress"`
	//收货地址
	ToAddress CoreSQLAddress.FieldsAddress `db:"to_address" json:"toAddress"`
	//订单ID
	OrderID int64 `db:"order_id" json:"orderID"`
	//货物ID
	Goods FieldsTransportGoods `db:"goods" json:"goods"`
	//快递总重量
	Weight int `db:"weight" json:"weight"`
	//长宽
	Length int `db:"length" json:"length"`
	Width  int `db:"width" json:"width"`
	//货币
	Currency int `db:"currency" json:"currency"`
	//配送费用
	Price int64 `db:"price" json:"price"`
	//完成缴费时间
	PayFinishAt time.Time `db:"pay_finish_at" json:"payFinishAt"`
	//缴费交易ID
	PayID int64 `db:"pay_id" json:"payID"`
	//历史支付请求
	PayIDs pq.Int64Array `db:"pay_ids" json:"payIDs" check:"ids" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// FieldsTransportGood 货物
type FieldsTransportGood struct {
	//来源系统
	System string `json:"system"`
	//来源ID
	ID int64 `json:"id"`
	//标识码
	Mark string `json:"mark"`
	//名称
	Name string `json:"name"`
	//数量
	Count int `json:"count"`
}

// Value sql底层处理器
func (t FieldsTransportGood) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsTransportGood) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsTransportGoods []FieldsTransportGood

// Value sql底层处理器
func (t FieldsTransportGoods) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsTransportGoods) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}
