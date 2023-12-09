package BaseStyle

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

// FieldsStyle 样式主表
type FieldsStyle struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//样式库名称
	Name string `db:"name" json:"name"`
	//关联标识码
	// 用于识别代码片段
	Mark string `db:"mark" json:"mark"`
	//样式使用渠道
	// app APP；wxx 小程序等，可以任意定义，模块内不做限制
	SystemMark string `db:"system_mark" json:"systemMark"`
	//分栏样式结构设计
	Components pq.Int64Array `db:"components" json:"components"`
	//默认标题
	// 标题是展示给用户的，样式库名称和该标题不是一个
	Title string `db:"title" json:"title"`
	//默认描述
	Des string `db:"des" json:"des"`
	//默认封面
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID"`
	//默认描述文件组
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID"`
	//标签组
	Tags pq.Int64Array `db:"tags" json:"tags"`
	//附加参数
	// 布局相关的自定会参数，都将在此处定义
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
