package ClassComment

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

type FieldsComment struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//评论ID
	// 评论被删除后出现，指向新的评论
	// 下级评论的上级江全部改为该新的ID
	CommentID int64 `db:"comment_id" json:"commentID"`
	//上级ID
	// 评论的上下级关系，一旦建立无法修改
	ParentID int64 `db:"parent_id" json:"parentID"`
	//绑定组织
	// 该组织根据资源来源设定
	// 如果是平台资源，则为0
	OrgID int64 `db:"org_id" json:"orgID"`
	//所属用户
	UserID int64 `db:"user_id" json:"userID"`
	//绑定内容
	BindID int64 `db:"bind_id" json:"bindID"`
	//评价类型
	// 0 好评 1 中立 2 差评
	LevelType int `db:"level_type" json:"levelType"`
	//分数
	Level int `db:"level" json:"level"`
	//标题
	Title string `db:"title" json:"title"`
	//内容
	Des string `db:"des" json:"des"`
	//介绍图文
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
