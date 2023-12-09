package UserTicket

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
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
	//默认过期时间
	DefaultExpireTime int64 `db:"default_expire_time" json:"defaultExpireTime"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//标识码
	Mark string `db:"mark" json:"mark"`
	//标题
	Title string `db:"title" json:"title"`
	//描述
	Des string `db:"des" json:"des"`
	//封面
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID"`
	//描述文件组
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles"`
	//默认减免的费用比例、费用金额
	ExemptionPrice int64 `db:"exemption_price" json:"exemptionPrice"`
	// 1-100% 百分比
	ExemptionDiscount int64 `db:"exemption_discount" json:"exemptionDiscount"`
	//费用低于多少时，将失效
	// 依赖于订单的总金额判断
	ExemptionMinPrice int64 `db:"exemption_min_price" json:"exemptionMinPrice"`
	//是否可用于订单抵扣
	// 否则只能用于一件商品的抵扣
	UseOrder bool `db:"use_order" json:"useOrder"`
	//领取周期类型
	// 0 不限制; 1 一次性; 2 每天限制; 3 每周限制; 4 每月限制; 5 每季度限制; 6 每年限制
	LimitTimeType int `db:"limit_time_type" json:"limitTimeType"`
	//领取次数
	LimitCount int `db:"limit_count" json:"limitCount"`
	//样式ID
	// 关联到样式库后，本记录的图片和文本将交给样式库布局实现
	StyleID int64 `db:"style_id" json:"styleID"`
	//附加参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
