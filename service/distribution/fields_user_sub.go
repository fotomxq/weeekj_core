package ServiceDistribution

import "time"

//FieldsUserSub 特定用户会员模式
type FieldsUserSub struct {
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
	//对应分销商
	DistributionID int64 `db:"distribution_id" json:"distributionID"`
	//会员配置
	SubConfigID int64 `db:"sub_config_id" json:"subConfigID"`
	//强制约定价格
	UnitPrice int64 `db:"unit_price" json:"unitPrice"`
	//指定奖励
	MarketGivingSubID int64 `db:"market_giving_sub_id" json:"marketGivingSubID"`
	//宣传海报
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID"`
	//描述
	Des string `db:"des" json:"des"`
	//进入次数
	InCount int64 `db:"in_count" json:"inCount"`
	//交易次数
	OrderCount int64 `db:"order_count" json:"orderCount"`
	//交易金额
	OrderPrice int64 `db:"order_price" json:"orderPrice"`
}
