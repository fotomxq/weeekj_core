package ServiceUserInfo

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

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
	OrgID int64 `db:"org_id" json:"orgID"`
	//名称
	Name string `db:"name" json:"name"`
	//描述
	Des string `db:"des" json:"des"`
	//封面图
	CoverFiles pq.Int64Array `db:"cover_files" json:"coverFiles"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID"`
	//标签
	Tags pq.Int64Array `db:"tags" json:"tags"`
	//文件数据
	FileData string `db:"file_data" json:"fileData"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
