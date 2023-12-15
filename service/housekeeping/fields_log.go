package ServiceHousekeeping

import (
	CoreSQLAddress "github.com/fotomxq/weeekj_core/v5/core/sql/address"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

// FieldsLog 服务申请和记录
type FieldsLog struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//编号
	// 商户下唯一
	SN int64 `db:"sn" json:"sn"`
	//今日编号
	SNDay int64 `db:"sn_day" json:"snDay"`
	//预计上门时间
	NeedAt time.Time `db:"need_at" json:"needAt"`
	//完成时间
	FinishAt time.Time `db:"finish_at" json:"finishAt"`
	//服务用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//服务负责人
	BindID int64 `db:"bind_id" json:"bindID"`
	//其他服务参与人员
	OtherBinds pq.Int64Array `db:"other_binds" json:"otherBinds"`
	//服务项目商品ID
	MallProductID int64 `db:"mall_product_id" json:"mallProductID"`
	//订单ID
	OrderID int64 `db:"order_id" json:"orderID"`
	//货币
	Currency int `db:"currency" json:"currency"`
	//单位报价
	// 服务负责人在收款前可以协商变更
	Price int64 `db:"price" json:"price"`
	//支付ID
	PayID int64 `db:"pay_id" json:"payID"`
	//支付时间
	PayAt time.Time `db:"pay_at" json:"payAt"`
	//客户备注
	Des string `db:"des" json:"des"`
	//客户地址
	Address CoreSQLAddress.FieldsAddress `db:"address" json:"address"`
	//服务配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
