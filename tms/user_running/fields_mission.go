package TMSUserRunning

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	CoreSQLAddress "gitee.com/weeekj/weeekj_core/v5/core/sql/address"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

// FieldsMission 跑腿单数据
type FieldsMission struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//期望上门时间
	WaitAt time.Time `db:"wait_at" json:"waitAt" check:"isoTime"`
	//物品类型
	GoodType string `db:"good_type" json:"goodType" check:"mark"`
	//取货时间
	TakeAt time.Time `db:"take_at" json:"takeAt"`
	//是否完结
	FinishAt time.Time `db:"finish_at" json:"finishAt"`
	//取货码
	TakeCode string `db:"take_code" json:"takeCode"`
	//跑腿单类型
	// 0 帮我送 ; 1 帮我买; 2 帮我取
	RunType int `db:"run_type" json:"runType" check:"intThan0" empty:"true"`
	//关联组织
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//关联用户
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//关联订单ID
	// 可能没有关联订单
	OrderID int64 `db:"order_id" json:"orderID" check:"id" empty:"true"`
	//履约跑腿员
	// 用户角色ID
	RoleID int64 `db:"role_id" json:"roleID" check:"id" empty:"true"`
	//跑腿付费时间
	RunPayAt time.Time `db:"run_pay_at" json:"runPayAt" check:"isoTime" empty:"true"`
	//是否完成跑腿费支付
	RunPayID int64 `db:"run_pay_id" json:"runPayID" check:"id" empty:"true"`
	//跑腿费用总计
	// 已经支付的部分
	RunPrice int64 `db:"run_price" json:"runPrice" check:"price" empty:"true"`
	//等待缴纳的费用
	RunWaitPrice int64 `db:"run_wait_price" json:"runWaitPrice" check:"price" empty:"true"`
	//跑腿追加费用清单
	// 开始的跑腿费和追加费用，都会被列入此列表
	RunPayList pq.Int64Array `db:"run_pay_list" json:"runPayList" check:"ids" empty:"true"`
	//跑腿费是否货到付款
	RunPayAfter bool `db:"run_pay_after" json:"runPayAfter" check:"bool"`
	//订单是否货到付款
	OrderPayAfter bool `db:"order_pay_after" json:"orderPayAfter" check:"bool"`
	//商品等待缴纳费用
	OrderWaitPrice int64 `db:"order_wait_price" json:"orderWaitPrice" check:"price" empty:"true"`
	//订单费用
	OrderPrice int64 `db:"order_price" json:"orderPrice" check:"price" empty:"true"`
	//订单是否已经支付
	OrderPayAt time.Time `db:"order_pay_at" json:"orderPayAt" check:"isoTime" empty:"true"`
	//订单支付ID
	// 如果>0，则关联到订单；否则其他参数将用于商品的价格等信息描述
	OrderPayID int64 `db:"order_pay_id" json:"orderPayID" check:"id" empty:"true"`
	//平台服务费
	ServicePrice int64 `db:"service_price" json:"servicePrice" check:"price"`
	//跑腿单描述信息
	Des string `db:"des" json:"des" check:"des" min:"1" max:"1200" empty:"true"`
	//跑腿单核对订单数据包
	OrderDesFiles pq.Int64Array `db:"order_des_files" json:"orderDesFiles" check:"ids" empty:"true"`
	//跑腿单追加订单描述
	OrderDes string `db:"order_des" json:"orderDes" check:"des" min:"1" max:"3000" empty:"true"`
	//物品重量
	GoodWidget int `db:"good_widget" json:"goodWidget" check:"intThan0" empty:"true"`
	//发货地址
	FromAddress CoreSQLAddress.FieldsAddress `db:"from_address" json:"fromAddress" check:"address_data" empty:"true"`
	//送货地址
	ToAddress CoreSQLAddress.FieldsAddress `db:"to_address" json:"toAddress" check:"address_data" empty:"true"`
	//日志
	Logs FieldsMissionLogs `db:"logs" json:"logs"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params" check:"params" empty:"true"`
}

type FieldsMissionLogs []FieldsMissionLog

// Value sql底层处理器
func (t FieldsMissionLogs) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsMissionLogs) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsMissionLog struct {
	//时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//日志描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"600"`
	//附加文件
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
}

// Value sql底层处理器
func (t FieldsMissionLog) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsMissionLog) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}
