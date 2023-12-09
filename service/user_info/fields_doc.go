package ServiceUserInfo

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

// FieldsDoc 归档档案结构
type FieldsDoc struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//组织ID
	// 允许平台方的0数据，该数据可能来源于其他领域
	OrgID int64 `db:"org_id" json:"orgID"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID"`
	//标签
	Tags pq.Int64Array `db:"tags" json:"tags"`
	//档案ID
	InfoID int64 `db:"info_id" json:"infoID"`
	//文档标题
	Title string `db:"title" json:"title"`
	//模版ID
	TemplateID int64 `db:"template_id" json:"templateID"`
	//文件数据
	FileData string `db:"file_data" json:"fileData"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
