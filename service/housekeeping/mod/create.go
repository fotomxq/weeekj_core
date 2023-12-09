package ServiceHousekeepingMod

import (
	CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"
	CoreSQLAddress "gitee.com/weeekj/weeekj_core/v5/core/sql/address"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

// ArgsCreateLog 创建新的请求参数
type ArgsCreateLog struct {
	//服务用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//预计上门时间
	NeedAt time.Time `db:"need_at" json:"needAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//服务负责人
	BindID int64 `db:"bind_id" json:"bindID" check:"id"`
	//其他服务参与人员
	OtherBinds pq.Int64Array `db:"other_binds" json:"otherBinds" check:"ids" empty:"true"`
	//服务项目商品ID
	MallProductID int64 `db:"mall_product_id" json:"mallProductID" check:"id"`
	//订单ID
	OrderID int64 `db:"order_id" json:"orderID" check:"id"`
	//货币
	Currency int `db:"currency" json:"currency" check:"currency"`
	//单位报价
	// 服务负责人在收款前可以协商变更
	Price int64 `db:"price" json:"price" check:"price" empty:"true"`
	//是否支付
	PayAt time.Time `db:"pay_at" json:"payAt"`
	//客户备注
	Des string `db:"des" json:"des" check:"des" min:"1" max:"500" empty:"true"`
	//客户地址
	Address CoreSQLAddress.FieldsAddress `db:"address" json:"address"`
	//服务配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateLog 创建新的请求
func CreateLog(args ArgsCreateLog) {
	CoreNats.PushDataNoErr("/service/housekeeping/create", "create", 0, "", args)
}
