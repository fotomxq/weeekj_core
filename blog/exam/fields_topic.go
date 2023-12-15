package BlogExam

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

// FieldsTopic 考试题目
type FieldsTopic struct {
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
	//描述
	Des string `db:"des" json:"des"`
	//题目类型
	// 0 单选； 1 多选； 2 判断； 3 填空题； 4 问答题
	TopicType int `db:"topic_type" json:"topicType"`
	//正确得分
	Score int `db:"score" json:"score"`
	//选项
	// 单选、多选、判断
	Options FieldsTopicOptions `db:"options" json:"options"`
	//正确选项
	// 可能是多个，用于支持单选、多选、判断、填空
	// 注意填空题此处如果设置，则需des配合填写{marK}的字符，方便直到具体是哪个mark
	Answer pq.StringArray `db:"answer" json:"answer"`
	//解析
	AnswerAnalysis string `db:"answer_analysis" json:"answerAnalysis"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

type FieldsTopicOptions []FieldsTopicOption

// Value sql底层处理器
func (t FieldsTopicOptions) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsTopicOptions) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsTopicOption struct {
	//选项
	Mark string `db:"mark" json:"mark" check:"mark"`
	//考题内容
	Des string `db:"des" json:"des" check:"des" min:"1" max:"600"`
}

// Value sql底层处理器
func (t FieldsTopicOption) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsTopicOption) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}
