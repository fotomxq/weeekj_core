package BaseStyle

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

// FieldsComponent 组件
// 按照顺序被样式主表定义引用
// 样式需要被申明
type FieldsComponent struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//关联标识码
	// 必填
	// 页面内独特的代码片段，声明后可直接定义该组件的默认参数形式
	Mark string `db:"mark" json:"mark"`
	//组件名称
	Name string `db:"name" json:"name"`
	//组件介绍
	Des string `db:"des" json:"des"`
	//组件封面
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID"`
	//组件描述文件组
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID"`
	//标签组
	Tags pq.Int64Array `db:"tags" json:"tags"`
	//附加参数
	// 布局相关的自定会参数，都将在此处定义
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
