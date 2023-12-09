package ServiceAD

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"time"
)

// FieldsBind 分区和广告的绑定关系
type FieldsBind struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//启动时间
	StartAt time.Time `db:"start_at" json:"startAt"`
	//结束时间
	EndAt time.Time `db:"end_at" json:"endAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//分区ID
	AreaID int64 `db:"area_id" json:"areaID"`
	//广告ID
	AdID int64 `db:"ad_id" json:"adID"`
	//权重因子
	// 当一个分区存在多个广告绑定时，权重将作为随机的重要计算参数，影响随机的倾向性
	// 同一个分区下将所有广告进行合计，然后取随机数，在因子+偏移值范围的广告将作为最终投放标定
	// 注意，广告也可能会被用户数据分析模块影响，该模块将注入新的因子作为投放参数，作为因子2存在。
	// 投放概率 = 广告因子 + 用户数据分析因子 / (因子总数)
	Factor int `db:"factor" json:"factor"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
