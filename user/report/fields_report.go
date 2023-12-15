package UserReport

import (
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	"github.com/lib/pq"
	"time"
)

// FieldsReport 反馈结构
type FieldsReport struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//反馈时间
	ReportAt time.Time `db:"report_at" json:"reportAt"`
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//举报内容来源
	FromInfo CoreSQLFrom.FieldsFrom `db:"from_info" json:"fromInfo"`
	//IP
	IP string `db:"ip" json:"ip"`
	//用户昵称
	UserName string `db:"user_name" json:"userName"`
	//绑定手机号的国家代码
	NationCode string `db:"nation_code" json:"nationCode"`
	//手机号码，绑定后的手机
	Phone string `db:"phone" json:"phone"`
	//邮箱，如果不存在手机则必须存在的
	Email string `db:"email" json:"email"`
	//截图
	Files pq.Int64Array `db:"files" json:"files"`
	//举报内容描述
	Content string `db:"content" json:"content"`
	//反馈人
	ReportUserID int64 `db:"report_user_id" json:"reportUserID"`
	//反馈内容
	ReportContent string `db:"report_content" json:"reportContent"`
}
