package BlogStuRead

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

type FieldsContent struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//文章类型
	// 0 普通文章; 1 挂靠视频; 3 第三方跳转
	// 挂靠或跳转内容，将在des中做特殊描述
	ContentType int `db:"content_type" json:"contentType"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//访问量
	VisitCount int64 `db:"visit_count" json:"visitCount"`
	//标题
	Title string `db:"title" json:"title"`
	//小标题
	TitleDes string `db:"title_des" json:"titleDes"`
	//封面文件
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID"`
	//附加封面图
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//内容
	Des string `db:"des" json:"des"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
