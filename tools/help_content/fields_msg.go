package ToolsHelpContent

import (
	"github.com/lib/pq"
	"time"
)

//交互式交流服务
// 用户可采用文本框形式浏览，系统记录并反馈符合条件的帮助服务内容
// 转入客户服务将接入所在分区的客服代表，技术开发商不对客户提供直接的沟通服务，只对被服务的企业提供技术咨询服务

/**
技术原理：
1、用户发起消息
2、系统收到消息，对消息进行分词拆分
3、将分词模糊搜索内容系统的tag，反馈匹配项到语料库
4、语料库进行学习分析，取出最大概率接近的tag组，并按照概率顺序从高到低排列
5、搜索符合条件的tag词条，出现最多的词条tag作为符合条件反馈给用户
6、【可选】如存在多个词条符合条件，抽取任意一个反馈。同时该数据将列入学习机等待训练
 */

type FieldsMsg struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//是否为用户消息，否为系统消息
	IsSelf bool `db:"is_self" json:"isSelf"`
	//聊天内容
	Content string `db:"content" json:"content"`
	//携带文件
	FileID int64 `db:"file_id" json:"fileID"`
	//关联阅读引导ID
	BindIDs pq.Int64Array `db:"bind_ids" json:"bind_ids"`
}
