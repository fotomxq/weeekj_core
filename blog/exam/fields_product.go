package BlogExam

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

// FieldsProduct 考试项目
type FieldsProduct struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//开始时间
	StartAt time.Time `db:"start_at" json:"startAt"`
	//结束时间
	EndAt time.Time `db:"end_at" json:"endAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//标题
	Title string `db:"title" json:"title"`
	//题目列
	TopicIDs pq.Int64Array `db:"topic_ids" json:"topicIDs"`
	//是否直接反馈正确答案
	ReturnAnswerNow bool `db:"return_answer_now" json:"returnAnswerNow"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
