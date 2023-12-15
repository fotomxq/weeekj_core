package TMSTransportMod

import (
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQLAddress "github.com/fotomxq/weeekj_core/v5/core/sql/address"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
)

// ArgsCreateTransport 创建新配送单参数
type ArgsCreateTransport struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//当前配送人员
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//客户档案ID
	InfoID int64 `db:"info_id" json:"infoID" check:"id" empty:"true"`
	//客户用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//取货地址
	FromAddress CoreSQLAddress.FieldsAddress `db:"from_address" json:"fromAddress"`
	//收货地址
	ToAddress CoreSQLAddress.FieldsAddress `db:"to_address" json:"toAddress"`
	//订单ID
	OrderID int64 `db:"order_id" json:"orderID" check:"id" empty:"true"`
	//货物ID
	Goods FieldsTransportGoods `db:"goods" json:"goods"`
	//快递总重量
	Weight int `db:"weight" json:"weight" check:"intThan0" empty:"true"`
	//长宽
	Length int `db:"length" json:"length" check:"intThan0" empty:"true"`
	Width  int `db:"width" json:"width" check:"intThan0" empty:"true"`
	//货币
	Currency int `db:"currency" json:"currency" check:"currency"`
	//配送费用
	Price int64 `db:"price" json:"price" check:"price" empty:"true"`
	//是否完成了缴费
	PayFinish bool `db:"pay_finish" json:"payFinish" check:"bool" empty:"true"`
	//期望送货时间
	TaskAt string `db:"task_at" json:"taskAt" check:"isoTime" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateTransport 创建新配送单
func CreateTransport(args ArgsCreateTransport) {
	CoreNats.PushDataNoErr("/tms/transport/create", "", 0, "", args)
}
