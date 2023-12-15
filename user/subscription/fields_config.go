package UserSubscription

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

type FieldsConfig struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//标识码
	Mark string `db:"mark" json:"mark"`
	//时间类型
	// 0 小时 1 天 2 周 3 月 4 年
	TimeType int `db:"time_type" json:"timeType"`
	//时间长度
	TimeN int `db:"time_n" json:"timeN"`
	//开通价格
	Currency int   `db:"currency" json:"currency"`
	Price    int64 `db:"price" json:"price"`
	//折扣前费用，用于展示
	PriceOld int64 `db:"price_old" json:"priceOld"`
	//标题
	Title string `db:"title" json:"title"`
	//描述
	Des string `db:"des" json:"des"`
	//封面
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID"`
	//描述文件组
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles"`
	//关联的用户组
	// 只有为平台配置时，该数据才可修改并会生效
	UserGroups pq.Int64Array `db:"user_groups" json:"userGroups"`
	//默认减免的费用比例、费用金额
	ExemptionPrice int64 `db:"exemption_price" json:"exemptionPrice"`
	// 1-100% 百分比
	ExemptionDiscount int64 `db:"exemption_discount" json:"exemptionDiscount"`
	//费用低于多少时，将失效
	// 依赖于订单的总金额判断
	ExemptionMinPrice int64 `db:"exemption_min_price" json:"exemptionMinPrice"`
	//限制设计
	// 允许设置多个条件，如1天限制一次的同时、30天能使用10次
	Limits FieldsLimits `db:"limits" json:"limits"`
	//周期价格
	ExemptionTime FieldsExemptionTimes `db:"exemption_time" json:"exemptionTime"`
	//样式ID
	// 关联到样式库后，本记录的图片和文本将交给样式库布局实现
	StyleID int64 `db:"style_id" json:"styleID"`
	//附加参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// FieldsLimits 限制措施
type FieldsLimits []FieldsLimit

// Value sql底层处理器
func (t FieldsLimits) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsLimits) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsLimit struct {
	//时间类型
	// 0 小时 1 天 2 周 3 月 4 年
	TimeType int `db:"time_type" json:"timeType" check:"intThan0"`
	//时间长度
	TimeN int `db:"time_n" json:"timeN" check:"intThan0"`
	//限制的次数
	Count int `db:"count" json:"count" check:"intThan0"`
}

// Value sql底层处理器
func (t FieldsLimit) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsLimit) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

// FieldsExemptionTimes 限制措施
type FieldsExemptionTimes []FieldsExemptionTime

// Value sql底层处理器
func (t FieldsExemptionTimes) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsExemptionTimes) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsExemptionTime struct {
	//时间长度
	TimeN int `db:"time_n" json:"timeN" check:"intThan0"`
	//价格
	Price int64 `db:"price" json:"price" check:"price"`
}

// Value sql底层处理器
func (t FieldsExemptionTime) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsExemptionTime) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}
