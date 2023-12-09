package TMSUserRunningMod

import (
	CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"
	CoreSQLAddress "gitee.com/weeekj/weeekj_core/v5/core/sql/address"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
)

// ArgsCreateMission 创建新的任务参数
type ArgsCreateMission struct {
	//跑腿单类型
	// 0 帮我送 ; 1 帮我买; 2 帮我取
	RunType int `db:"run_type" json:"runType" check:"intThan0" empty:"true"`
	//期望上门时间
	WaitAt string `db:"wait_at" json:"waitAt" check:"isoTime"`
	//物品类型
	GoodType string `db:"good_type" json:"goodType" check:"mark"`
	//关联组织
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//关联用户
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//关联订单ID
	// 可能没有关联订单
	OrderID int64 `db:"order_id" json:"orderID" check:"id" empty:"true"`
	//等待缴纳的费用
	// -1为自动重新计费，0为不需要缴纳费用，其他为需要缴纳的费用
	RunWaitPrice int64 `db:"run_wait_price" json:"runWaitPrice" check:"price" empty:"true"`
	//跑腿费是否货到付款
	RunPayAfter bool `db:"run_pay_after" json:"runPayAfter" check:"bool"`
	//订单是否已经缴纳了所有费用
	OrderPayAllPrice bool `json:"orderPayAllPrice" check:"bool"`
	//跑腿单描述信息
	Des string `db:"des" json:"des" check:"des" min:"1" max:"1200" empty:"true"`
	//物品重量
	GoodWidget int `db:"good_widget" json:"goodWidget" check:"intThan0" empty:"true"`
	//发货地址
	FromAddress CoreSQLAddress.FieldsAddress `db:"from_address" json:"fromAddress" check:"address_data" empty:"true"`
	//送货地址
	ToAddress CoreSQLAddress.FieldsAddress `db:"to_address" json:"toAddress" check:"address_data" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params" check:"params" empty:"true"`
}

// CreateMission 创建新的任务
func CreateMission(args ArgsCreateMission) {
	CoreNats.PushDataNoErr("/tms/user_running/new", "", 0, "", args)
}
