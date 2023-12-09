package BlogExam

import (
	"github.com/lib/pq"
	"time"
)

// FieldsLog 考试记录
type FieldsLog struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//结束时间
	EndAt time.Time `db:"end_at" json:"endAt"`
	//耗时
	RunTime int `db:"run_time" json:"runTime"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//参加考试
	ProductID int64 `db:"product_id" json:"productID"`
	//抽取的题目
	// 考试可以为空，此设计不能为空
	TopicIDs pq.Int64Array `db:"topic_ids" json:"topicIDs"`
	//最终得分
	// 得分根据正确得分分数递增得到
	// 全局配置中，如果BlogExamScore100=true，则按照百分比计算：(正确数量 / 试卷总分数) x 10000
	Score int `db:"score" json:"score"`
	//错误数量
	ErrCount int `db:"err_count" json:"errCount"`
	//正确数量
	CorrectCount int `db:"correct_count" json:"correctCount"`
}
