package ToolsHelpContent

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

// FieldsContent 内容设计
// 设置一整套内容和介绍文案，方便调取查阅
type FieldsContent struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//文本唯一标识码
	// none为预留值，指定后可以重复且该值将无效
	// 用于在不同页面使用
	// 删除后，将不会占用该mark设置，一个mark可以指定一个正常数据和多个已删除数据
	Mark string `db:"mark" json:"mark"`
	//是否公开
	// 非公开数据将作为草稿或私有数据存在，只有管理员可以看到
	IsPublic bool `db:"is_public" json:"isPublic"`
	//分类
	SortID int64 `db:"sort_id" json:"sortID"`
	//标签
	Tags pq.Int64Array `db:"tags" json:"tags"`
	//标题
	Title string `db:"title" json:"title"`
	//封面文件
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID"`
	//内容
	Des string `db:"des" json:"des"`
	//关联阅读引导ID
	BindIDs pq.Int64Array `db:"bind_ids" json:"bindIDs"`
	//关联阅读引导mark
	BindMarks pq.StringArray `db:"bind_marks" json:"bindMarks"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
