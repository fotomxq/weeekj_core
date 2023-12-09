package ServiceAD

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

// FieldsApply 商户广告申请
type FieldsApply struct {
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
	//审核时间
	AuditAt time.Time `db:"audit_at" json:"auditAt"`
	//申请描述
	AuditDes string `db:"audit_des" json:"auditDes"`
	//拒绝原因
	AuditBanDes string `db:"audit_ban_des" json:"auditBanDes"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户ID
	// 可以由用户发起，用户发起的广告主要向商户投放，由商户进行审核处理
	UserID int64 `db:"user_id" json:"userID"`
	//投放分区ID列
	AreaIDs pq.Int64Array `db:"area_ids" json:"areaIDs"`
	//已经投放的广告ID
	AdID int64 `db:"ad_id" json:"adID"`
	//分区标识码
	// 作为前端抽取数据类型使用，可以重复指定多个
	Mark string `db:"mark" json:"mark"`
	//名称
	Name string `db:"name" json:"name"`
	//描述
	Des string `db:"des" json:"des"`
	//封面
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID"`
	//描述组图
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles"`
	//展示次数
	Count int64 `db:"count" json:"count"`
	//点击次数
	ClickCount int64 `db:"click_count" json:"clickCount"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
