package ServiceAppointment

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

type FieldsAppointment struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//预约时间，从订单的时间抽取得到，此处方便核对
	// 修改时间后，将自动同步修改订单的预约时间
	WaitAt time.Time `db:"wait_at" json:"waitAt"`
	//是否核算
	FinishAt time.Time `db:"finish_at" json:"finishAt"`
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//创建来源
	// 0 商户平台下单、1 用户APP下单、2 微信小程序下单、3 支付宝小程序下单、4 电话系统下单、5 线下柜台下单、6 辅助设备下单、7 其他渠道
	CreateFrom int `db:"create_from" json:"createFrom"`
	//编号器，提供累计编号
	SerialNumber int64 `db:"serial_number" json:"serialNumber"`
	//当天的编号
	SerialNumberDay int64 `db:"serial_number_day" json:"serialNumberDay"`
	//支付ID
	PayID int64 `db:"pay_id" json:"payID"`
	//支付ID列
	// 所有关联请求，最后一条为最新的匹配数据
	PayList pq.Int64Array `db:"pay_list" json:"payList"`
	//费用
	Price int64 `db:"price" json:"price"`
	//是否支付
	PayAt time.Time `db:"pay_at" json:"payAt"`
	//附加费用
	AddPrice int64     `db:"add_price" json:"addPrice"`
	AddPayID int64     `db:"add_pay_id" json:"addPayID"`
	AddPayAt time.Time `db:"add_pay_at" json:"addPayAt"`
	//预约服务内容
	MallProductID int64 `db:"mall_product_id" json:"mallProductID"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
