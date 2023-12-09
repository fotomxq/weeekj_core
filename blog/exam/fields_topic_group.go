package BlogExam

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

// FieldsTopicGroup 题库设计
type FieldsTopicGroup struct {
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
	//访问量
	VisitCount int64 `db:"visit_count" json:"visitCount"`
	//题库名称
	Title string `db:"title" json:"title"`
	//题目列
	TopicIDs pq.Int64Array `db:"topic_ids" json:"topicIDs"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
