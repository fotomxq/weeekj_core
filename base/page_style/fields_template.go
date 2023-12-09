package BasePageStyle

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

// FieldsTemplate 模版
type FieldsTemplate struct {
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
	//页面识别码
	// 同一个系统下唯一
	Page string `db:"page" json:"page"`
	//名称
	Name string `db:"name" json:"name"`
	//介绍
	Des string `db:"des" json:"des"`
	//封面
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
	//可用组件列
	ComponentIDs pq.Int64Array `db:"component_ids" json:"componentIDs"`
	//默认呈现的组件排序
	DefaultComponentIDs pq.Int64Array `db:"default_component_ids" json:"defaultComponentIDs"`
	//样式结构内容
	Data string `db:"data" json:"data"`
	//默认样式结构
	DefaultData string `db:"default_data" json:"defaultData"`
	//附加参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
