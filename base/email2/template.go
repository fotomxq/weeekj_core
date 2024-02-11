package BaseEmail2

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
)

// ArgsCreateTemplate 创建模板参数
type ArgsCreateTemplate struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id" check:"id"`
	//计划使用的邮箱配置列
	// 多个配置将随机抽取一个发送
	ServerIDs pq.Int64Array `db:"server_ids" json:"serverIDs" check:"ids" empty:"true"`
	//标题
	Title string `db:"title" json:"title" check:"name"`
	//内容
	// 将强制邮件采用HTML模式发送，此处存放HTML内容
	// 相关变量根据模块的约定执行
	Content string `db:"content" json:"content" check:"des" min:"1" max:"3000" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateTemplate 修改模板
func UpdateTemplate(args *ArgsCreateTemplate) (err error) {
	//更新数据
	err = baseEmail2SQL.Update().SetFields([]string{"server_ids", "title", "content", "params"}).NeedUpdateTime().AddWhereID(args.ID).NamedExec(map[string]interface{}{
		"server_ids": args.ServerIDs,
		"title":      args.Title,
		"content":    args.Content,
		"params":     args.Params,
	})
	//反馈
	return
}
