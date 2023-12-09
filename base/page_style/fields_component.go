package BasePageStyle

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

// FieldsComponent 组件
type FieldsComponent struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//系统
	System string `db:"system" json:"system"`
	//关联标识码
	// 全局唯一
	Mark string `db:"mark" json:"mark"`
	//组件名称
	Name string `db:"name" json:"name"`
	//组件介绍
	Des string `db:"des" json:"des"`
	//组件封面
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID"`
	//标签组
	Tags pq.Int64Array `db:"tags" json:"tags"`
	//商户订阅
	// 必须存在商户订阅配置的订阅，才能使用该组件
	OrgSubConfigID pq.Int64Array `db:"org_sub_config_id" json:"orgSubConfigID"`
	//商户功能
	// 只有开通相关功能，才能使用使用该组件
	OrgFuncList pq.StringArray `db:"org_func_list" json:"orgFuncList"`
	//样式结构内容
	Data string `db:"data" json:"data"`
	//附加参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
