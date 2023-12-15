package OrgReport

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

// FieldsReport 投诉建议
type FieldsReport struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//投诉来源
	// 可以全部为0或空，则代表匿名
	FromSystem string `db:"from_system" json:"fromSystem"`
	FromID     int64  `db:"from_id" json:"fromID"`
	FromName   string `db:"from_name" json:"fromName"`
	//建议内容
	Des string `db:"des" json:"des"`
	//建议附图
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles"`
	//投诉目标
	// 可以不含投诉目标，具体看投诉人意愿和业务逻辑需要
	TargetSystem string `db:"target_system" json:"targetSystem"`
	TargetID     int64  `db:"target_id" json:"targetID"`
	TargetName   string `db:"target_name" json:"targetName"`
	//反馈内容
	ReportAt    time.Time     `db:"report_at" json:"reportAt"`
	ReportDes   string        `db:"report_des" json:"reportDes"`
	ReportFiles pq.Int64Array `db:"report_files" json:"reportFiles"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
