package BaseEmail

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

// FieldsTemplate 模板组件
// 可以为各种模块声明模板，提供特定的邮件内容
type FieldsTemplate struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID"`
	//计划使用的邮箱配置列
	// 多个配置将随机抽取一个发送
	ServerIDs pq.Int64Array `db:"server_ids" json:"serverIDs"`
	//标题
	Title string `db:"title" json:"title"`
	//内容
	// 将强制邮件采用HTML模式发送，此处存放HTML内容
	// 相关变量根据模块的约定执行
	Content string `db:"content" json:"content"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
