package ServiceAD

import "time"

//FieldsAnalysis 统计记录
type FieldsAnalysis struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//统计周期
	// 完全相同的一个来源体系，1小时仅构建一条数据
	DayTime time.Time `db:"day_time" json:"dayTime"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//分区ID
	AreaID int64 `db:"area_id" json:"areaID"`
	//广告ID
	AdID int64 `db:"ad_id" json:"adID"`
	//投放次数
	Count int64 `db:"count" json:"count"`
	//点击次数
	ClickCount int64 `db:"click_count" json:"clickCount"`
}
